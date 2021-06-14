package pagination

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrInvalidOffsetParameter = errors.New("invalid offset query parameter")
	ErrInvalidLimitParameter  = errors.New("invalid limit query parameter")
	ErrLimitOverMax           = errors.New("limit query parameter is larger than the maximum allowed")
)

type Paginator struct {
	DefaultLimit    int
	DefaultOffset   int
	DefaultMaxLimit int
}

type PaginatedResponse struct {
	Count      int `json:"count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
}

func NewPaginator(defaultLimit, defaultOffset, defaultMaxLimit int) *Paginator {
	return &Paginator{
		DefaultLimit:    defaultLimit,
		DefaultOffset:   defaultOffset,
		DefaultMaxLimit: defaultMaxLimit,
	}
}

func (p *Paginator) ReadPaginationParameters(r *http.Request) (offset int, limit int, err error) {

	offsetParameter := r.URL.Query().Get("offset")
	limitParameter := r.URL.Query().Get("limit")

	offset = p.DefaultOffset
	limit = p.DefaultLimit

	if offsetParameter != "" {
		offset, err = strconv.Atoi(offsetParameter)
		if err != nil || offset < 0 {
			return 0, 0, ErrInvalidOffsetParameter
		}
	}

	if limitParameter != "" {
		limit, err = strconv.Atoi(limitParameter)
		if err != nil || limit < 0 {
			return 0, 0, ErrInvalidLimitParameter
		}
	}

	if limit > p.DefaultMaxLimit {
		return 0, 0, ErrLimitOverMax
	}

	return
}
