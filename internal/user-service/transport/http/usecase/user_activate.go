package usecase

import (
	"context"
	"log"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

type UserActivateUseCase interface {
	Execute(ctx context.Context, activationID uuid.UUID, activationCode string) error
}
type UseractivateUseCase struct {
	userRepository domain.UserRepository
}

func NewUserActivateUseCase(repository domain.UserRepository) UserActivateUseCase {
	return &UseractivateUseCase{
		userRepository: repository,
	}
}

func (u *UseractivateUseCase) Execute(ctx context.Context, activationID uuid.UUID, activationCode string) error {

	user, err := u.userRepository.UserActivate(ctx, activationID, activationCode)
	if err != nil {
		return err
	}
	log.Printf("Activated user: %+v\n", user)

	return nil
}
