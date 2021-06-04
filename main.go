package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ONSdigital/dp-collection-api/config"
	"github.com/ONSdigital/dp-collection-api/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const serviceName = "dp-collection-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(nil, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)

	log.Info(ctx, "dp-collection-api version", log.Data{"version": Version})

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}
	log.Event(ctx, "loaded config", log.INFO, log.Data{"config": cfg})

	service, err := service.New(ctx, cfg, BuildTime, GitCommit, Version)
	if err != nil {
		return errors.Wrap(err, "error creating service")
	}

	// Start service
	service.Start(ctx, svcErrors)
	if err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return service.Close(ctx)
}
