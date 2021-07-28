package mongo

import (
	"github.com/ONSdigital/dp-collection-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

// AnyETag represents the wildchar that corresponds to not check the ETag value for update requests
const AnyETag = "*"

func newETagForUpdate(currentCollection *models.Collection, update *models.Collection) (eTag string, err error) {
	b, err := bson.Marshal(update)
	if err != nil {
		return "", err
	}
	return currentCollection.Hash(b)
}
