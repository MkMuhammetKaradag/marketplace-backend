package session

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
)

func (sm *SessionRepository) DeleteSession(ctx context.Context, token string) error {
	sessionData, err := sm.GetSessionData(ctx, token)
	if err != nil {
		return err
	}

	if sessionData == nil {
		return domain.ErrUnauthorized
	}

	pipe := sm.client.Pipeline()
	pipe.Del(ctx, token)
	pipe.SRem(ctx, sm.userSessionsKey(sessionData.UserID), token)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil

}
