package models

import "fmt"

// AppError carries transport-safe error metadata plus an internal cause.
type AppError struct {
	Status  int
	Message string
	Cause   error
}

func NewAppError(status int, message string, cause error) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
		Cause:   cause,
	}
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		return fmt.Sprintf("status=%d message=%s", e.Status, e.Message)
	}
	return fmt.Sprintf("status=%d message=%s: %v", e.Status, e.Message, e.Cause)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
