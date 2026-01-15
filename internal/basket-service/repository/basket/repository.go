// internal/basket-service/repository/basket/repository.go
package basket

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/basket-service/config"
	"marketplace/internal/basket-service/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type BasketRedisRepository struct {
	client *redis.Client
}

func NewBasketRedisRepository(cfg config.Config) (*BasketRedisRepository, error) {
	client, err := newRedisDB(cfg)
	if err != nil {
		return nil, err
	}

	return &BasketRedisRepository{
		client: client,
	}, nil
}

func (r *BasketRedisRepository) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Helper function to create Redis key
func (r *BasketRedisRepository) getBasketKey(userID string) string {
	return fmt.Sprintf("basket:%s", userID)
}

// UpdateBasket is used to both add new products and update quantities
func (r *BasketRedisRepository) UpdateBasket(ctx context.Context, basket *domain.Basket) error {
	// 1. Convert the basket to a JSON string
	jsonData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("failed to marshal basket: %w", err)
	}

	// 2. Save to Redis (Example: We're giving it a 7-day lifespan)
	key := r.getBasketKey(basket.UserID.String())
	err = r.client.Set(ctx, key, jsonData, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save basket to redis: %w", err)
	}

	return nil
}

// GetBasket returns the user's basket based on user id
func (r *BasketRedisRepository) GetBasket(ctx context.Context, userID string) (*domain.Basket, error) {
	key := r.getBasketKey(userID)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// If there's no basket, we're returning an empty basket (Not a mistake)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var basket domain.Basket
	err = json.Unmarshal([]byte(val), &basket)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal basket: %w", err)
	}

	return &basket, nil
}

// ClearBasket completely deletes the user's basket key from Redis.
func (r *BasketRedisRepository) ClearBasket(ctx context.Context, userID string) error {
	key := r.getBasketKey(userID)

	// Delete the key from Redis
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to clear basket from redis: %w", err)
	}

	return nil
}
