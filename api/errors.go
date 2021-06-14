package api

import (
	"context"
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
	}
)

func handleError(ctx context.Context, err error, w http.ResponseWriter, data log.Data) {
	var status int
	switch {

	case badRequest[err]:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	if data == nil {
		data = log.Data{}
	}

	log.Error(ctx, "request unsuccessful", err, data)
	http.Error(w, err.Error(), status)
}
