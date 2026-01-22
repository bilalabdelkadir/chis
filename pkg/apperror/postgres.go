package apperror

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	// Integrity / constraints
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgNotNullViolation    = "23502"
	pgCheckViolation      = "23514"

	// Data issues
	pgInvalidTextFormat         = "22P02"
	pgNumericOutOfRange         = "22003"
	pgStringDataRightTruncation = "22001"

	// Transaction / locking
	pgSerializationFailure = "40001"
	pgDeadlockDetected     = "40P01"

	// Auth / access
	pgInsufficientPrivilege = "42501"
	pgInvalidAuthorization  = "28000"
)

// newAppError creates an AppError without timestamp.
// Timestamp will be set by WithContext()/Wrap.
func newAppError(code int, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// FromPostgres converts a Postgres error into an AppError.
// Returns nil if the error is not a Postgres error.
func FromPostgres(err error) *AppError {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {

	// ---------- 409 Conflict ----------
	case pgUniqueViolation:
		return newAppError(http.StatusConflict, "Resource already exists.", err)

	// ---------- 400 Bad Request ----------
	case pgForeignKeyViolation:
		return newAppError(http.StatusBadRequest, "Invalid reference.", err)

	case pgNotNullViolation:
		return newAppError(http.StatusBadRequest, "Required field is missing.", err)

	case pgCheckViolation:
		return newAppError(http.StatusBadRequest, "Data validation failed.", err)

	case pgInvalidTextFormat:
		return newAppError(http.StatusBadRequest, "Invalid input format.", err)

	case pgNumericOutOfRange:
		return newAppError(http.StatusBadRequest, "Numeric value out of range.", err)

	case pgStringDataRightTruncation:
		return newAppError(http.StatusBadRequest, "Input exceeds maximum length.", err)

	// ---------- 503 Service Unavailable ----------
	case pgSerializationFailure, pgDeadlockDetected:
		return newAppError(http.StatusServiceUnavailable, "Database is busy, please try again later.", err)

	// ---------- 403 / 401 ----------
	case pgInsufficientPrivilege:
		return newAppError(http.StatusForbidden, "Database permission denied.", err)

	case pgInvalidAuthorization:
		return newAppError(http.StatusUnauthorized, "Invalid database credentials.", err)

	// ---------- 500 Internal Server Error ----------
	default:
		return newAppError(http.StatusInternalServerError, "Database error occurred.", err)
	}
}
