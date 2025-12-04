package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Close() error
	SignUp(ctx context.Context, user *User) (uuid.UUID, string, error)
	UserActivate(ctx context.Context, activationID uuid.UUID, code string) (*User, error)
}
