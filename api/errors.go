package api

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

var (
	// errors that should return a 400 status
	badRequest = map[error]bool{
		pagination.ErrInvalidLimitParameter:  true,
		pagination.ErrInvalidOffsetParameter: true,
		pagination.ErrLimitOverMax:           true,
		collections.ErrInvalidOrderBy:        true,
		collections.ErrNameSearchTooLong:     true,
		collections.ErrCollectionNameEmpty:   true,
		collections.ErrInvalidID:             true,
		collections.ErrNoIfMatchHeader:       true,
		ErrUnableToParseJSON:                 true,
	}

	notFound = map[error]bool{
		collections.ErrCollectionNotFound: true,
	}

	conflictRequest = map[error]bool{
		collections.ErrCollectionNameAlreadyExists: true,
		collections.ErrCollectionConflict:          true,
	}

	ErrUnableToParseJSON = errors.New("failed to parse json body")
)

func handleError(ctx context.Context, err error, w http.ResponseWriter, logData log.Data) {
	var status int
	switch {

	case badRequest[err]:
		status = http.StatusBadRequest
	case notFound[err]:
		status = http.StatusNotFound
	case conflictRequest[err]:
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	if logData == nil {
		logData = log.Data{}
	}

	response := models.ErrorsResponse{
		Errors: []models.ErrorResponse{
			{
				Message: err.Error(),
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	WriteJSONBody(ctx, response, w, logData)
	log.Error(ctx, "request unsuccessful", err, logData)
}
