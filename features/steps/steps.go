package steps

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *CollectionComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^there are no collections`, c.thereAreNoCollections)
	ctx.Step(`^I have these collections:$`, c.iHaveTheseCollections)
}

func (c *CollectionComponent) iHaveTheseCollections(input *godog.DocString) error {
	var collections []models.Collection

	err := json.Unmarshal([]byte(input.Content), &collections)
	if err != nil {
		return err
	}

	for _, collection := range collections {
		if err := c.putDocumentInDatabase(collection, collection.ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectionComponent) thereAreNoCollections() error {
	// nothing to do
	return nil
}

func (c *CollectionComponent) putDocumentInDatabase(document interface{}, id string) error {

	update := bson.M{
		"$set": document,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := c.mongoClient.Connection.GetConfiguredCollection().UpsertId(context.Background(), id, update)
	if err != nil {
		return err
	}
	return nil
}
