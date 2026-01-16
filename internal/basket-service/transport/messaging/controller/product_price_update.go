package controller

import (
	"context"

	"fmt"

	"marketplace/internal/basket-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"
)

type ProductPriceUpdatedHandler struct {
	usecase usecase.ProductPriceUpdatedUseCase
}

func NewProductPriceUpdatedHandler(usecase usecase.ProductPriceUpdatedUseCase) *ProductPriceUpdatedHandler {
	return &ProductPriceUpdatedHandler{
		usecase: usecase,
	}
}

func (h *ProductPriceUpdatedHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetProductPriceUpdatedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not ProductPriceUpdatedData for message ID: %s", msg.Id)
	}

	return h.usecase.Execute(ctx, data.ProductId, float64(data.Price))
}
