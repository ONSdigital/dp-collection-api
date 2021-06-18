package collections

import "errors"

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

// String returns a string representation of the OrderBy instance
func (ob OrderBy) String() string {
	return []string{"default", "publish_date"}[ob]
}
