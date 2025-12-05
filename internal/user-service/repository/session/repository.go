package session

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(cfg config.Config) (*SessionRepository, error) {
	client, err := newRedisDB(cfg)
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		client: client,
	}, nil
}

func (sm *SessionRepository) CreateSession(ctx context.Context, token string, duration time.Duration, data *domain.SessionData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	pipe := sm.client.Pipeline()
	pipe.Set(ctx, token, jsonData, duration)
	pipe.SAdd(ctx, sm.userSessionsKey(data.UserID), token)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}
func (sm *SessionRepository) userSessionsKey(userID string) string {
	return fmt.Sprintf("user-service:user_sessions:%s", userID)
}
