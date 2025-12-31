package controller

import (
	"fmt"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadProductImagesRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type UploadProductImagesResponse struct {
	Message string `json:"message"`
}
type UploadProductImagesController struct {
	usecase usecase.UploadProductImagesUseCase
}

func NewUploadProductImagesController(usecase usecase.UploadProductImagesUseCase) *UploadProductImagesController {
	return &UploadProductImagesController{
		usecase: usecase,
	}
}

func (c *UploadProductImagesController) Handle(fiberCtx *fiber.Ctx, req *UploadProductImagesRequest) (*UploadProductImagesResponse, error) {
	form, err := fiberCtx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("form error: %w", err)
	}

	files := form.File["images"]
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}
	fmt.Println(files)

	err = c.usecase.Execute(fiberCtx, req.ProductID, files)
	if err != nil {
		return nil, err
	}

	return &UploadProductImagesResponse{Message: "Product images uploaded successfully"}, nil

}
