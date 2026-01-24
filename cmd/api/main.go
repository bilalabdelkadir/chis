package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/logger"
	"github.com/bilalabdelkadir/chis/internal/middleware"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/internal/router"
	pb "github.com/bilalabdelkadir/chis/proto/delivery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func healthHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	return nil
}

func main() {
	logger.Setup()
	cfg, err := config.LoadEnv()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.Connect(cfg.DbUrl)
	if err != nil {
		slog.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("database connected.")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	orgRepo := repository.NewOrganizationRepository(pool)
	membershipRepo := repository.NewMembershipRepository(pool)
	apiKeyRepo := repository.NewApiKeyRepository(pool)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to connect to delivery service", "error", err)
		os.Exit(1)
	}

	deliveryClient := pb.NewDeliveryServiceClient(conn)

	// Handlers
	authHandler := handler.NewAuthHandler(userRepo, accountRepo, orgRepo, membershipRepo, cfg.JwtSecret)
	apiKeyHandler := handler.NewApiKeyHandler(membershipRepo, apiKeyRepo)
	webhookHandler := handler.NewWebhookHandler(deliveryClient)

	// Router
	r := router.NewRouter()
	r.Use(middleware.Logging)
	r.Get("/health", healthHandler)

	router.Setup(r, authHandler, apiKeyHandler, webhookHandler, apiKeyRepo, cfg.JwtSecret)

	slog.Info("server starting", "port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		slog.Error("failed to start http", "error", err)
		os.Exit(1)
	}
}
