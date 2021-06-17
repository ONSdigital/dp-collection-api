package collections

import "errors"

var ErrInvalidOrderBy = errors.New("invalid order_by")

type OrderBy int

const (
	OrderByDefault OrderBy = iota
	OrderByPublishDate
)

func (ob OrderBy) String() string {
	return []string{"default", "publish_date"}[ob]
}
