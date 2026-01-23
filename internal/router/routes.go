package router

import (
	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/repository"
)

func Setup(r *Router,
	authHandler *handler.AuthHandler,
	apiKeyHandler *handler.ApiKeyHandler,
	webhookHandler *handler.WebhookHandler,
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
		r.Post("/send", webhookHandler.Send)
	})

	r.Route("/api", func(r *Router) {
		r.Use(middleware.ValidateToken(secret))
		r.Route("/api-key", func(r *Router) {
			r.Post("/create", apiKeyHandler.Create)
		})
	})

}
