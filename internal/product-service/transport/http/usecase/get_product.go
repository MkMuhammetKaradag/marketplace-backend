package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type GetProductUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*domain.Product, error)
}

type getProductUseCase struct {
	productRepository domain.ProductRepository
}

func NewGetProductUseCase(productRepository domain.ProductRepository) GetProductUseCase {
	return &getProductUseCase{
		productRepository: productRepository,
	}
}

func (c *getProductUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*domain.Product, error) {

	product, err := c.productRepository.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	fmt.Println(product)
	// Repository'deki o meşhur vektör benzerliği sorgusunu çağırıyoruz
	go func(p *domain.Product, uID uuid.UUID) {

		if uID == uuid.Nil {
			// Log alabilirsin ama işleme devam etmeye gerek yok
			fmt.Println("Anonim kullanıcı: Tracking atlandı.")
			return
		}
		// Ürünün embedding'i yoksa (AI henüz üretmediyse) bir şey yapamayız
		if len(p.Embedding) == 0 {
			return
		}

		// Repository'deki o meşhur fonksiyonu çağırıyoruz
		// Context.Background kullanıyoruz çünkü ana istek (request) bitse bile bu sürsün
		err := c.productRepository.TrackProductView(context.Background(), uID, p.Embedding)
		if err != nil {
			fmt.Printf("Tracking hatası: %v\n", err)
		}

		// Opsiyonel: user_product_interactions tablosuna da kayıt atabilirsin
		c.productRepository.AddInteraction(context.Background(), uID, p.ID, "view")
	}(product, userID)

	return product, nil
}
