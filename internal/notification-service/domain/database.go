package domain

import (
	"context"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	AddUser(ctx context.Context, userID uuid.UUID, username string, email string) error
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	Close() error
}
