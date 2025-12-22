package usecase

import (
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddUserRolerUseCase interface {
	Execute(ctx *fiber.Ctx, userID uuid.UUID, role string) error
}
type addUserRolerUseCase struct {
	userRepository domain.UserRepository
}

func NewAddUserRolerUseCase(repository domain.UserRepository) AddUserRolerUseCase {
	return &addUserRolerUseCase{
		userRepository: repository,
	}
}

func (u *addUserRolerUseCase) Execute(ctx *fiber.Ctx, userID uuid.UUID, role string) error {

	err := u.userRepository.AddUserRole(ctx.UserContext(), userID, role)
	if err != nil {
		return err
	}

	return nil
}
