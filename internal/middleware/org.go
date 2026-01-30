package middleware

import (
	"context"
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/google/uuid"
)

func ValidateOrgAccess(membershipRepo repository.MembershipRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgHeader := r.Header.Get("X-Org-ID")
			if orgHeader == "" {
				response.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "missing X-Org-ID header"})
				return
			}

			orgID, err := uuid.Parse(orgHeader)
			if err != nil {
				response.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid X-Org-ID header"})
				return
			}

			userIDVal := r.Context().Value(UserIDKey)
			if userIDVal == nil {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "user not found in context"})
				return
			}
			userID, ok := userIDVal.(uuid.UUID)
			if !ok {
				response.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid user ID type"})
				return
			}

			_, err = membershipRepo.FindByUserAndOrgID(r.Context(), userID, orgID)
			if err != nil {
				response.WriteJSON(w, http.StatusForbidden, map[string]string{"message": "you are not a member of this organization"})
				return
			}

			ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
