package session

import "context"

func (s *SessionRepository) DeleteUserAllSession(ctx context.Context, token string) error {
	session, err := s.GetSessionData(ctx, token)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}

	sessionSetKey := s.userSessionsKey(session.UserID)
	tokens, err := s.client.SMembers(ctx, sessionSetKey).Result()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}
	pipe := s.client.Pipeline()

	keysToDele := make([]string, 0, len(tokens)+1)
	keysToDele = append(keysToDele, sessionSetKey)
	keysToDele = append(keysToDele, tokens...)
	pipe.Del(ctx, keysToDele...)

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}
