package postgres

import (
	"context"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

const (
	GET_USER = `
        SELECT id, username, email 
        FROM local_users 
        WHERE id = $1`
)

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, GET_USER, userID)
	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}
