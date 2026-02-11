package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) CreateWithOutbox(ctx context.Context, user *domain.User, outbox *domain.OutboxMessage) error {
	// 1. Transaction başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	// Hata durumunda geri al
	defer tx.Rollback()

	// 2. Kullanıcıyı Kaydet (Users Tablosuna)
	// Not: Şifre hashleme işlemini UseCase'de yapmış olduğunu varsayıyorum.
	userQuery := `
        INSERT INTO users (username, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	row := tx.QueryRowContext(ctx, userQuery, user.Username, user.Email, user.Password)
	if err := row.Scan(&user.ID); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// 3. Outbox Mesajını Kaydet (Outbox_messages Tablosuna)
	outboxQuery := `
        INSERT INTO outbox_messages (payload, topic, status, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	row = tx.QueryRowContext(ctx, outboxQuery, outbox.Payload, outbox.Topic, outbox.Status)
	if err := row.Scan(&outbox.ID); err != nil {
		return fmt.Errorf("failed to insert outbox message: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) SignUpWithOutbox(ctx context.Context, user *domain.User, payload []byte) (uuid.UUID, string, error) {
	// 1. Şifre Hashleme ve Kod Üretme (Mevcut mantığın)
	hashedPassword, err := r.hashPassword(user.Password)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("hashing error: %w", err)
	}

	if user.ActivationCode == "" {
		code, err := generateRandomCode(6)
		if err != nil {
			return uuid.Nil, "", err
		}
		user.ActivationCode = code
	}

	if user.ActivationExpiry.IsZero() {
		user.ActivationExpiry = time.Now().Add(10 * time.Minute)
	}

	// 2. Transaction Başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	// 3. Kullanıcıyı Kaydet (Mevcut mantığın)
	userQuery := `
        INSERT INTO users (
            username, email, password,
            activation_code, activation_expiry,activation_id,
            failed_login_attempts, account_locked
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, userQuery,
		user.Username, user.Email, hashedPassword,
		user.ActivationCode, user.ActivationExpiry, user.ActivationID, 0, false,
	)

	if err != nil {
		if r.isDuplicateKeyError(err) {
			return uuid.Nil, "", fmt.Errorf("%w: username or email already exists", ErrDuplicateResource)
		}
		return uuid.Nil, "", fmt.Errorf("insert user error: %w", err)
	}

	// 4. YENİ: Outbox Tablosuna Mesajı Ekle (AYNI TRANSACTION)
	outboxQuery := `
        INSERT INTO outbox_messages (
            id, payload, topic, status, created_at
        ) VALUES ($1, $2, $3, $4, NOW())`

	outboxID := uuid.New()
	topic := "user-events" // Mesajın gideceği Kafka topic'i
	status := "PENDING"

	_, err = tx.ExecContext(ctx, outboxQuery, outboxID, payload, topic, status)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("outbox insert error: %w", err)
	}

	// 5. Her şey tamamsa COMMIT
	if err := tx.Commit(); err != nil {
		return uuid.Nil, "", fmt.Errorf("commit error: %w", err)
	}

	return user.ActivationID, user.ActivationCode, nil
}
