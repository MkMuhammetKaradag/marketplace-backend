package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type SellerApprovedUseCase interface {
	Execute(ctx context.Context, sellerID uuid.UUID) error
}
type sellerApprovedUseCase struct {
	repository domain.ProductRepository
}

func NewSellerApprovedUseCase(repository domain.ProductRepository) SellerApprovedUseCase {
	return &sellerApprovedUseCase{
		repository: repository,
	}
}

func (u *sellerApprovedUseCase) Execute(ctx context.Context, sellerID uuid.UUID) error {

	err := u.repository.AddSeller(ctx, sellerID, "approved")
	if err != nil {
		return fmt.Errorf("failed to add seller: %w", err)
	}

	return nil
}
