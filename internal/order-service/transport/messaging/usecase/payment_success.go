package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

type PaymentSuccessUseCase interface {
	Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, amount float64, stripeSessionID string) error
}
type paymentSuccessUseCase struct {
	repository domain.OrderRepository
}

func NewPaymentSuccessUseCase(repository domain.OrderRepository) PaymentSuccessUseCase {
	return &paymentSuccessUseCase{
		repository: repository,
	}
}

func (u *paymentSuccessUseCase) Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, amount float64, stripeSessionID string) error {

	fmt.Println("PaymentSuccessUseCase Execute", orderID, userID, amount, stripeSessionID)
	return nil
}
