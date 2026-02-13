package messaginghandler

import (
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/messaging/controller"
	"marketplace/internal/basket-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	ProductPriceUpdated domain.MessageHandler
	PaymentSuccess      domain.MessageHandler
}

func NewHandlers(repository domain.BasketRedisRepository) *Handlers {
	return &Handlers{
		ProductPriceUpdated: controller.NewProductPriceUpdatedHandler(
			usecase.NewProductPriceUpdatedUseCase(repository),
		),
		PaymentSuccess: controller.NewPaymentSuccessHandler(
			usecase.NewPaymentSuccessUseCase(repository),
		),
	}
}

func SetupMessageHandlers(repository domain.BasketRedisRepository) map[pb.MessageType]domain.MessageHandler {
	h := NewHandlers(repository)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_PRODUCT_PRICE_UPDATED: h.ProductPriceUpdated,
		pb.MessageType_PAYMENT_SUCCESSFUL:    h.PaymentSuccess,
	}
}
