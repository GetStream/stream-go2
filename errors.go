package stream

import (
	"fmt"
)

var (
	errMissingCredentials = fmt.Errorf("missing API key or secret")
)

// APIError is an error returned by Stream API when the request cannot be
// performed or errored server side.
type APIError struct {
	Code            int                      `json:"code,omitempty"`
	Detail          string                   `json:"detail,omitempty"`
	Duration        Duration                 `json:"duration,omitempty"`
	Exception       string                   `json:"exception,omitempty"`
	ExceptionFields map[string][]interface{} `json:"exception_fields,omitempty"`
	StatusCode      int                      `json:"status_code,omitempty"`
}

func (e APIError) Error() string {
	return e.Detail
}

// ToAPIError tries to cast the provided error to APIError type, returning the
// obtained APIError and whether the operation was successful.
func ToAPIError(err error) (APIError, bool) {
	se, ok := err.(APIError)
	return se, ok
}
