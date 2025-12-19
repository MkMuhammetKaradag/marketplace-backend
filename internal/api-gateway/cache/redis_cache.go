// internal/api-gateway/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionCache struct {
	UserID string
	Role   string
}

type CacheManager struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheManager(redisAddr string, password string, db int, ttl time.Duration) (*CacheManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})

	// Bağlantıyı test et
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis bağlantısı başarısız: %w", err)
	}

	return &CacheManager{
		client: client,
		ttl:    ttl,
	}, nil
}

// GetSession cache'den session bilgisini al
func (c *CacheManager) GetSession(ctx context.Context, token string) (*SessionCache, error) {
	key := c.sessionKey(token)

	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("seesion not found")
	}
	if err != nil {
		return nil, fmt.Errorf("session read error : %w", err)
	}

	var session SessionCache
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("session deserialize error : %w", err)
	}

	return &session, nil
}

// SetSession cache'e session bilgisini kaydet
func (c *CacheManager) SetSession(ctx context.Context, token string, userID string, role string) error {
	key := c.sessionKey(token)

	session := SessionCache{
		UserID: userID,
		Role:   role,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("session serialize error : %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("session write error : %w", err)
	}

	return nil
}

// InvalidateSession cache'den session'ı sil (logout için)
func (c *CacheManager) InvalidateSession(ctx context.Context, token string) error {
	key := c.sessionKey(token)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("session delete error : %w", err)
	}

	return nil
}

// sessionKey prefix ekleyerek key oluştur
func (c *CacheManager) sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

// Close Redis bağlantısını kapat
func (c *CacheManager) Close() error {
	return c.client.Close()
}
