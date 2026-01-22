package apperror

import (
	"errors"
	"log"
	"net/http"

	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
)

// HandlerFunc is your custom handler that can return an error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Wrap converts a HandlerFunc into a standard http.HandlerFunc
func Wrap(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := h(w, r)

		if err == nil {
			return
		}

		var appErr *AppError

		if errors.As(err, &appErr) {
			// already set

		} else if pgErr := FromPostgres(err); pgErr != nil {
			appErr = pgErr
		} else if validator.IsValidationError(err) {
			details := validator.FormatErrors(err)
			appErr = ValidationFailed(details)
		} else {
			appErr = Internal("Internal server error")
		}

		appErr.WithContext(r)

		response.WriteJSON(w, appErr.Code, appErr)

		if appErr.Cause != nil {
			log.Printf("Error: %+v\n", appErr.Cause)
		}
	}
}
