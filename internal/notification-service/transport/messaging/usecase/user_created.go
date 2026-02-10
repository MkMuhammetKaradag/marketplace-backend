// internal/notification-service/transport/messaging/usecase/user_created.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type UserCreatedUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, username string, email string) error
}
type userCreatedUseCase struct {
	repository domain.NotificationRepository
}

func NewUserCreatedUseCase(repository domain.NotificationRepository) UserCreatedUseCase {
	return &userCreatedUseCase{
		repository: repository,
	}
}

func (u *userCreatedUseCase) Execute(ctx context.Context, userID uuid.UUID, username string, email string) error {

	err := u.repository.AddUser(ctx, userID, username, email)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}
