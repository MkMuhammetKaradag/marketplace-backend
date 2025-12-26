package controller

import (
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadStoreBannerRequest struct {
	SellerID uuid.UUID `params:"seller_id"`
}

type UploadStoreBannerController struct {
	usecase usecase.UploadStoreBannerUseCase
}

type UploadStoreBannerResponse struct {
	Message string `json:"message"`
}

func NewUploadStoreBannerController(usecase usecase.UploadStoreBannerUseCase) *UploadStoreBannerController {
	return &UploadStoreBannerController{
		usecase: usecase,
	}
}

func (h *UploadStoreBannerController) Handle(fbr *fiber.Ctx, req *UploadStoreBannerRequest) (*UploadStoreBannerResponse, error) {

	userID, err := uuid.Parse(fbr.Get("X-User-ID"))

	if err != nil {
		return nil, err
	}
	fileHeader, err := fbr.FormFile("store_banner")
	if err != nil {
		return nil, err
	}

	if err := h.usecase.Execute(fbr.UserContext(), userID, req.SellerID, fileHeader); err != nil {
		return nil, err
	}

	return &UploadStoreBannerResponse{Message: "Store banner uploaded successfully"}, nil
}
