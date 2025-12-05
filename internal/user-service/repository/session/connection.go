package session

import (
	"context"
	"fmt"

	"marketplace/internal/user-service/config"

	"github.com/redis/go-redis/v9"
)

func newRedisDB(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisSession.Addr,
		Password: cfg.RedisSession.Password,
		DB:       cfg.RedisSession.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Successfully connected to Redis.")

	return client, nil
}
