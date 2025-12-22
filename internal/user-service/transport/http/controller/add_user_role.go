package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddUserRolerRequest struct {
	UserID uuid.UUID `params:"user_id"`
	Role   string    `json:"role"`
}

type AddUserRolerResponse struct {
	Message string `json:"message"`
}

type AddUserRolerController struct {
	usecase usecase.AddUserRolerUseCase
}

func NewAddUserRolerController(usecase usecase.AddUserRolerUseCase) *AddUserRolerController {
	return &AddUserRolerController{
		usecase: usecase,
	}
}

func (h *AddUserRolerController) Handle(fiberCtx *fiber.Ctx, req *AddUserRolerRequest) (*AddUserRolerResponse, error) {
	// currenUserRole := fiberCtx.Get("X-User-Role")
	// if currenUserRole != "admin" {
	// 	return nil, fiberCtx.SendStatus(fiber.StatusUnauthorized)
	// }

	err := h.usecase.Execute(fiberCtx, req.UserID, req.Role)
	if err != nil {
		return nil, err
	}

	return &AddUserRolerResponse{Message: "User role added successfully"}, nil
}
