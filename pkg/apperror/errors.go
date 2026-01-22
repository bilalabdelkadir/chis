package apperror

import (
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/pkg/shared"
)

type AppError struct {
	Code      int                 `json:"-"`
	Message   string              `json:"message"`
	Path      string              `json:"path"`
	Timestamp time.Time           `json:"timestamp"`
	Cause     error               `json:"-"`
	Details   []shared.FieldError `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func ValidationFailed(details []shared.FieldError) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Validation failed",
		Details: details,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func (e *AppError) WithContext(r *http.Request) *AppError {
	e.Path = r.URL.Path
	e.Timestamp = time.Now().UTC()

	return e
}
