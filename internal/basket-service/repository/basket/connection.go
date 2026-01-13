// internal/basket-service/repository/redis/connection.go
package basket

import (
	"context"
	"fmt"

	"marketplace/internal/basket-service/config"

	"github.com/redis/go-redis/v9"
)

func newRedisDB(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Successfully connected to Redis.")

	return client, nil
}
