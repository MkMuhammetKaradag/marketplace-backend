package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type GetProductUseCase interface {
	Execute(ctx context.Context, productID uuid.UUID, userID *uuid.UUID) (*domain.Product, error)
}

type getProductUseCase struct {
	productRepository domain.ProductRepository
	distributorWorker domain.Worker
}

func NewGetProductUseCase(productRepository domain.ProductRepository, distributor domain.Worker) GetProductUseCase {
	return &getProductUseCase{
		productRepository: productRepository,
		distributorWorker: distributor,
	}
}

func (c *getProductUseCase) Execute(ctx context.Context, productID uuid.UUID, userID *uuid.UUID) (*domain.Product, error) {

	product, err := c.productRepository.GetProduct(ctx, productID, userID)
	if err != nil {
		return nil, err
	}

	if userID != nil && len(product.Embedding) > 0 {
		// Hızlıca kuyruğa atıyoruz, kullanıcı beklemiyor
		_ = c.distributorWorker.EnqueueTrackView(domain.TrackProductViewPayload{
			UserID:    *userID,
			Embedding: product.Embedding,
			ProductID: product.ID,
		})
	}

	return product, nil
}
