package usecase

import (
	"fmt"
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateRoleUsecase interface {
	Execute(fiberCtx *fiber.Ctx, createdBy uuid.UUID, permissions int64, name string) error
}
type createRoleUsecase struct {
	repo domain.UserRepository
}

func NewCreateRoleUseCase(repo domain.UserRepository) CreateRoleUsecase {
	return &createRoleUsecase{
		repo: repo,
	}
}

func (u *createRoleUsecase) Execute(fiberCtx *fiber.Ctx, createdBy uuid.UUID, permissions int64, name string) error {

	if !domain.IsValidPermission(permissions) {
		return fmt.Errorf("invalid permissions")
	}

	roleID, err := u.repo.CreateRole(fiberCtx.UserContext(), createdBy, name, permissions)
	if err != nil {
		return fmt.Errorf("role not created: %w", err)
	}

	fmt.Println("roleid", roleID)

	return nil
}
