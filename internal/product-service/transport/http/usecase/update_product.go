package usecase

import (
	"context"
	"errors"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type UpdateProductUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, req *domain.UpdateProduct) error
}

type updateProductUseCase struct {
	productRepository domain.ProductRepository
	aiProvider        domain.AiProvider
}

func NewUpdateProductUseCase(productRepository domain.ProductRepository, aiProvider domain.AiProvider) UpdateProductUseCase {
	return &updateProductUseCase{
		productRepository: productRepository,
		aiProvider:        aiProvider,
	}
}

func (u *updateProductUseCase) Execute(ctx context.Context, userID uuid.UUID, p *domain.UpdateProduct) error {
	sellerid, err := u.productRepository.GetSellerIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	existingProduct, err := u.productRepository.GetProductByID(ctx, p.ProductID)
	if err != nil {
		return err
	}

	// 2. Yetki kontrolü
	if existingProduct.SellerID != sellerid {
		return errors.New("unauthorized to update this product")
	}

	// 3. Değişiklikleri uygula (Sadece nil olmayanları)
	contentChanged := false
	if p.Name != nil && *p.Name != existingProduct.Name {
		existingProduct.Name = *p.Name
		contentChanged = true
	}
	if p.Description != nil && *p.Description != existingProduct.Description {
		existingProduct.Description = *p.Description
		contentChanged = true
	}
	if p.Price != nil {
		existingProduct.Price = *p.Price
	}
	if p.StockCount != nil {
		existingProduct.StockCount = *p.StockCount
	}
	if p.CategoryID != nil {
		existingProduct.CategoryID = *p.CategoryID
	}
	if p.Attributes != nil {
		existingProduct.Attributes = p.Attributes
	}
	err = u.productRepository.UpdateProduct(ctx, existingProduct)
	if err != nil {
		return err
	}
	if contentChanged {
		go func(pID uuid.UUID, name, desc string) {

			text := fmt.Sprintf("%s %s", name, desc)
			vector, err := u.aiProvider.GetVector(text)
			if err != nil {
				fmt.Println("Update AI Error:", err)
				return
			}

			u.productRepository.UpdateProductEmbedding(context.Background(), pID, vector)
		}(existingProduct.ID, existingProduct.Name, existingProduct.Description)
	}

	return nil
}
