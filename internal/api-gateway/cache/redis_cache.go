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
	UserID      string
	Permissions int64
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
func (c *CacheManager) SetSession(ctx context.Context, token string, userID string, permissions int64) error {
	sessionKey := c.sessionKey(token)
	userIndexKey := fmt.Sprintf("user_sessions:%s", userID)

	session := SessionCache{
		UserID:      userID,
		Permissions: permissions,
	}

	data, _ := json.Marshal(session)

	// Redis Pipeline kullanarak atomik işlem yapalım
	pipe := c.client.TxPipeline()

	// 1. Session verisini kaydet
	pipe.Set(ctx, sessionKey, data, c.ttl)

	// 2. Bu token'ı kullanıcının aktif listesine ekle
	pipe.SAdd(ctx, userIndexKey, token)
	pipe.Expire(ctx, userIndexKey, c.ttl) // Liste de session süresi kadar yaşasın

	_, err := pipe.Exec(ctx)
	return err
}

// InvalidateSession cache'den session'ı sil (logout için)
func (c *CacheManager) InvalidateSession(ctx context.Context, token string) error {
	// Önce session'ı çekip userID'yi bulmamız lazım ki listeden silelim
	session, err := c.GetSession(ctx, token)
	if err != nil {
		return c.client.Del(ctx, c.sessionKey(token)).Err() // Bulamazsa bile silmeyi dene
	}

	userIndexKey := fmt.Sprintf("user_sessions:%s", session.UserID)

	pipe := c.client.TxPipeline()
	pipe.Del(ctx, c.sessionKey(token))
	pipe.SRem(ctx, userIndexKey, token) // Listeden bu tokenı çıkar

	_, err = pipe.Exec(ctx)
	return err
}
func (c *CacheManager) InvalidateAllUserSessions(ctx context.Context, userID string) error {
	userIndexKey := fmt.Sprintf("user_sessions:%s", userID)

	// 1. Kullanıcıya ait tüm tokenları al
	tokens, err := c.client.SMembers(ctx, userIndexKey).Result()
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	// 2. Tüm tokenları ve index listesini sil
	pipe := c.client.TxPipeline()
	for _, token := range tokens {
		pipe.Del(ctx, c.sessionKey(token))
	}
	pipe.Del(ctx, userIndexKey)

	_, err = pipe.Exec(ctx)
	return err
}

// sessionKey prefix ekleyerek key oluştur
func (c *CacheManager) sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

// Close Redis bağlantısını kapat
func (c *CacheManager) Close() error {
	return c.client.Close()
}
