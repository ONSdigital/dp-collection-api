package api

import (
	"context"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"net/http"
)

//go:generate moq -out mock/paginator.go -pkg mock . Paginator
//go:generate moq -out mock/collectionstore.go -pkg mock . CollectionStore

// Paginator defines the required methods from the paginator package
type Paginator interface {
	ReadPaginationParameters(r *http.Request) (offset int, limit int, err error)
}

// CollectionStore defines the required methods from the data store of collections
type CollectionStore interface {
	GetCollections(ctx context.Context, offset, limit int, orderBy collections.OrderBy) (collections []models.Collection, totalCount int, err error)
}
