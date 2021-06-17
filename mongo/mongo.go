package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/globalsign/mgo"
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	healthClient          *dpMongoHealth.CheckMongoClient
	Database              string
	CollectionsCollection string
	Session               *mgo.Session
	URI                   string
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context) error {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	var err error
	if m.Session, err = mgo.Dial(m.URI); err != nil {
		return err
	}
	m.Session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	m.Session.SetMode(mgo.Strong, true)

	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.CollectionsCollection)}

	// Create client and health client from session AND collections
	client := dpMongoHealth.NewClientWithCollections(m.Session, databaseCollectionBuilder)

	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Session == nil {
		return errors.New("cannot close a mongoDB connection without a valid session")
	}
	return dpMongoDriver.Close(ctx, m.Session)
}

// Checker is called by the health check library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetCollections retrieves all collection documents
func (m *Mongo) GetCollections(ctx context.Context, offset, limit int, orderBy collections.OrderBy) ([]models.Collection, int, error) {

	s := m.Session.Copy()
	defer s.Close()

	var q *mgo.Query

	q = s.DB(m.Database).C(m.CollectionsCollection).Find(nil)

	switch orderBy {
	case collections.OrderByPublishDate:
		q.Sort("publish_date")
	}

	totalCount, err := q.Count()
	if err != nil {
		log.Error(ctx, "error getting count of collections from mongo db", err)
		if err == mgo.ErrNotFound {
			return []models.Collection{}, totalCount, nil
		}
		return nil, totalCount, err
	}

	var values []models.Collection

	if limit > 0 {
		iter := q.Skip(offset).Limit(limit).Iter()

		defer func() {
			err := iter.Close()
			if err != nil {
				log.Error(ctx, "error closing iterator", err)
			}
		}()

		if err := iter.All(&values); err != nil {
			if err == mgo.ErrNotFound {
				return []models.Collection{}, totalCount, nil
			}
			return nil, totalCount, err
		}
	}

	return values, totalCount, nil
}
