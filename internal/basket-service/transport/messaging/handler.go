package messaginghandler

import (
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/messaging/controller"
	"marketplace/internal/basket-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	basketRepository domain.BasketRedisRepository
}

func NewMessageHandlers(repository domain.BasketRedisRepository) *Handlers {
	return &Handlers{basketRepository: repository}
}

func SetupMessageHandlers(repository domain.BasketRedisRepository) map[pb.MessageType]domain.MessageHandler {

	productPriceUpdatedUseCase := usecase.NewProductPriceUpdatedUseCase(repository)
	productPriceUpdatedHandler := controller.NewProductPriceUpdatedHandler(productPriceUpdatedUseCase)
	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_PRODUCT_PRICE_UPDATED: productPriceUpdatedHandler,
	}
}
