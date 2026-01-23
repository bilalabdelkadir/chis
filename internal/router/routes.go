package router

import (
	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/middleware"
)

func Setup(r *Router,
	authHandler *handler.AuthHandler,
	secret string,
) {
	// Public routes - no auth required
	r.Route("/auth", func(r *Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/api", func(r *Router) {
		r.Use(middleware.ValidateToken(secret)) // Only this group
	})

}
