package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/grpc_client"

	"github.com/google/uuid"
)

type AddItemUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, p *domain.BasketItem) error
}

type addItemUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewAddItemUseCase(basketRepository domain.BasketRedisRepository) AddItemUseCase {
	return &addItemUseCase{
		basketRepository: basketRepository,
	}
}

func (u *addItemUseCase) Execute(ctx context.Context, userID uuid.UUID, p *domain.BasketItem) error {
	// 1. Ürünü gRPC ile doğrula
	product, err := grpc_client.GetProductForBasket(p.ProductID.String())
	if err != nil {
		return err
	}
	if product == nil { // || !product.IsActive
		return fmt.Errorf("product is not available")
	}

	// 2. Fiyat Güvenliği: Kullanıcının gönderdiği fiyatı değil,
	// Product Service'den gelen orijinal fiyatı set ediyoruz.
	p.Price = product.Price
	p.Name = product.Name // İsim de değişmiş olabilir, güncellemek iyidir.

	// 3. Mevcut sepeti çek
	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}
	if basket == nil {
		basket = &domain.Basket{UserID: userID, Items: []domain.BasketItem{}}
	}

	// 4. Miktar ve Stok Kontrolü
	totalRequestedQuantity := p.Quantity
	foundIndex := -1

	for i, item := range basket.Items {
		if item.ProductID == p.ProductID {
			totalRequestedQuantity += item.Quantity
			foundIndex = i
			break
		}
	}

	// Stok yetersizse hata döndür
	if int32(totalRequestedQuantity) > product.Stock {
		return fmt.Errorf("insufficient stock: requested %d, available %d", totalRequestedQuantity, product.Stock)
	}

	// 5. Sepeti Güncelle
	if foundIndex != -1 {
		basket.Items[foundIndex].Quantity = totalRequestedQuantity
		basket.Items[foundIndex].Price = product.Price // Fiyat güncellenmiş olabilir
	} else {
		basket.Items = append(basket.Items, *p)
	}
	fmt.Println("Basket updated:", basket)

	return u.basketRepository.UpdateBasket(ctx, basket)
}
