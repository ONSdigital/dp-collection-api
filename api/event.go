package api

import (
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"net/http"
)

// GetEventsHandler handles HTTP requests for the get collection events endpoint
func (api *API) GetEventsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logData := log.Data{}

	queryParams, err := readEventsQueryParams(req, api.paginator)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}
	logData["query_params"] = queryParams

	// todo check collection exists

	events, totalCount, err := api.collectionStore.GetCollectionEvents(ctx, *queryParams)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	response := models.EventsResponse{
		Items: events,
		PaginatedResponse: pagination.PaginatedResponse{
			Count:      len(events),
			Offset:     queryParams.Offset,
			Limit:      queryParams.Limit,
			TotalCount: totalCount,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	WriteJSONBody(ctx, response, w, logData)
}

func readEventsQueryParams(req *http.Request, paginator Paginator) (*collections.EventsQueryParams, error) {

	offset, limit, err := paginator.ReadPaginationParameters(req)
	if err != nil {
		return nil, err
	}

	vars := mux.Vars(req)
	collectionID := vars["collection_id"]

	if len(collectionID) == 0 {
		return nil, collections.ErrCollectionIDEmpty
	}

	return &collections.EventsQueryParams{
		Offset:       offset,
		Limit:        limit,
		CollectionID: collectionID,
	}, nil
}
