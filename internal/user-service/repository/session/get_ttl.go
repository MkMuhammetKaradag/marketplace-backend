package session

import (
	"context"
	"marketplace/internal/user-service/domain"
	"time"
)

func (sm *SessionRepository) GetTTL(ctx context.Context, token string) (time.Duration, error) {
	// Redis'ten anahtarın kalan süresini al
	ttl, err := sm.client.TTL(ctx, token).Result()
	if err != nil {
		return 0, err
	}
	// Redis TTL komutu özel değerler döndürür:
	// -1 (anahtar var ama süresiz)
	// -2 (anahtar yok)
	if ttl == -2 {

		return 0, domain.ErrUnauthorized
	}
	if ttl == -1 {

		return 0, nil
	}

	return ttl, nil
}
