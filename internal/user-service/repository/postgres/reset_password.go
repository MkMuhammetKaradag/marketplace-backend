package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	selectTokenQuery    = `SELECT user_id,expires_at,attempt_count FROM forgot_passwords WHERE id = $1`
	deleteTokenQuery    = `DELETE FROM forgot_passwords WHERE id = $1`
	updatePasswordQuery = `UPDATE users SET password = $1 WHERE id = $2`
)

func (r *Repository) ResetPassword(ctx context.Context, recordID uuid.UUID, newPassword string) (uuid.UUID, error) {

	var (
		userID       uuid.UUID
		expiresAt    time.Time
		attemptCount int
	)

	// Token'ı kontrol et
	err := r.db.QueryRowContext(ctx, selectTokenQuery, recordID).Scan(&userID, &expiresAt, &attemptCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userID, ErrInvalidToken
		}
		return userID, fmt.Errorf("select token error: %w", err)
	}

	// Token süresi dolmuş mu?
	if time.Now().After(expiresAt) {
		_ = r.DeleteForgotPassword(ctx, recordID)
		return userID, ErrTokenExpired
	}

	// 3 kez denediyse token geçersiz
	if attemptCount >= 3 {
		_ = r.DeleteForgotPassword(ctx, recordID)
		return userID, ErrMaxAttemptsReached
	}

	hashedPassword, err := r.hashPassword(newPassword)
	if err != nil {
		return userID, fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback()

	// Şifreyi güncelle
	if _, err := tx.ExecContext(ctx, updatePasswordQuery, hashedPassword, userID); err != nil {
		return userID, fmt.Errorf("update password error: %w", err)
	}

	// Token'ı sil (Artık geçersiz kıl)
	if _, err := tx.ExecContext(ctx, deleteTokenQuery, recordID); err != nil {
		return userID, fmt.Errorf("failed to delete used token: %w", err)
	}

	// Her şey tamamsa onayla
	if err := tx.Commit(); err != nil {
		return userID, err
	}
	return userID, nil
}
func (r *Repository) DeleteForgotPassword(ctx context.Context, recordID uuid.UUID) error {
	const query = `DELETE FROM forgot_passwords WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, recordID)
	return err
}
