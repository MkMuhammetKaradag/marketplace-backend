package controller

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	StockCount  int       `json:"stock_count"`
	// Status      string    `json:"status"`
	SellerID    uuid.UUID `json:"seller_id"`
}

type CreateProductResponse struct {
	Message string `json:"message"`
}
type CreateProductController struct {
	usecase usecase.CreateProductUseCase
}

func NewCreateProductController(usecase usecase.CreateProductUseCase) *CreateProductController {
	return &CreateProductController{
		usecase: usecase,
	}
}

func (c *CreateProductController) Handle(fiberCtx *fiber.Ctx, req *CreateProductRequest) (*CreateProductResponse, error) {

	p := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		StockCount:  req.StockCount,
		// Status:      req.Status,
		SellerID:    req.SellerID,
	}
	err := c.usecase.Execute(fiberCtx, p)
	if err != nil {
		return nil, err
	}

	return &CreateProductResponse{Message: "Product created successfully"}, nil

}
