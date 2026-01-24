package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/delivery"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	pb "github.com/bilalabdelkadir/chis/proto/delivery"
	"google.golang.org/grpc"
)

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

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	q := queue.NewQueue(rdsClient, QueueName)

	messageRepo := repository.NewMessageRepository(pool)

	deliveryService := delivery.NewDeliveryService(messageRepo, q)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delivery service starting on port 50051")

	grpcServer := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(grpcServer, deliveryService)
	grpcServer.Serve(lis)
}
