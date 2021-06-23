package api

import (
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/log.go/v2/log"
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
