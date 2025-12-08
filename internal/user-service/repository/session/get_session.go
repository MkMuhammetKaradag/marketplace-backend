package session

import (
	"context"
	"encoding/json"
	"marketplace/internal/user-service/domain"

	"github.com/redis/go-redis/v9"
)

func (sm *SessionRepository) GetSessionData(ctx context.Context, token string) (*domain.SessionData, error) {
	val, err := sm.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return nil, domain.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	var sessionData domain.SessionData
	err = json.Unmarshal([]byte(val), &sessionData)
	if err != nil {
		return nil, err
	}
	return &sessionData, nil
}
