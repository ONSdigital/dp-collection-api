package api

import (
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

func (api *API) GetCollectionsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logData := log.Data{}

	offset, limit, err := api.paginator.ReadPaginationParameters(req)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}
	logData["offset"] = offset
	logData["limit"] = limit

	collections, totalCount, err := api.collectionStore.GetCollections(ctx, offset, limit)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}

	response := models.CollectionsResponse{
		Items: collections,
		PaginatedResponse: pagination.PaginatedResponse{
			Count:      len(collections),
			Offset:     offset,
			Limit:      limit,
			TotalCount: totalCount,
		},
	}

	WriteJSONBody(ctx, response, w, logData)
}
