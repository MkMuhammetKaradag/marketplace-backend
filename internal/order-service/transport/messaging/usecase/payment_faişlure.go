package usecase

import (
	"context"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

type PaymentFailureUseCase interface {
	Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, errorMessage string) error
}
type paymentFailureUseCase struct {
	repository domain.OrderRepository
}

func NewPaymentFailureUseCase(repository domain.OrderRepository) PaymentFailureUseCase {
	return &paymentFailureUseCase{
		repository: repository,
	}
}

func (u *paymentFailureUseCase) Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, errorMessage string) error {

	err := u.repository.UpdateOrderStatus(ctx, orderID, domain.OrderFailed)
	if err != nil {
		return err
	}
	return nil
}
