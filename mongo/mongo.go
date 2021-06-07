package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongo-driver"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	healthClient          *dpMongoHealth.CheckMongoClient
	Database              string
	CollectionsCollection string
	Connection            *dpMongoDriver.MongoConnection
	Username              string
	Password              string
	CAFilePath            string
	URI                   string
}

func (m *Mongo) getConnectionConfig() *dpMongoDriver.MongoConnectionConfig {
	return &dpMongoDriver.MongoConnectionConfig{
		CaFilePath:              m.CAFilePath,
		ConnectTimeoutInSeconds: connectTimeoutInSeconds,
		QueryTimeoutInSeconds:   queryTimeoutInSeconds,

		Username:             m.Username,
		Password:             m.Password,
		ClusterEndpoint:      m.URI,
		Database:             m.Database,
		Collection:           m.CollectionsCollection,
		SkipCertVerification: true,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context) error {
	if m.Connection != nil {
		return errors.New("datastore connection already exists")
	}
	mongoConnection, err := dpMongoDriver.Open(m.getConnectionConfig())
	if err != nil {
		return err
	}
	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.CollectionsCollection)}

	// Create client and health client from session AND collections
	client := dpMongoHealth.NewClientWithCollections(mongoConnection, databaseCollectionBuilder)

	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Connection == nil {
		return errors.New("cannot close a empty connection")
	}
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}
