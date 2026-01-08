package postgres

import (
	"context"

	"github.com/google/uuid"
)

const (
	ADD_USER = `
        INSERT INTO local_users (id, username, email, updated_at) 
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (id) 
        DO UPDATE SET 
            username = EXCLUDED.username, 
            email = EXCLUDED.email, 
            updated_at = NOW()`
)

func (r *Repository) AddUser(ctx context.Context, userID uuid.UUID, username string, email string) error {
	_, err := r.db.ExecContext(ctx, ADD_USER, userID, username, email)
	if err != nil {
		return err
	}
	return nil
}
