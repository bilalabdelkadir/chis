package middleware

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/google/uuid"
)

func RequireAdmin(membershipRepo repository.MembershipRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			orgIDVal := r.Context().Value(OrgIDKey)
			if orgIDVal == nil {
				response.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "organization not found in context"})
				return
			}
			orgID, ok := orgIDVal.(uuid.UUID)
			if !ok {
				response.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid organization ID type"})
				return
			}

			membership, err := membershipRepo.FindByUserAndOrgID(r.Context(), userID, orgID)
			if err != nil {
				response.WriteJSON(w, http.StatusForbidden, map[string]string{"message": "access denied"})
				return
			}

			if membership.Role != "admin" {
				response.WriteJSON(w, http.StatusForbidden, map[string]string{"message": "admin access required"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
