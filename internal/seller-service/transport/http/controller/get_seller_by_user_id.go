package controller

import (
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetSellerByUserIDRequest struct {
}

type GetSellerByUserIDResponse struct {
	Message string         `json:"message"`
	Seller  *domain.Seller `json:"seller"`
}

type GetSellerByUserIDController struct {
	usecase usecase.GetSellerByUserIDUseCase
}

func NewGetSellerByUserIDController(usecase usecase.GetSellerByUserIDUseCase) *GetSellerByUserIDController {
	return &GetSellerByUserIDController{
		usecase: usecase,
	}
}

func (h *GetSellerByUserIDController) Handle(fbrCtx *fiber.Ctx, req *GetSellerByUserIDRequest) (*GetSellerByUserIDResponse, error) {
	parsedUserID, err := uuid.Parse(fbrCtx.Get("X-User-ID"))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	 seller, err := h.usecase.Execute(fbrCtx.UserContext(), parsedUserID);
	 
	 if err != nil {
		return nil, err
	}
	return &GetSellerByUserIDResponse{Message: "Seller fetched successfully", Seller: seller}, nil
}
