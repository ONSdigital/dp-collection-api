package models

// ErrorsResponse represents a slice of errors in a JSON response body
type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

// ErrorResponse represents a single error in a JSON response body
type ErrorResponse struct {
	Message string `json:"message"`
}
