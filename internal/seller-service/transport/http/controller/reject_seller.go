package controller

import (
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type RejectSellerRequest struct {
	SellerId string `params:"seller_id" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
}

type RejectSellerResponse struct {
	Message string `json:"message"`
}

type RejectSellerController struct {
	usecase usecase.RejectSellerUseCase
}

func NewRejectSellerController(usecase usecase.RejectSellerUseCase) *RejectSellerController {
	return &RejectSellerController{
		usecase: usecase,
	}
}

func (h *RejectSellerController) Handle(fbrCtx *fiber.Ctx, req *RejectSellerRequest) (*RejectSellerResponse, error) {

	approvedBy := fbrCtx.Get("X-User-ID")
	if approvedBy == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := h.usecase.Execute(fbrCtx.UserContext(), req.SellerId, approvedBy, req.Reason); err != nil {
		return nil, err
	}
	return &RejectSellerResponse{Message: "Seller rejected successfully"}, nil
}
