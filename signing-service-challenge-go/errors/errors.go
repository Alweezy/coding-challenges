package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError defines a structured error with a message and an associated HTTP status code.
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Predefined errors that can be reused across the app.
var (
	ErrNotFound       = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInvalidRequest = &AppError{Code: http.StatusBadRequest, Message: "Invalid request data"}
	ErrInternal       = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
)

// WrapError wraps an existing error with a custom message and HTTP status code.
func WrapError(err error, message string, code int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// FromError converts a generic error to an AppError, returning an internal server error if none is found.
func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return ErrInternal
}
