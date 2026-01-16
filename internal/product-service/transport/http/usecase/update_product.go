package usecase

import (
	"context"
	"errors"
	"fmt"
	"marketplace/internal/product-service/domain"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type UpdateProductUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, req *domain.UpdateProduct) error
}

type updateProductUseCase struct {
	productRepository domain.ProductRepository
	aiProvider        domain.AiProvider
	messaging         domain.Messaging
}

func NewUpdateProductUseCase(productRepository domain.ProductRepository, aiProvider domain.AiProvider, messaging domain.Messaging) UpdateProductUseCase {
	return &updateProductUseCase{
		productRepository: productRepository,
		aiProvider:        aiProvider,
		messaging:         messaging,
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

	if p.Price != nil {
		go func(pID uuid.UUID, price float64) {
			u.messaging.PublishMessage(context.Background(), &pb.Message{
				Type:        pb.MessageType_PRODUCT_PRICE_UPDATED,
				FromService: pb.ServiceType_PRODUCT_SERVICE,
				Critical:    false,
				RetryCount:  2,
				ToServices:  []pb.ServiceType{pb.ServiceType_BASKET_SERVICE},
				Payload: &pb.Message_ProductPriceUpdatedData{ProductPriceUpdatedData: &pb.ProductPriceUpdatedData{
					ProductId: pID.String(),
					Price:     float32(price),
				}},
			})
		}(existingProduct.ID, *p.Price)
	}

	return nil
}
