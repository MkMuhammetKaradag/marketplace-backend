package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type CreateProductUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, req *domain.Product) error
}

type createProductUseCase struct {
	productRepository domain.ProductRepository
}

func NewCreateProductUseCase(productRepository domain.ProductRepository) CreateProductUseCase {
	return &createProductUseCase{
		productRepository: productRepository,
	}
}

func (c *createProductUseCase) Execute(ctx context.Context, userID uuid.UUID, req *domain.Product) error {

	sellerID, err := c.productRepository.GetSellerIDByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get seller ID: %w", err)
	}

	req.SellerID = sellerID
	return c.productRepository.CreateProduct(ctx, req)
}
