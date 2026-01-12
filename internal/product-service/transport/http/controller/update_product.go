package controller

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdateProductRequest struct {
	ProductID   uuid.UUID              `params:"product_id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Price       *float64               `json:"price"`
	StockCount  *int                   `json:"stock_count"`
	CategoryID  *uuid.UUID             `json:"category_id"`
	Attributes  map[string]interface{} `json:"attributes"`
}

type UpdateProductResponse struct {
	Message string `json:"message"`
}
type UpdateProductController struct {
	usecase usecase.UpdateProductUseCase
}

func NewUpdateProductController(usecase usecase.UpdateProductUseCase) *UpdateProductController {
	return &UpdateProductController{
		usecase: usecase,
	}
}

// Handle godoc
// @Summary update product
// @Description update a product
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param product body UpdateProductRequest true "Product"
// @Success 200 {object} UpdateProductResponse
// @Router /products/update/{product_id} [put]
func (c *UpdateProductController) Handle(fiberCtx *fiber.Ctx, req *UpdateProductRequest) (*UpdateProductResponse, error) {
	parsedUserID, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	p := &domain.UpdateProduct{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		StockCount:  req.StockCount,
		CategoryID:  req.CategoryID,
		Attributes:  req.Attributes,
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), parsedUserID, p)
	if err != nil {
		return nil, err
	}

	return &UpdateProductResponse{Message: "Product updated successfully"}, nil

}
