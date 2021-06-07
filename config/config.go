package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config represents service configuration for dp-collection-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	MongoConfig                MongoConfig
}

// MongoConfig contains the config required to connect to MongoDB.
type MongoConfig struct {
	BindAddr              string `envconfig:"MONGODB_BIND_ADDR"           json:"-"` // This line contains sensitive data and the json:"-" tells the json marshaller to skip serialising it.
	CollectionsDatabase   string `envconfig:"MONGODB_COLLECTIONS_DATABASE"`
	CollectionsCollection string `envconfig:"MONGODB_COLLECTIONS_COLLECTION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:26000",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MongoConfig: MongoConfig{
			BindAddr:              "localhost:27017",
			CollectionsDatabase:   "collections",
			CollectionsCollection: "collections",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
