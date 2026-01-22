package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bilalabdelkadir/chis/pkg/shared"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// DecodeAndValidate parses JSON body and validates struct in one step
func DecodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return validate.Struct(dst)
}

// FormatErrors converts validation errors to readable FieldError structs
func FormatErrors(err error) []shared.FieldError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Not a validation error → return generic FieldError
		return []shared.FieldError{
			{Field: "unknown", Message: err.Error()},
		}
	}

	var errors []shared.FieldError
	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())
		msg := ""

		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
		case "gte":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
		case "lte":
			msg = fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
		case "oneof":
			msg = fmt.Sprintf("%s must be one of: %s", field, e.Param())
		case "uuid":
			msg = fmt.Sprintf("%s must be a valid UUID", field)
		case "url":
			msg = fmt.Sprintf("%s must be a valid URL", field)
		default:
			msg = fmt.Sprintf("%s is invalid", field)
		}

		errors = append(errors, shared.FieldError{
			Field:   field,
			Message: msg,
		})
	}

	return errors
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
