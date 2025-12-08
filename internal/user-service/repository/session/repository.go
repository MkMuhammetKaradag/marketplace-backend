package session

import (
	"fmt"
	"marketplace/internal/user-service/config"

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

func (sm *SessionRepository) userSessionsKey(userID string) string {
	return fmt.Sprintf("user-service:user_sessions:%s", userID)
}
