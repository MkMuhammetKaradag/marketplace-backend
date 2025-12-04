// internal/user-service/transport/http/usecase/signup.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
)

type SignUpUseCase interface {
	Execute(ctx context.Context, user *domain.User) error
}
type signUpUseCase struct {
	userRepository domain.UserRepository
}

type SignUpRequest struct {
	Username string
	Email    string
	Password string
}

func NewSignUpUseCase(repository domain.UserRepository) SignUpUseCase {
	return &signUpUseCase{
		userRepository: repository,
	}
}

func (u *signUpUseCase) Execute(ctx context.Context, user *domain.User) error {

	id, code, err := u.userRepository.SignUp(ctx, user)
	if err != nil {
		return err
	}

	fmt.Printf("id:%v ,    code:%v \n", id, code)

	return nil
}
