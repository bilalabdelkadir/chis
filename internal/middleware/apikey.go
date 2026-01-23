package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/helper"
	"github.com/bilalabdelkadir/chis/pkg/response"
)

const OrgIDKey contextKey = "orgId"

func ValidateApiKey(apiKeyRepo *repository.ApiKeyRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKeyHeader := r.Header.Get("X-API-Key")
			if apiKeyHeader == "" {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid ApiKey"})
				return
			}

			hashedApiKey := helper.HashSHA256(apiKeyHeader)

			apiKey, err := apiKeyRepo.FindByHashedKey(r.Context(), hashedApiKey)

			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid ApiKey"})
				return
			}

			if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "expired ApiKey"})
				return
			}
			ctx := context.WithValue(r.Context(), OrgIDKey, apiKey.OrgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
