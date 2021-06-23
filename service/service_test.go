package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-collection-api/config"
	"github.com/ONSdigital/dp-collection-api/service"
	"github.com/ONSdigital/dp-collection-api/service/mock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx            = context.Background()
	testBuildTime  = "1622727733"
	testGitCommit  = "GitCommit"
	testVersion    = "Version"
	errServer      = errors.New("HTTP Server error")
	errAddCheck    = errors.New("health check add check error")
	errHealthCheck = errors.New("healthCheck error")
)

func TestNew(t *testing.T) {

	Convey("Having a set of mocked dependencies", t, func() {

		cfg := &config.Config{}

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		}
		service.GetHealthCheck = func(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
			return hcMock
		}

		serverMock := &mock.HTTPServerMock{}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		mongoDBMock := &mock.MongoDBMock{}
		service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
			return mongoDBMock, nil
		}

		Convey("Given that health check versionInfo cannot be created due to a wrong build time", func() {
			wrongBuildTime := "wrongFormat"

			Convey("When service.New is called", func() {
				_, err := service.New(ctx, cfg, wrongBuildTime, testGitCommit, testVersion)

				Convey("Then the expected error is returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldResemble, "failed to parse build time")
				})

				Convey("Then the expected health checks are registered", func() {
					So(len(hcMock.AddCheckCalls()), ShouldEqual, 0)
				})
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			Convey("When service.New is called", func() {
				svc, err := service.New(ctx, cfg, testBuildTime, testGitCommit, testVersion)
				So(err, ShouldBeNil)
				So(svc, ShouldNotBeNil)

				Convey("Then the expected health checks are registered", func() {
					So(len(hcMock.AddCheckCalls()), ShouldEqual, 1)
					So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "Mongo DB")
				})
			})
		})
	})

	Convey("Given an error occurs when MongoDB is initialised", t, func() {

		cfg := &config.Config{}

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		}
		service.GetHealthCheck = func(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
			return hcMock
		}

		serverMock := &mock.HTTPServerMock{}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		expectedErr := errors.New("mongodb failed to initialised")
		service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
			return nil, expectedErr
		}

		Convey("When service.New is called", func() {
			svc, err := service.New(ctx, cfg, testBuildTime, testGitCommit, testVersion)
			So(svc, ShouldBeNil)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedErr)
			})
		})
	})
}

func TestStart(t *testing.T) {

	Convey("Having a correctly initialised Service with mocked dependencies", t, func() {

		cfg := &config.Config{}

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}
		service.GetHealthCheck = func(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
			return hcMock
		}

		serverMock := &mock.HTTPServerMock{}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		mongoDBMock := &mock.MongoDBMock{}
		service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
			return mongoDBMock, nil
		}

		serverWg := &sync.WaitGroup{}

		svc, err := service.New(ctx, cfg, testBuildTime, testGitCommit, testVersion)
		So(err, ShouldBeNil)
		So(svc, ShouldNotBeNil)

		Convey("When a service with a successful HTTP server is started", func() {
			serverMock.ListenAndServeFunc = func() error {
				serverWg.Done()
				return nil
			}
			serverWg.Add(1)
			svc.Start(ctx, make(chan error, 1))

			Convey("Then health check is started and HTTP server starts listening", func() {
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})

		Convey("When a service with a failing HTTP server is started", func() {
			serverMock.ListenAndServeFunc = func() error {
				serverWg.Done()
				return errServer
			}
			errChan := make(chan error, 1)
			serverWg.Add(1)
			svc.Start(ctx, errChan)

			Convey("Then HTTP server errors are reported to the provided errors channel", func() {
				rxErr := <-errChan
				So(rxErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
			})
		})
	})
}

func TestClose(t *testing.T) {

	Convey("Having a correctly initialised service with mocked dependencies", t, func() {
		cfg := &config.Config{
			GracefulShutdownTimeout: 5 * time.Second,
		}
		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			StopFunc: func() { hcStopped = true },
			AddCheckFunc: func(name string, checker healthcheck.Checker) error {
				return nil
			},
		}
		service.GetHealthCheck = func(version healthcheck.VersionInfo, criticalTimeout, interval time.Duration) service.HealthChecker {
			return hcMock
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server was stopped before healthcheck")
				}
				return nil
			},
		}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		mongoDBMock := &mock.MongoDBMock{
			CloseFunc: func(ctx context.Context) error {
				return nil
			},
		}
		service.GetMongoDB = func(ctx context.Context, cfg config.MongoConfig) (service.MongoDB, error) {
			return mongoDBMock, nil
		}

		svc, err := service.New(ctx, cfg, testBuildTime, testGitCommit, testVersion)
		So(err, ShouldBeNil)
		So(svc, ShouldNotBeNil)

		Convey("Given that all dependencies succeed to close", func() {
			Convey("Closing the service results in all the initialised dependencies being closed in the expected order", func() {
				err = svc.Close(context.Background())
				So(err, ShouldBeNil)
				So(len(hcMock.StopCalls()), ShouldEqual, 1)
				So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
				So(len(mongoDBMock.CloseCalls()), ShouldEqual, 1)
			})
		})

		Convey("Given that all dependencies fail to close", func() {
			serverMock.ShutdownFunc = func(ctx context.Context) error {
				return errServer
			}

			Convey("Then closing the service fails with the expected error and further dependencies are attempted to close", func() {
				err = svc.Close(context.Background())
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, "failed to shutdown gracefully")
				So(len(hcMock.StopCalls()), ShouldEqual, 1)
				So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
			})
		})

		Convey("Given that a dependency takes more time to close than the graceful shutdown timeout", func() {
			cfg.GracefulShutdownTimeout = 5 * time.Millisecond
			serverMock.ShutdownFunc = func(ctx context.Context) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			}

			Convey("Then closing the service fails with context.DeadlineExceeded error and no further dependencies are attempted to close", func() {
				err = svc.Close(context.Background())
				So(err, ShouldResemble, context.DeadlineExceeded)
				So(len(hcMock.StopCalls()), ShouldEqual, 1)
				So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
			})
		})
	})
}
