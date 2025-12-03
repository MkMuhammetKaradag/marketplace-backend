// internal/user-service/database/postgres/repository.go
package postgres

import (
	"database/sql"
	"errors"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/domain"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrMaxAttemptsReached = errors.New("maximum number of reset attempts reached")
	ErrTokenExpired       = errors.New("token expired")
	ErrActivationExpired  = errors.New("activation code expired")
	ErrInvalidActivation  = errors.New("invalid activation link or code expired")
	ErrInvalidCredentials = errors.New("invalid username, email or password")
	ErrAccountLocked      = errors.New("account is locked, please try again later")
	ErrDuplicateResource  = errors.New("duplicate resource")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg config.Config) (domain.PostgresRepository, error) {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	go repo.startCleanupJob(10 * time.Minute)

	return repo, nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
