package security

import (
	"fmt"
)

// ClientError is an extension of the error interface that provides a
// `SafeError` method. This should return a string representation of the error
// that is safe to show to clients. This allows us to properly log the cause of
// the error whilst migitigating potential security breaches.
type ClientError interface {
	error
	SafeError() string
}

type clientError struct {
	Message string
	Cause error
}

// NewClientError creates a new client error.
func NewClientError(message string, err error) ClientError {
	return &clientError{
		Message: message,
		Cause: err,
	}
}

// Error returns clientError as a string. This is suitable for logging but not
// returning to the client.
func (e *clientError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

// SafeError returns a representation of client error that is safe to return to
// the client.
func (e *clientError) SafeError() string {
	return e.Message
}
