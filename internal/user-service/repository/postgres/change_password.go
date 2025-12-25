package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	selectUserQuery = `
	SELECT password
	FROM users
	WHERE id = $1`
	changePasswordQuery = `
	UPDATE users SET password = $1,updated_at = NOW() WHERE id = $2`
)

func (r *Repository) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {

	var currentHash string

	err := r.db.QueryRowContext(ctx, selectUserQuery, userID).Scan(&currentHash)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash failed: %w", err)
	}

	_, err = r.db.ExecContext(ctx, changePasswordQuery, string(hashedPassword), userID)
	if err != nil {
		return fmt.Errorf("password update failed: %w", err)
	}

	return nil
}
