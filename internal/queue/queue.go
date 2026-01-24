package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	rdsClient *redis.Client
	name      string
}

func NewQueue(client *redis.Client, queueName string) *Queue {
	return &Queue{
		rdsClient: client,
		name:      queueName,
	}
}

func (q *Queue) Push(ctx context.Context, messageID string) error {
	if err := q.rdsClient.LPush(ctx, q.name, messageID).Err(); err != nil {
		return err
	}
	return nil
}

func (q *Queue) Pop(ctx context.Context) (string, error) {
	val := q.rdsClient.BRPop(ctx, 3*time.Second, q.name)
	if val.Err() != nil {
		return "", val.Err()
	}
	return val.Val()[1], nil
}
