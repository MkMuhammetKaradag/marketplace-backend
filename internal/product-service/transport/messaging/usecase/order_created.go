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

	var reserveItems []domain.OrderItemReserve
	for _, item := range items {
		pID, _ := uuid.Parse(item.ProductId)
		reserveItems = append(reserveItems, domain.OrderItemReserve{
			ProductID: pID,
			Quantity:  int(item.Quantity),
		})
	}

	err := u.repository.ReserveStocks(ctx, orderID, reserveItems)
	if err != nil {
		fmt.Println("The reservation failed, a cancel order event can be triggered:", err)

		return err
	}

	fmt.Printf("Products reserved for order %s.\n", orderID)
	return nil

}
