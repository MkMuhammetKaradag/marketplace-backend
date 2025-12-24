package controller

import (
	"context"
	"marketplace/internal/user-service/transport/http/usecase"
)

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"`
}

type ForgotPasswordController struct {
	usecase usecase.ForgotPasswordUseCase
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

func NewForgotPasswordController(usecase usecase.ForgotPasswordUseCase) *ForgotPasswordController {
	return &ForgotPasswordController{
		usecase: usecase,
	}
}

func (h *ForgotPasswordController) Handle(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {

	err := h.usecase.Execute(ctx, req.Identifier)
	if err != nil {
		return nil, err
	}

	return &ForgotPasswordResponse{Message: "Password reset email sent successfully"}, nil
}
