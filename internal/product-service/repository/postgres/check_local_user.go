package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const CHECK_LOCAL_USER_EXISTS = `SELECT EXISTS(SELECT 1 FROM local_users WHERE id = $1)`

func (r *Repository) CheckLocalUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, CHECK_LOCAL_USER_EXISTS, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking local user: %w", err)
	}
	return exists, nil
}
