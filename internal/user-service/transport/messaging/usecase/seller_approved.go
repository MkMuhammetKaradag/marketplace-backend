package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

type SellerApprovedUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}
type sellerApprovedUseCase struct {
	repository domain.UserRepository
}

func NewSellerApprovedUseCase(repository domain.UserRepository) SellerApprovedUseCase {
	return &sellerApprovedUseCase{
		repository: repository,
	}
}

func (u *sellerApprovedUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	fmt.Println("Seller approved use case executed", userID)

	return nil
}
