package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type PaymentSuccessUseCase interface {
	Execute(ctx context.Context, orderID uuid.UUID) error
}
type paymentSuccessUseCase struct {
	repository domain.ProductRepository
}

func NewPaymentSuccessUseCase(repository domain.ProductRepository) PaymentSuccessUseCase {
	return &paymentSuccessUseCase{
		repository: repository,
	}
}

func (u *paymentSuccessUseCase) Execute(ctx context.Context, orderID uuid.UUID) error {

	err := u.repository.ConfirmStock(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to confirm stock for order %s: %w", orderID, err)
	}
	return nil
}
