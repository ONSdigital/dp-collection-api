package collections

import (
	"errors"
	"strings"
)

// ErrInvalidOrderBy is the error used when an invalid order by value is identified
var ErrInvalidOrderBy = errors.New("invalid order_by")

// OrderBy represents the subset of values that can be used for ordering results
type OrderBy int

const (
	// OrderByDefault is used when no specific order has been specified
	OrderByDefault OrderBy = iota

	// OrderByPublishDate is used to order results by publish date
	OrderByPublishDate
)

// supportedOrderBy defines the supported order by values for ordering collection results.
// the string represents the value as it is in the request query string.
var supportedOrderBy = map[string]OrderBy{
	"publish_date": OrderByPublishDate,
}

// String returns a string representation of the OrderBy instance
func (ob OrderBy) String() string {
	return []string{"default", "publish_date"}[ob]
}

// ParseOrderBy parses the given string as an orderBy value
func ParseOrderBy(orderByInput string) (OrderBy, error) {
	orderByInput = strings.ToLower(orderByInput)
	if len(orderByInput) > 0 {
		orderBy, ok := supportedOrderBy[orderByInput]
		if !ok {
			return 0, ErrInvalidOrderBy
		}

		return orderBy, nil
	}

	return OrderByDefault, nil
}
