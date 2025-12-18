package controller

import (
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChangeUserRolerRequest struct {
	UserID uuid.UUID       `params:"user_id"`
	Role   domain.UserRole `json:"role"`
}

type ChangeUserRolerResponse struct {
	Message string `json:"message"`
}

type ChangeUserRolerController struct {
	usecase usecase.ChangeUserRolerUseCase
}

func NewChangeUserRolerController(usecase usecase.ChangeUserRolerUseCase) *ChangeUserRolerController {
	return &ChangeUserRolerController{
		usecase: usecase,
	}
}

func (h *ChangeUserRolerController) Handle(fiberCtx *fiber.Ctx, req *ChangeUserRolerRequest) (*ChangeUserRolerResponse, error) {
	currenUserRole := fiberCtx.Get("X-User-Role")
	if currenUserRole != "admin" {
		return nil, fiberCtx.SendStatus(fiber.StatusUnauthorized)
	}

	err := h.usecase.Execute(fiberCtx, req.UserID, req.Role)
	if err != nil {
		return nil, err
	}

	return &ChangeUserRolerResponse{Message: "User role changed successfully"}, nil
}
