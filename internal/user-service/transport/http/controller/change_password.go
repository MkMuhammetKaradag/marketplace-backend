package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChangePasswordRequest struct {
	OldPassword      string `json:"old_password", binding:"required"`
	NewPassword      string `json:"new_password", binding:"required"`
	CloseAllSessions bool   `json:"close_all_sessions"`
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

	err = h.usecase.Execute(fbr.UserContext(), userID, req.OldPassword, req.NewPassword, req.CloseAllSessions)
	if err != nil {
		return nil, err
	}
	if req.CloseAllSessions {
		fbr.Set("X-Invalidate-User-All-Sessions", userID.String())
	}

	return &ChangePasswordResponse{Message: "Password changed successfully"}, nil
}
