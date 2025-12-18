package usecase

import (
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChangeUserRolerUseCase interface {
	Execute(ctx *fiber.Ctx, userID uuid.UUID, role domain.UserRole) error
}
type changeUserRolerUseCase struct {
	userRepository domain.UserRepository
}

func NewChangeUserRolerUseCase(repository domain.UserRepository) ChangeUserRolerUseCase {
	return &changeUserRolerUseCase{
		userRepository: repository,
	}
}

func (u *changeUserRolerUseCase) Execute(ctx *fiber.Ctx, userID uuid.UUID, role domain.UserRole) error {

	err := u.userRepository.ChangeUserRole(ctx.UserContext(), userID, role)
	if err != nil {
		return err
	}

	return nil
}
