package collections

import "github.com/pkg/errors"

// ErrNameSearchTooLong is the error used when an name search value is larger than the maximum allowed
var ErrNameSearchTooLong = errors.New("name search text is >64 chars")

var ErrCollectionNameAlreadyExists = errors.New("a collection with this name already exists")

// QueryParams represents the query parameters that can be sent to get collections
type QueryParams struct {
	Offset     int
	Limit      int
	OrderBy    OrderBy
	NameSearch string
}

// ValidateNameSearchInput returns an error if the given input is not valid as a name search term
func ValidateNameSearchInput(input string) error {
	if len(input) > 64 {
		return ErrNameSearchTooLong
	}
	return nil
}
