package controller

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetProductRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type GetProductResponse struct {
	Product *domain.Product `json:"product"`
}

type GetProductController struct {
	usecase usecase.GetProductUseCase
}

func NewGetProductController(usecase usecase.GetProductUseCase) *GetProductController {
	return &GetProductController{
		usecase: usecase,
	}
}

// Handle godoc
// @Summary Get product
// @Description Get product by product id
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} GetProductResponse
// @Router /products/{product_id} [get]
func (c *GetProductController) Handle(fiberCtx *fiber.Ctx, req *GetProductRequest) (*GetProductResponse, error) {

	userIDStr := fiberCtx.Get("X-User-ID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = uuid.Nil
	}
	product, err := c.usecase.Execute(fiberCtx.UserContext(), userID, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &GetProductResponse{
		Product: product,
	}, nil
}
