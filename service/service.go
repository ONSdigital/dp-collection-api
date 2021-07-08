package service

import (
	"context"
	"github.com/ONSdigital/dp-collection-api/mongo"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-collection-api/api"
	"github.com/ONSdigital/dp-collection-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var GetHealthCheck = func(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) HealthChecker {
	hc := healthcheck.New(version, criticalTimeout, interval)
	return &hc
}

var GetHTTPServer = func(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

var GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (MongoDB, error) {
	mongodb := &mongo.Mongo{
		CollectionsCollection: cfg.CollectionsCollection,
		EventsCollection:      cfg.EventsCollection,
		Database:              cfg.CollectionsDatabase,
		Username:              cfg.Username,
		Password:              cfg.Password,
		IsSSL:                 cfg.IsSSL,
		URI:                   cfg.BindAddr,
	}
	err := mongodb.Init()
	if err != nil {
		return nil, err
	}
	return mongodb, nil
}

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	cfg         *config.Config
	server      HTTPServer
	router      *mux.Router
	api         *api.API
	healthCheck HealthChecker
	mongoDB     MongoDB
}

// New initialises all the service dependencies
func New(ctx context.Context, cfg *config.Config, buildTime, gitCommit, version string) (*Service, error) {

	// Get HealthCheck and register checkers
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "error creating version info", err)
		return nil, err
	}

	mongoDB, err := GetMongoDB(ctx, cfg.MongoConfig)
	if err != nil {
		log.Fatal(ctx, "failed to initialise mongo db", err)
		return nil, err
	}

	healthCheck := GetHealthCheck(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err := registerHealthChecks(ctx, healthCheck, mongoDB); err != nil {
		return nil, errors.Wrap(err, "unable to register health checks")
	}

	// Get HTTP router and server with middleware
	r := mux.NewRouter()
	r.StrictSlash(true).Path("/health").HandlerFunc(healthCheck.Handler)
	server := GetHTTPServer(cfg.BindAddr, r)

	paginator := pagination.NewPaginator(cfg.DefaultLimit, cfg.DefaultOffset, cfg.DefaultMaxLimit)

	api := api.Setup(ctx, r, paginator, mongoDB)

	return &Service{
		cfg:         cfg,
		server:      server,
		router:      r,
		api:         api,
		healthCheck: healthCheck,
		mongoDB:     mongoDB,
	}, nil
}

// Start the service, allowing it to serve HTTP requests
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {

	svc.healthCheck.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		log.Event(ctx, "starting api", log.INFO)
		if err := svc.server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.cfg.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop health check, as it depends on everything else
		if svc.healthCheck != nil {
			svc.healthCheck.Stop()
		}

		// stop any incoming requests
		if svc.server != nil {
			if err := svc.server.Shutdown(ctx); err != nil {
				log.Error(ctx, "failed to shutdown http server", err)
				hasShutdownError = true
			}
		}

		if svc.mongoDB != nil {
			if err := svc.mongoDB.Close(ctx); err != nil {
				log.Error(ctx, "error closing mongo db client", err)
				hasShutdownError = true
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

// registerHealthChecks adds the checkers for the service clients to the health check object.
func registerHealthChecks(ctx context.Context, hc HealthChecker, mongoDB MongoDB) (err error) {

	hasErrors := false

	if err = hc.AddCheck("Mongo DB", mongoDB.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for mongo db client", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for health check")
	}
	return nil
}
