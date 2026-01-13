// internal/basket-service/repository/postgres/repository.go
package postgres

import (
	"database/sql"
	"errors"
	"marketplace/internal/basket-service/config"
	"marketplace/internal/basket-service/domain"

	_ "github.com/lib/pq"
)

var (
	ErrBasketNotFound     = errors.New("basket not found")
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

func NewRepository(cfg config.Config) (domain.BasketPostgresRepository, error) {
	db, err := newPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	repo := &Repository{db: db}

	return repo, nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
