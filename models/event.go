package models

import (
	"github.com/ONSdigital/dp-collection-api/pagination"
	"time"
)

// Event represents the data for a single collection event
type Event struct {
	ID           string    `bson:"_id,omitempty"   json:"-"`
	Type         string    `bson:"type,omitempty"  json:"type,omitempty"`
	Email        string    `bson:"email,omitempty" json:"email,omitempty"`
	Date         time.Time `bson:"date,omitempty"  json:"date,omitempty"`
	CollectionID string    `bson:"collection_id,omitempty"   json:"-"`
}

// EventsResponse represents a paginated list of collection events
type EventsResponse struct {
	Items []Event `json:"items"`
	pagination.PaginatedResponse
}
