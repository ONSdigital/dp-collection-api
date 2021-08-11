package collections

import "github.com/pkg/errors"

// ErrNameSearchTooLong is the error used when an name search value is larger than the maximum allowed
var ErrNameSearchTooLong = errors.New("name search text is >64 chars")

// ErrCollectionNameAlreadyExists is the error used when an existing collection is already using the collection name
var ErrCollectionNameAlreadyExists = errors.New("a collection with this name already exists")

// ErrCollectionNameEmpty is the error used when an empty collection name is provided
var ErrCollectionNameEmpty = errors.New("the collection name field must be specified")

// ErrCollectionIDEmpty is the error used when an empty collection ID is provided
var ErrCollectionIDEmpty = errors.New("the collection id field must be specified")

// ErrNilCollection is the error used when nil collection is identified
var ErrNilCollection = errors.New("could not validate a nil collection")

// ErrCollectionNotFound is the error used when a particular collection is not found
var ErrCollectionNotFound = errors.New("collection not found")

// ErrInvalidID is the error used when an invalid collection ID format is used
var ErrInvalidID = errors.New("collection id must be valid UUID")

// ErrCollectionConflict is the error used when the provided eTag does not match
var ErrCollectionConflict = errors.New("out of date collection resource")

// ErrNoIfMatchHeader is the error used when an If-Match is required but not provided
var ErrNoIfMatchHeader = errors.New("required If-Match header not provided")

// QueryParams represents the query parameters that can be sent to get collections
type QueryParams struct {
	Offset     int
	Limit      int
	OrderBy    OrderBy
	NameSearch string
}

// EventsQueryParams represents the parameters to query a collection's events
type EventsQueryParams struct {
	CollectionID string
	Offset       int
	Limit        int
}

// ValidateNameSearchInput returns an error if the given input is not valid as a name search term
func ValidateNameSearchInput(input string) error {
	if len(input) > 64 {
		return ErrNameSearchTooLong
	}
	return nil
}
