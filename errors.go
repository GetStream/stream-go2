package stream

import (
	"errors"
)

var (
	errMissingCredentials = errors.New("missing API key or secret")
	errInvalidUserID      = errors.New("invalid userID provided")
	errToTargetsNoChanges = errors.New("no changes specified, please supply new targets or added/removed targets")
)

// Rate limit headers
const (
	HeaderRateLimit     = "X-Ratelimit-Limit"
	HeaderRateRemaining = "X-Ratelimit-Remaining"
	HeaderRateReset     = "X-Ratelimit-Reset"
)

// APIError is an error returned by Stream API when the request cannot be
// performed or errored server side.
type APIError struct {
	Code            int              `json:"code,omitempty"`
	Detail          string           `json:"detail,omitempty"`
	Duration        Duration         `json:"duration,omitempty"`
	Exception       string           `json:"exception,omitempty"`
	ExceptionFields map[string][]any `json:"exception_fields,omitempty"`
	StatusCode      int              `json:"status_code,omitempty"`
	Rate            *Rate            `json:"-"`
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
