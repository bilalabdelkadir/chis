package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, url string) (*redis.Client, error) {
	options := redis.Options{
		Addr: url,
	}

	rdb := redis.NewClient(&options)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
