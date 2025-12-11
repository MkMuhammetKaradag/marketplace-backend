package session

import (
	"context"
	"time"
)

func (sm *SessionRepository) RefreshSession(ctx context.Context, token string, duration time.Duration) error {

	_, err := sm.client.Expire(ctx, token, duration).Result()
	if err != nil {
		return err
	}

	return nil
}
