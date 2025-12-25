// internal/user-service/transport/http/controller/all_signout.go
package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type AllSignOutRequest struct{}
type AllSignOutResponse struct {
	Message string `json:"message"`
}
type AllSignOutController struct {
	usecase usecase.AllSignOutUseCase
}

func NewAllSignOutController(usecase usecase.AllSignOutUseCase) *AllSignOutController {
	return &AllSignOutController{
		usecase: usecase,
	}
}

func (h *AllSignOutController) Handle(fbrCtx *fiber.Ctx, req *AllSignOutRequest) (*AllSignOutResponse, error) {
	err := h.usecase.Execute(fbrCtx)
	if err != nil {
		return nil, err
	}
	fbrCtx.Set("X-Invalidate-User-All-Sessions", "true")
	return &AllSignOutResponse{Message: "all signout  successfully"}, nil
}
