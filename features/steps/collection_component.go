package steps

import (
	"context"
	"github.com/ONSdigital/dp-collection-api/config"
	"github.com/ONSdigital/dp-collection-api/mongo"
	"github.com/ONSdigital/dp-collection-api/service"
	"github.com/ONSdigital/dp-collection-api/service/mock"
	"github.com/benweissmann/memongo"
	"net/http"
	"time"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CollectionComponent struct {
	componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
	apiFeature     *componenttest.APIFeature
}

func NewCollectionComponent(mongoFeature *componenttest.MongoFeature) (*CollectionComponent, error) {

	c := &CollectionComponent{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	mongodb := &mongo.Mongo{
		Database:              memongo.RandomDatabase(),
		URI:                   mongoFeature.Server.URI(),
		CollectionsCollection: c.Config.MongoConfig.CollectionsCollection,
	}

	if err := mongodb.Init(context.TODO()); err != nil {
		return nil, err
	}

	service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
		return mongodb, nil
	}

	return c, nil
}

func (c *CollectionComponent) Reset() *CollectionComponent {
	c.apiFeature.Reset()
	return c
}

func (c *CollectionComponent) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *CollectionComponent) InitialiseService() (http.Handler, error) {
	var err error
	ctx := context.Background()

	service.GetHTTPServer = c.GetHTTPServer
	service.GetHealthCheck = c.GetHealthCheck

	c.svc, err = service.New(ctx, c.Config, "1", "", "")
	if err != nil {
		return nil, err
	}

	c.svc.Start(ctx, c.errorChan)
	c.ServiceRunning = true
	return c.HTTPServer.Handler, nil
}

func (c *CollectionComponent) GetHealthCheck(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}
}

func (c *CollectionComponent) GetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}
