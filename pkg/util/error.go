package util

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common Errors
var (
	ErrNotFound      = NewAppError(http.StatusNotFound, "Resource not found", nil)
	ErrUnauthorized  = NewAppError(http.StatusUnauthorized, "Unauthorized access", nil)
	ErrBadRequest    = NewAppError(http.StatusBadRequest, "Invalid request", nil)
	ErrInternal      = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
	ErrUnverified    = NewAppError(http.StatusForbidden, "Email not verified", nil)
	ErrAlreadyExists = NewAppError(http.StatusConflict, "Resource already exists", nil)
)

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
