package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	URI                   string
	IsSSL                 bool
}

func (m *Mongo) getConnectionConfig() *dpMongoDriver.MongoConnectionConfig {
	return &dpMongoDriver.MongoConnectionConfig{
		ConnectTimeoutInSeconds:       connectTimeoutInSeconds,
		QueryTimeoutInSeconds:         queryTimeoutInSeconds,
		Username:                      m.Username,
		Password:                      m.Password,
		ClusterEndpoint:               m.URI,
		Database:                      m.Database,
		Collection:                    m.CollectionsCollection,
		IsSSL:                         m.IsSSL,
		IsWriteConcernMajorityEnabled: false,
		IsStrongReadConcernEnabled:    false,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init() error {
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

// Checker is called by the health check library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetCollections retrieves all collection documents
func (m *Mongo) GetCollections(ctx context.Context, queryParams collections.QueryParams) ([]models.Collection, int, error) {

	var q *dpMongoDriver.Find
	query := bson.D{}

	if len(queryParams.NameSearch) > 0 {
		query = bson.D{{"name", primitive.Regex{Pattern: queryParams.NameSearch, Options: "i"}}}
	}

	q = m.Connection.
		GetConfiguredCollection().
		Find(query)

	switch queryParams.OrderBy {
	case collections.OrderByPublishDate:
		q.Sort(bson.D{{"publish_date", 1}})
	}

	totalCount, err := q.Count(ctx)
	if err != nil {
		log.Error(ctx, "error getting count of collections from mongo db", err)
		if dpMongoDriver.IsErrNoDocumentFound(err) {
			return []models.Collection{}, totalCount, nil
		}
		return nil, totalCount, err
	}

	var values []models.Collection

	if queryParams.Limit > 0 {
		err = q.Skip(queryParams.Offset).Limit(queryParams.Limit).IterAll(ctx, &values)
		if err != nil {
			if dpMongoDriver.IsErrNoDocumentFound(err) {
				return []models.Collection{}, totalCount, nil
			}
			return nil, totalCount, err
		}
	}

	return values, totalCount, nil
}
