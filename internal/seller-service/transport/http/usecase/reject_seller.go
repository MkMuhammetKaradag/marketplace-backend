package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
)

type RejectSellerUseCase interface {
	Execute(ctx context.Context, sellerId, rejectedBy string, reason string) error
}
type rejectSellerUseCase struct {
	sellerRepository domain.SellerRepository
}

func NewRejectSellerUseCase(repository domain.SellerRepository) RejectSellerUseCase {
	return &rejectSellerUseCase{
		sellerRepository: repository,
	}
}

func (u *rejectSellerUseCase) Execute(ctx context.Context, sellerId, rejectedBy string, reason string) error {

	sellerUserId, err := u.sellerRepository.RejectSeller(ctx, sellerId, rejectedBy, reason)
	if err != nil {
		return err
	}
	fmt.Println("seller user id ", sellerUserId)

	return nil
}
