package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type PaymentSuccessUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}
type paymentSuccessUseCase struct {
	repository domain.BasketRedisRepository
}

func NewPaymentSuccessUseCase(repository domain.BasketRedisRepository) PaymentSuccessUseCase {
	return &paymentSuccessUseCase{
		repository: repository,
	}
}

func (u *paymentSuccessUseCase) Execute(ctx context.Context, userID uuid.UUID) error {

	err := u.repository.ClearBasket(ctx, userID.String())
	if err != nil {
		return fmt.Errorf("failed to clear basket for user %s: %w", userID, err)
	}
	return nil
}
