package router

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/response"
)

func webhookTestHandler(w http.ResponseWriter, r *http.Request) error {
	orgId := r.Context().Value(middleware.OrgIDKey)
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "API key valid",
		"orgId":   orgId,
	})
	return nil
}

func Setup(r *Router,
	authHandler *handler.AuthHandler,
	apiKeyHandler *handler.ApiKeyHandler,
	apiKeyRepo *repository.ApiKeyRepository,
	secret string,
) {
	// Public routes - no auth required
	r.Route("/auth", func(r *Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/webhook", func(r *Router) {
		r.Use(middleware.ValidateApiKey(apiKeyRepo))
		r.Get("/test", webhookTestHandler)
	})

	r.Route("/api", func(r *Router) {
		r.Use(middleware.ValidateToken(secret))
		r.Route("/api-key", func(r *Router) {
			r.Post("/create", apiKeyHandler.Create)
		})
	})

}
