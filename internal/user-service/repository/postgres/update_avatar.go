package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	return err
}