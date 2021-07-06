package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/v2/log"
	"io"
	"io/ioutil"
	"net/http"
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

func (api *API) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {

	defer dphttp.DrainBody(r)
	ctx := r.Context()
	logData := log.Data{}

	collection, err := ParseCollection(ctx, r.Body)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	err = api.validateCollection(ctx, collection)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	if err = api.collectionStore.UpsertCollection(ctx, collection); err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = WriteJSONBody(ctx, collection, w, logData)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	log.Event(ctx, "add collection request completed successfully", log.INFO, logData)
}

func (api *API) validateCollection(ctx context.Context, collection *models.Collection) error {

	if len(collection.Name) == 0 {
		return collections.ErrCollectionNameEmpty
	}

	_, err := api.collectionStore.GetCollectionByName(ctx, collection.Name)
	if err != nil && !mongodb.IsErrNoDocumentFound(err) {
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

	collection.ID, err = NewID()
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
