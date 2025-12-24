package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/user-service/domain"
	"time"

	"github.com/google/uuid"
)

const (
	userQuery          = `SELECT id, username, email FROM users WHERE username = $1 OR email = $1`
	checkExistingToken = `SELECT created_at FROM forgot_passwords WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	deleteOldToken     = `DELETE FROM forgot_passwords WHERE user_id = $1`
	insertQuery        = `INSERT INTO forgot_passwords (user_id, expires_at) VALUES ($1, $2) RETURNING id`
)

func (r *Repository) ForgotPassword(ctx context.Context, identifier string) (*domain.ForgotPasswordResult, error) {

	var (
		userID    uuid.UUID
		username  string
		email     string
		tokenID   uuid.UUID
		createdAt time.Time
	)

	err := r.db.QueryRowContext(ctx, userQuery, identifier).Scan(&userID, &username, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user by identifier: %w", err)
	}

	err = r.db.QueryRowContext(ctx, checkExistingToken, userID).Scan(&createdAt)
	if err == nil && time.Since(createdAt) < time.Minute {
		return nil, fmt.Errorf("please wait at least 1 minute before requesting a new link")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, deleteOldToken, userID); err != nil {
		return nil, fmt.Errorf("failed to delete old tokens: %w", err)
	}
	
	expiresAt := time.Now().Add(time.Minute * 15)
	err = tx.QueryRowContext(ctx, insertQuery, userID, expiresAt).Scan(&tokenID)
	if err != nil {
		return nil, err
	}

	// Transaction'ı onayla
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.ForgotPasswordResult{
		Username: username,
		Email:    email,
		Token:    tokenID.String(),
	}, nil
}
