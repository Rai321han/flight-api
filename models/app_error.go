package models

import "fmt"

// AppError carries transport-safe error metadata alongside an optional
// internal cause. It is used across the service and controller layers to
// propagate structured errors that can be mapped directly to HTTP responses
// without leaking internal implementation details to the caller.
type AppError struct {
	Status  int    // HTTP status code to return to the client
	Message string // human-readable message safe to include in the response body
	Cause   error  // underlying error for server-side logging; never sent to the client
}

// NewAppError constructs an AppError with the given HTTP status code,
// client-safe message, and optional underlying cause.
//
// Parameters:
//   - status  : HTTP status code (e.g. 400, 502, 503)
//   - message : human-readable description safe to expose in the response body
//   - cause   : underlying error for internal logging; may be nil
//
// Returns:
//   - *AppError: fully initialised error value
func NewAppError(status int, message string, cause error) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the error interface. It formats the status code and message
// together, appending the cause when one is present.
//
// Returns:
//   - string: formatted error string for logging and debugging purposes;
//     this value is never sent to HTTP clients
func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		return fmt.Sprintf("status=%d message=%s", e.Status, e.Message)
	}
	return fmt.Sprintf("status=%d message=%s: %v", e.Status, e.Message, e.Cause)
}

// Unwrap returns the underlying cause, enabling errors.As and errors.Is to
// traverse the error chain.
//
// Returns:
//   - error: the wrapped cause, or nil when no cause was provided
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}