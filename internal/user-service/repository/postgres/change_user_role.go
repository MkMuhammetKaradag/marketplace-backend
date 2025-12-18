package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

const query = `
	UPDATE users
	SET user_role = $1
	WHERE id = $2

`

func (r *Repository) ChangeUserRole(ctx context.Context, userID uuid.UUID, role domain.UserRole) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to change user role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
