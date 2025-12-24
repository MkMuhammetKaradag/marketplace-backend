package usecase

import (
	"context"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, recordID uuid.UUID, password string) error
}
type resetPasswordUseCase struct {
	repo domain.UserRepository
}

func NewResetPasswordUseCase(repo domain.UserRepository) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		repo: repo,
	}
}

func (u *resetPasswordUseCase) Execute(ctx context.Context, recordID uuid.UUID, password string) error {

	_, err := u.repo.ResetPassword(ctx, recordID, password)
	if err != nil {
		return err
	}

	return nil
}
