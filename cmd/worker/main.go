package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/logger"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/internal/worker"
)

var QueueName = "main"

func main() {
	logger.Setup()

	cfg, err := config.LoadEnv()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.Connect(cfg.DbUrl)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("database_connected")

	messageRepo := repository.NewMessageRepository(pool)
	attemptRepo := repository.NewDeliveryAttemptsRepository(pool)

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	queue := queue.NewQueue(rdsClient, QueueName)

	slog.Info("worker_started")

	w := worker.NewWorker(messageRepo, attemptRepo, queue)
	w.Start(context.Background())
}
