package controller

import (
	"context"
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserActivateRequest struct {
	ActivationID   uuid.UUID `json:"activation_id" validate:"required"`
	ActivationCode string    `json:"activation_code" validate:"required"`
}

type UserActivateResponse struct {
	Message string `json:"message"`
}
type UserActivateController struct {
	usecase usecase.UserActivateUseCase
}

func NewUserActivateController(usecase usecase.UserActivateUseCase) *UserActivateController {
	return &UserActivateController{
		usecase: usecase,
	}
}

func (h *UserActivateController) Handle(ctx context.Context, req *UserActivateRequest) (*UserActivateResponse, int, error) {
	status, err := h.usecase.Execute(ctx, req.ActivationID, req.ActivationCode)
	if err != nil {
		return nil, status, err
	}

	return &UserActivateResponse{Message: "user Useractivate"}, fiber.StatusOK, nil
}
