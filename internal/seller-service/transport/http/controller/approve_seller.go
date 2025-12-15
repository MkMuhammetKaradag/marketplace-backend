package controller

import (
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type ApproveSellerRequest struct {
	SellerId string `params:"seller_id" validate:"required"`
}

type ApproveSellerResponse struct {
	Message string `json:"message"`
}

type ApproveSellerController struct {
	usecase usecase.ApproveSellerUseCase
}

func NewApproveSellerController(usecase usecase.ApproveSellerUseCase) *ApproveSellerController {
	return &ApproveSellerController{
		usecase: usecase,
	}
}

func (h *ApproveSellerController) Handle(fbrCtx *fiber.Ctx, req *ApproveSellerRequest) (*ApproveSellerResponse, error) {

	approvedBy := fbrCtx.Get("X-User-ID")
	if approvedBy == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := h.usecase.Execute(fbrCtx.UserContext(), req.SellerId, approvedBy); err != nil {
		return nil, err
	}
	return &ApproveSellerResponse{Message: "Seller approved successfully"}, nil
}
