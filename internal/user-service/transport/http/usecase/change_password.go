package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error
}
type changePasswordUseCase struct {
	repo domain.UserRepository
}

func NewChangePasswordUseCase(repo domain.UserRepository) ChangePasswordUseCase {
	return &changePasswordUseCase{
		repo: repo,
	}
}

func (u *changePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {

	err := u.repo.ChangePassword(ctx, userID, oldPassword, newPassword)
	if err != nil {
		return err
	}
	fmt.Println("Password changed successfully")

	return nil
}
