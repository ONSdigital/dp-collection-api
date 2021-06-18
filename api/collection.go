package api

import (
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
	"strings"
)

// collectionsOrderBy defines the supported order by values for ordering collection results.
// the string represents the value as it is in the request query string.
var collectionsOrderBy = map[string]collections.OrderBy{
	"publish_date": collections.OrderByPublishDate,
}

// GetCollectionsHandler handles HTTP requests for the get collections endpoint
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

	orderBy, err := parseCollectionsOrderBy(req)
	if err != nil {
		handleError(ctx, err, w, logData)
		return
	}
	logData["order_by"] = orderBy

	collections, totalCount, err := api.collectionStore.GetCollections(ctx, offset, limit, orderBy)
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	WriteJSONBody(ctx, response, w, logData)
}

func parseCollectionsOrderBy(req *http.Request) (collections.OrderBy, error) {

	orderByInput := strings.ToLower(req.URL.Query().Get("order_by"))
	if len(orderByInput) > 0 {
		orderBy, ok := collectionsOrderBy[orderByInput]
		if !ok {
			return 0, collections.ErrInvalidOrderBy
		}

		return orderBy, nil
	}

	return collections.OrderByDefault, nil
}
