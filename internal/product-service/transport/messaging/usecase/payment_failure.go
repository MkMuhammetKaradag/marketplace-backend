package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type PaymentFailureUseCase interface {
	Execute(ctx context.Context, orderID uuid.UUID, errorMessage string) error
}
type paymentFailureUseCase struct {
	repository domain.ProductRepository
}

func NewPaymentFailureUseCase(repository domain.ProductRepository) PaymentFailureUseCase {
	return &paymentFailureUseCase{
		repository: repository,
	}
}

func (u *paymentFailureUseCase) Execute(ctx context.Context, orderID uuid.UUID, errorMessage string) error {

	err := u.repository.ReleaseStock(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to confirm stock for order %s: %w", orderID, err)
	}
	return nil
}
