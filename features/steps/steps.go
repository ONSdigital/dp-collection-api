package steps

import (
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/cucumber/godog"
	"github.com/globalsign/mgo/bson"
	"time"
)

func (c *CollectionComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I have these collections:$`, c.iHaveTheseCollections)
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

func (c *CollectionComponent) putDocumentInDatabase(document interface{}, id, collectionName string) error {
	s := c.mongoClient.Session.Copy()
	defer s.Close()

	update := bson.M{
		"$set": document,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := s.DB(c.mongoClient.Database).C(collectionName).UpsertId(id, update)
	if err != nil {
		return err
	}
	return nil
}
