package models

import (
	"crypto/sha1"
	"fmt"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/globalsign/mgo/bson"
	"time"
)

// AnyETag represents the wildcard that corresponds to not check the ETag value for update requests
const AnyETag = "*"

// Collection represents information related to a single collection
type Collection struct {
	ID          string     `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string     `bson:"name,omitempty"         json:"name,omitempty"`
	PublishDate *time.Time `bson:"publish_date,omitempty" json:"publish_date,omitempty"`
	LastUpdated time.Time  `bson:"last_updated,omitempty" json:"-"`
	ETag        string     `bson:"e_tag"                  json:"e_tag,omitempty"`
}

// CollectionsResponse represents a paginated list of collections
type CollectionsResponse struct {
	Items []Collection `json:"items"`
	pagination.PaginatedResponse
}

// Hash generates a SHA-1 hash of the collection struct. SHA-1 is not cryptographically safe,
// but it has been selected for performance as we are only interested in uniqueness.
// ETag field value is ignored when generating a hash.
// An optional byte array can be provided to append to the hash.
// This can be used, for example, to calculate a hash of this filter and an update applied to it.
func (c *Collection) Hash(extraBytes []byte) (string, error) {
	h := sha1.New()

	// copy by value to ignore ETag without affecting f
	c2 := *c
	c2.ETag = ""

	collectionBytes, err := bson.Marshal(c2)
	if err != nil {
		return "", err
	}

	if _, err := h.Write(append(collectionBytes, extraBytes...)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (c *Collection) NewETagForUpdate(update *Collection) (eTag string, err error) {
	b, err := bson.Marshal(update)
	if err != nil {
		return "", err
	}
	return c.Hash(b)
}
