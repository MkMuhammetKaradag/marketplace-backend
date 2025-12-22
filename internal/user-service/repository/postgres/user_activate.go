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

func (r *Repository) UserActivate(ctx context.Context, activationID uuid.UUID, code string) (*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}
	defer tx.Rollback()

	// 1. Adım: Kullanıcıyı aktif et
	query := `
        UPDATE users 
        SET is_active = true, is_email_verified = true
        WHERE activation_id = $1 AND activation_code = $2 AND activation_expiry > $3
        RETURNING id, username, email;`

	var user domain.User
	err = tx.QueryRowContext(ctx, query, activationID, code, time.Now()).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidActivation
		}
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// 2. Adım: Kullanıcıya "Buyer" rolünü ata
	// Alt sorgu ile 'Buyer' rolünün ID'sini bulup user_roles tablosuna ekliyoruz.
	roleQuery := `
        INSERT INTO user_roles (user_id, role_id)
        SELECT $1, id FROM roles WHERE name = 'Buyer' LIMIT 1
        ON CONFLICT (user_id, role_id) DO NOTHING;`

	_, err = tx.ExecContext(ctx, roleQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign default buyer role: %w", err)
	}

	// 3. Adım: İşlemi onayla
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &user, nil
}
