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
	aiProvider        domain.AiProvider
}

func NewCreateProductUseCase(productRepository domain.ProductRepository, aiProvider domain.AiProvider) CreateProductUseCase {
	return &createProductUseCase{
		productRepository: productRepository,
		aiProvider:        aiProvider,
	}
}

func (c *createProductUseCase) Execute(ctx context.Context, userID uuid.UUID, req *domain.Product) error {

	sellerID, err := c.productRepository.GetSellerIDByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get seller ID: %w", err)
	}

	req.SellerID = sellerID
	productID, err := c.productRepository.CreateProduct(ctx, req)
	if err != nil {
		return err
	} // 2. Gerçek AI servisini çağır
	go func(p *domain.Product) {
		// Ürünü tanımlayan metin: İsim + Açıklama
		text := fmt.Sprintf("%s %s", p.Name, p.Description)

		// Burası artık simüle değil, HuggingFace'e gidiyor
		vector, err := c.aiProvider.GetVector(text)
		if err != nil {
			fmt.Println("AI Error:", err)
			return
		}
		fmt.Println("Vector:", vector)
		// 3. Veritabanını güncelle
		c.productRepository.UpdateProductEmbedding(context.Background(), productID, vector)
	}(req)

	return nil
}
