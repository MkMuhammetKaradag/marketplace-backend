package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
)

type ApproveSellerUseCase interface {
	Execute(ctx context.Context, sellerId, approveBy string) error
}
type approveSellerUseCase struct {
	sellerRepository domain.SellerRepository
}

func NewApproveSellerUseCase(repository domain.SellerRepository) ApproveSellerUseCase {
	return &approveSellerUseCase{
		sellerRepository: repository,
	}
}

func (u *approveSellerUseCase) Execute(ctx context.Context, sellerId, approveBy string) error {

	sellerUserId, err := u.sellerRepository.ApproveSeller(ctx, sellerId, approveBy)
	if err != nil {
		return err
	}
	fmt.Println("seller user id ", sellerUserId)

	return nil
}
