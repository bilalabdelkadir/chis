package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/queue"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/internal/scheduler"
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

	ctx := context.Background()

	rdsClient, err := queue.NewRedisClient(ctx, cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	queue := queue.NewQueue(rdsClient, QueueName)

	fmt.Println("scheduler started")

	// Worker
	w := scheduler.NewScheduler(messageRepo, queue)
	w.Start(context.Background())
}
