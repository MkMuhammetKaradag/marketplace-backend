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

// Redis anahtarı oluşturma yardımcı fonksiyonu
func (r *BasketRedisRepository) getBasketKey(userID string) string {
	return fmt.Sprintf("basket:%s", userID)
}

// UpdateBasket hem yeni ürün eklemek hem de miktar güncellemek için kullanılır
func (r *BasketRedisRepository) UpdateBasket(ctx context.Context, basket *domain.Basket) error {
	// 1. Sepeti JSON stringine çevir
	jsonData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("failed to marshal basket: %w", err)
	}

	// 2. Redis'e kaydet (Örnek: 7 gün ömür biçiyoruz)
	key := r.getBasketKey(basket.UserID.String())
	err = r.client.Set(ctx, key, jsonData, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save basket to redis: %w", err)
	}

	return nil
}

// GetBasket kullanıcı id'sine göre sepeti getirir
func (r *BasketRedisRepository) GetBasket(ctx context.Context, userID string) (*domain.Basket, error) {
	key := r.getBasketKey(userID)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Sepet yoksa boş bir sepet dönüyoruz (Hata değil)
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
