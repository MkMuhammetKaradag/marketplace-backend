// internal/seller-service/transport/http/usecase/seller_onboard.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/util"
)

type CreateSellerUseCase interface {
	Execute(ctx context.Context, seller *domain.Seller) (string, error)
}
type createSellerUseCase struct {
	sellerRepository domain.SellerRepository
}

func NewCreateSellerUseCase(repository domain.SellerRepository) CreateSellerUseCase {
	return &createSellerUseCase{
		sellerRepository: repository,
	}
}

func (u *createSellerUseCase) Execute(ctx context.Context, seller *domain.Seller) (string, error) {

	seller.StoreSlug = util.Slugify(seller.StoreName)
	if seller.StoreSlug == "" {
		return "", fmt.Errorf("The store name cannot create a valid slug.")
	}

	id, err := u.sellerRepository.Create(ctx, seller)
	if err != nil {
		return "", err
	}

	return id, nil
}
