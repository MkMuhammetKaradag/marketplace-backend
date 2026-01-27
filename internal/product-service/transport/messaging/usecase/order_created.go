package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"
	eventsProto "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type OrderCreatedUseCase interface {
	Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, items []*eventsProto.OrderItemData) error
}
type orderCreatedUseCase struct {
	repository domain.ProductRepository
}

func NewOrderCreatedUseCase(repository domain.ProductRepository) OrderCreatedUseCase {
	return &orderCreatedUseCase{
		repository: repository,
	}
}

func (u *orderCreatedUseCase) Execute(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, items []*eventsProto.OrderItemData) error {

	fmt.Println("OrderCreatedUseCase Execute", orderID, userID, items)

	return nil
}
