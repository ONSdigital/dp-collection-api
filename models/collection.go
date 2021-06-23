package models

import (
	"github.com/ONSdigital/dp-collection-api/pagination"
	"time"
)

// Collection represents information related to a single collection
type Collection struct {
	ID          string    `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string    `bson:"name,omitempty"         json:"name,omitempty"`
	PublishDate time.Time `bson:"publish_date,omitempty" json:"publish_date,omitempty"`
	LastUpdated time.Time `bson:"last_updated,omitempty" json:"-"`
}

// CollectionsResponse represents a paginated list of collections
type CollectionsResponse struct {
	Items []Collection `json:"items"`
	pagination.PaginatedResponse
}
