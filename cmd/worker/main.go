package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/internal/worker"
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

	messageRepo := repository.NewMessageRepository(pool)
	attemptRepo := repository.NewDeliveryAttemptsRepository(pool)

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	queue := queue.NewQueue(rdsClient, QueueName)

	fmt.Println("worker started")

	// Worker
	w := worker.NewWorker(messageRepo, attemptRepo, queue)
	w.Start(context.Background())

}
