package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordController struct {
	usecase usecase.ChangePasswordUseCase
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

func NewChangePasswordController(usecase usecase.ChangePasswordUseCase) *ChangePasswordController {
	return &ChangePasswordController{
		usecase: usecase,
	}
}

func (h *ChangePasswordController) Handle(fbr *fiber.Ctx, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	userID, err := uuid.Parse(fbr.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = h.usecase.Execute(fbr.UserContext(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{Message: "Password changed successfully"}, nil
}
