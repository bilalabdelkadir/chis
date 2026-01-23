package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bilalabdelkadir/chis/internal/auth"
	"github.com/bilalabdelkadir/chis/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userId"

func ValidateToken(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
				return
			}

			tokenStr := parts[1]

			userId, err := auth.ValidateToken(tokenStr, secret)

			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
