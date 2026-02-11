package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Close() error
	SignUp(ctx context.Context, user *User) (uuid.UUID, string, error)
	UserActivate(ctx context.Context, activationID uuid.UUID, code string) (*User, error)
	SignIn(ctx context.Context, identifier, password string) (*User, error)
	AddUserRole(ctx context.Context, userID uuid.UUID, roleName string) error
	CreateRole(ctx context.Context, createdBy uuid.UUID, name string, permissions int64) (uuid.UUID, error)
	ForgotPassword(ctx context.Context, identifier string) (*ForgotPasswordResult, error)
	ResetPassword(ctx context.Context, recordID uuid.UUID, newPassword string) (uuid.UUID, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error

	CreateWithOutbox(ctx context.Context, user *User, outbox *OutboxMessage) error
	SignUpWithOutbox(ctx context.Context, user *User, payload []byte) (uuid.UUID, string, error)
	GetPendingOutboxMessages(ctx context.Context, limit int) ([]OutboxMessage, error)
	MarkOutboxAsProcessed(ctx context.Context, id uuid.UUID) error
}
type ForgotPasswordResult struct {
	UserID   uuid.UUID
	Username string
	Email    string
	Token    string
}

type OutboxMessage struct {
	ID        uuid.UUID
	Payload   []byte
	Topic     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
