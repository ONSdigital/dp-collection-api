package steps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-collection-api/config"
	"github.com/ONSdigital/dp-collection-api/mongo"
	"github.com/ONSdigital/dp-collection-api/service"
	"github.com/ONSdigital/dp-collection-api/service/mock"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type CollectionComponent struct {
	componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	config         *config.Config
	httpServer     *http.Server
	serviceRunning bool
	apiFeature     *componenttest.APIFeature
	mongoClient    *mongo.Mongo
}

func NewCollectionComponent(mongoFeature *componenttest.MongoFeature) (*CollectionComponent, error) {

	c := &CollectionComponent{
		httpServer:     &http.Server{},
		errorChan:      make(chan error),
		serviceRunning: false,
	}

	var err error

	c.config, err = config.Get()
	if err != nil {
		return nil, err
	}

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	mongoURI := fmt.Sprintf("localhost:%d", mongoFeature.Server.Port())
	c.mongoClient = &mongo.Mongo{
		Database:              utils.RandomDatabase(),
		URI:                   mongoURI,
		CollectionsCollection: c.config.MongoConfig.CollectionsCollection,
		EventsCollection:      c.config.MongoConfig.EventsCollection,
	}

	if err := c.mongoClient.Init(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CollectionComponent) Reset() *CollectionComponent {
	c.apiFeature.Reset()
	return c
}

func (c *CollectionComponent) Close() error {
	if c.svc != nil && c.serviceRunning {
		c.svc.Close(context.Background())
		c.serviceRunning = false
	}
	return nil
}

func (c *CollectionComponent) InitialiseService() (http.Handler, error) {
	var err error
	ctx := context.Background()

	service.GetHTTPServer = c.GetHTTPServer
	service.GetHealthCheck = c.GetHealthCheck
	service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
		return c.mongoClient, nil
	}

	c.svc, err = service.New(ctx, c.config, "1", "", "")
	if err != nil {
		return nil, err
	}

	c.svc.Start(ctx, c.errorChan)
	c.serviceRunning = true
	return c.httpServer.Handler, nil
}

func (c *CollectionComponent) GetHealthCheck(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}
}

func (c *CollectionComponent) GetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.httpServer.Addr = bindAddr
	c.httpServer.Handler = router
	return c.httpServer
}
