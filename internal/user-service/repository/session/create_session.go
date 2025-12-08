package session

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/user-service/domain"
	"time"
)

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
