package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type UserCreatedUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, username string, email string) error
}
type userCreatedUseCase struct {
	repository domain.ProductRepository
}

func NewUserCreatedUseCase(repository domain.ProductRepository) UserCreatedUseCase {
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
