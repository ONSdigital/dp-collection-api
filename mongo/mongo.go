package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	mongo "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/globalsign/mgo"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
	EventsCollection      string
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
		C(m.CollectionsCollection).
		Find(query)

	switch queryParams.OrderBy {
	case collections.OrderByPublishDate:
		q.Sort(bson.D{{"publish_date", 1}})
	}

	totalCount, err := q.Count(ctx)
	if err != nil {
		log.Error(ctx, "error getting count of collections from mongo db", err)
		return nil, totalCount, err
	}

	values := []models.Collection{}

	if queryParams.Limit > 0 {
		err = q.Skip(queryParams.Offset).Limit(queryParams.Limit).IterAll(ctx, &values)
		if err != nil {
			return nil, totalCount, err
		}
	}

	return values, totalCount, nil
}

// GetCollectionByName retrieves a single collection by name
func (m *Mongo) GetCollectionByName(ctx context.Context, name string) (*models.Collection, error) {

	query := bson.D{{"name", name}}
	result := &models.Collection{}

	err := m.Connection.
		C(m.CollectionsCollection).
		FindOne(ctx, query, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetCollectionByID retrieves a single collection by ID
func (m *Mongo) GetCollectionByID(ctx context.Context, id string, eTagSelector string) (*models.Collection, error) {

	query := bson.D{{"_id", id}}
	result := &models.Collection{}

	err := m.Connection.
		C(m.CollectionsCollection).
		FindOne(ctx, query, result)
	if err != nil {
		return nil, err
	}

	// If eTag was provided and did not match, return the corresponding error
	if eTagSelector != AnyETag && eTagSelector != result.ETag {
		return nil, errors.New("conflict when retrieving the Collection")
	}

	return result, nil
}

// AddCollection adds a collection
func (m *Mongo) AddCollection(ctx context.Context, collection *models.Collection) error {
	var err error
	// set eTag value to current hash of the collection
	collection.ETag, err = collection.Hash(nil)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": collection,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}

	_, err = m.Connection.C(m.CollectionsCollection).UpsertId(ctx, collection.ID, update)

	return err
}

// UpdateCollection updates a collection
func (m *Mongo) UpdateCollection(ctx context.Context, updatedCollection *models.Collection, eTagSelector string, currentCollection *models.Collection) (newETag string, err error) {
	// calculate the new eTag hash for the filter that would result from applying the update
	newETag, err = newETagForUpdate(currentCollection, updatedCollection)
	if err != nil {
		return "", err
	}

	selector := bson.M{
		"_id":   updatedCollection.ID,
		"e_tag": eTagSelector,
	}

	update, err := mongo.WithUpdates(bson2.M(bson.M{
		"$set": bson.M{
			"e_tag":        newETag,
		},
	}))

	if _, err := m.Connection.C(m.CollectionsCollection).Update(ctx, selector, update); err != nil {
		if err == mgo.ErrNotFound {
			return "", collections.ErrCollectionConflict
		}
		return "", err
	}

	return newETag, nil
}

// GetCollectionEvents retrieves all events for a collection
func (m *Mongo) GetCollectionEvents(ctx context.Context, queryParams collections.EventsQueryParams) ([]models.Event, int, error) {

	var q *dpMongoDriver.Find

	query := bson.D{{"collection_id", queryParams.CollectionID}}

	q = m.Connection.
		C(m.EventsCollection).
		Find(query).
		Sort(bson.D{{"date", 1}})

	totalCount, err := q.Count(ctx)
	if err != nil {
		log.Error(ctx, "error getting count of collection events from mongo db", err)
		return nil, totalCount, err
	}

	values := []models.Event{}

	if queryParams.Limit > 0 {
		err = q.Skip(queryParams.Offset).Limit(queryParams.Limit).IterAll(ctx, &values)
		if err != nil {
			return nil, totalCount, err
		}
	}

	return values, totalCount, nil
}
