package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"marketplace/internal/user-service/domain"

	"golang.org/x/crypto/bcrypt"
)

const maxLoginAttempts = 5

func (r *Repository) SignIn(ctx context.Context, identifier, password string) (*domain.User, error) {

	user := &domain.User{}
	query := `
        SELECT 
            u.id, u.username, u.email, u.password, u.failed_login_attempts, 
            u.account_locked, u.lock_until,
            COALESCE(BIT_OR(r.permissions), 0) as total_permissions
        FROM users u
        LEFT JOIN user_roles ur ON u.id = ur.user_id
        LEFT JOIN roles r ON ur.role_id = r.id
        WHERE (u.username = $1 OR u.email = $1) AND u.is_active = true
        GROUP BY u.id`

	var lockUntil sql.NullTime
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FailedLoginAttempts,
		&user.AccountLocked,
		&lockUntil,
		&user.Permissions,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	if user.AccountLocked && lockUntil.Valid && lockUntil.Time.After(time.Now()) {
		return nil, ErrAccountLocked
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Şifre yanlışsa, yeni bir işlem başlatarak deneme sayısını güncelleyin.
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("transaction begin failed: %w", err)
		}
		defer tx.Rollback()

		newAttempts := user.FailedLoginAttempts + 1
		var lockUntilTime sql.NullTime
		if newAttempts >= maxLoginAttempts {
			lockUntilTime.Time = time.Now().Add(1 * time.Minute)
			lockUntilTime.Valid = true
			user.AccountLocked = true
		}

		updateQuery := `
			UPDATE users SET failed_login_attempts = $1, account_locked = $2, lock_until = $3 WHERE id = $4`
		if _, updateErr := tx.ExecContext(ctx, updateQuery, newAttempts, user.AccountLocked, lockUntilTime, user.ID); updateErr != nil {
			fmt.Printf("Failed to update login attempts for user %s: %v\n", user.Email, updateErr)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("transaction commit failed: %w", err)
		}

		return nil, ErrInvalidCredentials
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE users SET failed_login_attempts = 0, account_locked = false, last_login = NOW(), lock_until = NULL WHERE id = $1`
	if _, updateErr := tx.ExecContext(ctx, updateQuery, user.ID); updateErr != nil {
		fmt.Printf("Failed to update last login for user %s: %v\n", user.Email, updateErr)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return user, nil
}
