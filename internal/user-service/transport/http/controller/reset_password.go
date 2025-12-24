package controller

import (
	"context"
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/google/uuid"
)

type ResetPasswordRequest struct {
	RecordID uuid.UUID `json:"record_id" binding:"required,min=20,max=20"`
	Password string    `json:"password" binding:"required,min=8,max=16"`
}

type ResetPasswordController struct {
	usecase usecase.ResetPasswordUseCase
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

func NewResetPasswordController(usecase usecase.ResetPasswordUseCase) *ResetPasswordController {
	return &ResetPasswordController{
		usecase: usecase,
	}
}

func (h *ResetPasswordController) Handle(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, error) {

	err := h.usecase.Execute(ctx, req.RecordID, req.Password)
	if err != nil {
		return nil, err
	}

	return &ResetPasswordResponse{Message: "Password reset successfully"}, nil
}
