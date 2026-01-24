package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/delivery"
	"github.com/bilalabdelkadir/chis/internal/logger"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	pb "github.com/bilalabdelkadir/chis/proto/delivery"
	"google.golang.org/grpc"
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

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	q := queue.NewQueue(rdsClient, QueueName)

	messageRepo := repository.NewMessageRepository(pool)

	deliveryService := delivery.NewDeliveryService(messageRepo, q)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("failed to listen", "error", err, "port", 50051)
		os.Exit(1)
	}

	slog.Info("delivery_service_starting", "port", 50051)

	grpcServer := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(grpcServer, deliveryService)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("grpc_server_failed", "error", err)
		os.Exit(1)
	}
}
