package usecase

import (
	"context"
	"errors"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type DeleteProductUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID, isAdmin bool) error
}

type deleteProductUseCase struct {
	productRepository domain.ProductRepository
}

func NewDeleteProductUseCase(productRepository domain.ProductRepository) DeleteProductUseCase {
	return &deleteProductUseCase{
		productRepository: productRepository,
	}
}

func (u *deleteProductUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID, isAdmin bool) error {

	existingProduct, err := u.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	if !isAdmin {
		sellerid, err := u.productRepository.GetSellerIDByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if existingProduct.SellerID != sellerid {
			return errors.New("unauthorized to update this product")
		}
	}

	return u.productRepository.SoftDeleteProduct(ctx, productID)
}
