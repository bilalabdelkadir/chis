package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/internal/router"
)

func healthHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	return nil
}

var QueueName = "main"

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.Connect(cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fmt.Println("database connected.")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	orgRepo := repository.NewOrganizationRepository(pool)
	membershipRepo := repository.NewMembershipRepository(pool)
	apiKeyRepo := repository.NewApiKeyRepository(pool)
	messageRepo := repository.NewMessageRepository(pool)

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	queue := queue.NewQueue(rdsClient, QueueName)

	fmt.Println("database connected.")

	// Handlers
	authHandler := handler.NewAuthHandler(userRepo, accountRepo, orgRepo, membershipRepo, cfg.JwtSecret)
	apiKeyHandler := handler.NewApiKeyHandler(membershipRepo, apiKeyRepo)
	webhookHandler := handler.NewWebhookHandler(messageRepo, queue)

	// Router
	r := router.NewRouter()
	r.Use(middleware.Logging)
	r.Get("/health", healthHandler)

	router.Setup(r, authHandler, apiKeyHandler, webhookHandler, apiKeyRepo, cfg.JwtSecret)

	fmt.Println("server starting on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
