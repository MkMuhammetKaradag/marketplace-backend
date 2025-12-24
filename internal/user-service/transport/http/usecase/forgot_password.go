package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, identifier string) error
}
type forgotPasswordUseCase struct {
	repo domain.UserRepository
}

func NewForgotPasswordUseCase(repo domain.UserRepository) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		repo: repo,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, identifier string) error {

	result, err := u.repo.ForgotPassword(ctx, identifier)
	if err != nil {
		return err
	}
	// fmt.Println("username", result.Username)
	// fmt.Println("email", result.Email)
	fmt.Println("token", result.Token)
	return nil
}
