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
	dashboardHandler *handler.DashboardHandler,
	orgHandler *handler.OrganizationHandler,
	apiKeyRepo repository.ApiKeyRepository,
	membershipRepo repository.MembershipRepository,
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

		// Org management routes (no org context needed)
		r.Get("/orgs", orgHandler.ListOrgs)
		r.Post("/orgs", orgHandler.CreateOrg)

		// Org-scoped routes (require X-Org-ID header)
		r.Route("/", func(r *Router) {
			r.Use(middleware.ValidateOrgAccess(membershipRepo))

			r.Route("/api-key", func(r *Router) {
				r.Post("/create", apiKeyHandler.Create)
				r.Get("/list", apiKeyHandler.List)
				r.Delete("/{id}", apiKeyHandler.Delete)
			})

			r.Route("/dashboard", func(r *Router) {
				r.Get("/stats", dashboardHandler.Stats)
			})

			r.Get("/webhook-logs", dashboardHandler.WebhookLogs)
			r.Get("/webhook-logs/{id}", dashboardHandler.WebhookLogDetail)
		})
	})
}
