package api

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetCollectionsHandler handles HTTP requests for the get collections endpoint
func (api *API) GetCollectionsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logData := log.Data{}

	queryParams, err := readCollectionsQueryParams(req, api.paginator)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}
	logData["query_params"] = queryParams

	collections, totalCount, err := api.collectionStore.GetCollections(ctx, *queryParams)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	response := models.CollectionsResponse{
		Items: collections,
		PaginatedResponse: pagination.PaginatedResponse{
			Count:      len(collections),
			Offset:     queryParams.Offset,
			Limit:      queryParams.Limit,
			TotalCount: totalCount,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	WriteJSONBody(ctx, response, w, logData)
}

func (api *API) GetCollectionHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logData := log.Data{}
	eTag := getIfMatch(req)
	logData["e_tag"] = eTag

	collectionID := mux.Vars(req)["collection_id"]
	logData["collection_id"] = collectionID

	err := ValidateUUID(collectionID)
	if err != nil {
		handleError(ctx, collections.ErrInvalidID, w, logData)
		return
	}

	collection, err := api.collectionStore.GetCollectionByID(ctx, collectionID, eTag)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	setETag(w, collection.ETag)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	WriteJSONBody(ctx, collection, w, logData)
}

func (api *API) PostCollectionHandler(w http.ResponseWriter, r *http.Request) {

	defer dphttp.DrainBody(r)
	ctx := r.Context()
	logData := log.Data{}

	collection, err := ParseCollection(ctx, r.Body)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	collection.ID, err = NewID()
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	err = api.validateCollection(ctx, collection)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	if err = api.collectionStore.AddCollection(ctx, collection); err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	setETag(w, collection.ETag)
	w.WriteHeader(http.StatusCreated)
	err = WriteJSONBody(ctx, collection, w, logData)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	log.Info(ctx, "add collection request completed successfully", logData)
}

func (api *API) PutCollectionHandler(w http.ResponseWriter, req *http.Request) {
	defer dphttp.DrainBody(req)

	ctx := req.Context()
	logData := log.Data{}
	collectionID := mux.Vars(req)["collection_id"]
	logData["collection_id"] = collectionID

	log.Info(ctx, "put collection", logData)

	err := ValidateUUID(collectionID)
	if err != nil {
		handleError(ctx, collections.ErrInvalidID, w, logData)
		return
	}

	_, err = api.collectionStore.GetCollectionByID(ctx, collectionID, models.AnyETag)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	// eTag value must be present in If-Match header
	eTag, err := getIfMatchForce(req)
	if err != nil {
		log.Error(ctx, "missing header", err, log.Data{"error": err.Error()})
		handleError(ctx, err, w, logData)
		return
	}
	logData["e_tag"] = eTag

	collection, err := ParseCollection(ctx, req.Body)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	collection.ID = collectionID

	if err := api.collectionStore.ReplaceCollection(ctx, collection, eTag); err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	setETag(w, collection.ETag)
	w.WriteHeader(http.StatusOK)
	err = WriteJSONBody(ctx, collection, w, logData)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	log.Info(ctx, "put collection request completed successfully", logData)
}

func (api *API) validateCollection(ctx context.Context, collection *models.Collection) error {

	if collection == nil {
		return collections.ErrNilCollection
	}

	if len(collection.Name) == 0 {
		return collections.ErrCollectionNameEmpty
	}

	_, err := api.collectionStore.GetCollectionByName(ctx, collection.Name)
	if err != nil && err != collections.ErrCollectionNotFound {
		return err
	}
	if err == nil {
		return collections.ErrCollectionNameAlreadyExists
	}

	return nil
}

func ParseCollection(ctx context.Context, reader io.Reader) (*models.Collection, error) {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var collection models.Collection

	err = json.Unmarshal(b, &collection)
	if err != nil {
		log.Error(ctx, "failed to parse collection json body", err)
		return nil, ErrUnableToParseJSON
	}

	// set eTag value to current hash of the collection
	collection.ETag, err = collection.Hash(nil)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func readCollectionsQueryParams(req *http.Request, paginator Paginator) (*collections.QueryParams, error) {

	offset, limit, err := paginator.ReadPaginationParameters(req)
	if err != nil {
		return nil, err
	}

	orderByInput := req.URL.Query().Get("order_by")
	orderBy, err := collections.ParseOrderBy(orderByInput)
	if err != nil {
		return nil, err
	}

	nameSearchInput := req.URL.Query().Get("name")
	err = collections.ValidateNameSearchInput(nameSearchInput)
	if err != nil {
		return nil, err
	}

	return &collections.QueryParams{
		Offset:     offset,
		Limit:      limit,
		OrderBy:    orderBy,
		NameSearch: nameSearchInput,
	}, nil
}

func getIfMatch(r *http.Request) string {
	ifMatch := r.Header.Get("If-Match")
	if ifMatch == "" {
		return models.AnyETag
	}
	return ifMatch
}

func setETag(w http.ResponseWriter, eTag string) {
	w.Header().Set("ETag", eTag)
}

func getIfMatchForce(r *http.Request) (string, error) {
	eTag := getIfMatch(r)
	if eTag == models.AnyETag {
		err := collections.ErrNoIfMatchHeader
		return "", err
	}
	return eTag, nil
}
