package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

//API provides a struct to wrap the api around
type API struct {
	Router          *mux.Router
	paginator       Paginator
	collectionStore CollectionStore
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, paginator Paginator, collectionStore CollectionStore) *API {
	api := &API{
		Router:          r,
		paginator:       paginator,
		collectionStore: collectionStore,
	}

	r.HandleFunc("/collections", api.PostCollectionHandler).Methods(http.MethodPost)
	r.HandleFunc("/collections", api.GetCollectionsHandler).Methods(http.MethodGet)
	r.HandleFunc("/collections/{collection_id}", api.GetCollectionHandler).Methods(http.MethodGet)
	r.HandleFunc("/collections/{collection_id}", api.PutCollectionHandler).Methods(http.MethodPut)
	r.HandleFunc("/collections/{collection_id}/events", api.GetEventsHandler).Methods(http.MethodGet)
	return api
}

// WriteJSONBody marshals the provided interface into json, and writes it to the response body.
func WriteJSONBody(ctx context.Context, v interface{}, w http.ResponseWriter, data log.Data) error {

	body, err := json.Marshal(v)
	if err != nil {
		handleError(ctx, err, w, data)
		return err
	}

	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}

// NewID returns a new UUID
var NewID = func() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// ValidateUUID returns an error if the string cannot be converted to a valid UUID
func ValidateUUID(s string) error {
	_, err := uuid.FromString(s)
	if err != nil {
		return err
	}
	return nil
}
