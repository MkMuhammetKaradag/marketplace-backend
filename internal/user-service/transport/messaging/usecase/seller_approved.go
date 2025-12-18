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

	err := u.repository.ChangeUserRole(ctx, userID, domain.Seller)
	if err != nil {
		return fmt.Errorf("failed to change user role: %w", err)
	}

	return nil
}
