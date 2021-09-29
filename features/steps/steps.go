package steps

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *CollectionComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^there are no collections`, c.thereAreNoCollections)
	ctx.Step(`^I have these collections:$`, c.iHaveTheseCollections)
	ctx.Step(`^I have a collection with ID "([^"]*)" with the following events:$`, c.iHaveCollectionWithEvents)
}

func (c *CollectionComponent) iHaveCollectionWithEvents(collectionID string, documentJson *godog.DocString) error {

	collection := models.Collection{
		ID:          collectionID,
		Name:        "collection name",
		PublishDate: &time.Time{},
		LastUpdated: time.Time{},
	}

	if err := c.putDocumentInDatabase(collection, collection.ID, c.config.MongoConfig.CollectionsCollection); err != nil {
		return err
	}

	var events []models.Event
	json.Unmarshal([]byte(documentJson.Content), &events)

	for _, event := range events {

		event.CollectionID = collectionID
		uuid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		if err := c.putDocumentInDatabase(event, uuid.String(), c.config.MongoConfig.EventsCollection); err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectionComponent) iHaveTheseCollections(input *godog.DocString) error {
	var collections []models.Collection

	err := json.Unmarshal([]byte(input.Content), &collections)
	if err != nil {
		return err
	}

	for _, collection := range collections {
		if err := c.putDocumentInDatabase(collection, collection.ID, c.config.MongoConfig.CollectionsCollection); err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectionComponent) thereAreNoCollections() error {
	// nothing to do
	return nil
}

func (c *CollectionComponent) putDocumentInDatabase(document interface{}, id, collectionName string) error {

	update := bson.M{
		"$set": document,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := c.mongoClient.Connection.C(collectionName).UpsertById(context.Background(), id, update)
	if err != nil {
		return err
	}
	return nil
}
