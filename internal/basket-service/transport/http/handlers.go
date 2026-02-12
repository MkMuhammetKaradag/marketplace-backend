package http

import (
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/http/controller"
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Basket *basketHandlers
}

type basketHandlers struct {
	AddItem       *controller.AddItemController
	RemoveItem    *controller.RemoveItemController
	DecrementItem *controller.DecrementItemController
	IncrementItem *controller.IncrementItemController
	ClearBasket   *controller.ClearBasketController
	GetBasket     *controller.GetBasketController
	Count         *controller.BasketCountController
}

func NewHandlers(
	postgresRepo domain.BasketPostgresRepository,
	redisRepo domain.BasketRedisRepository,
	grpcProductClient domain.ProductClient,
) *Handlers {

	return &Handlers{
		Basket: &basketHandlers{
			AddItem:       controller.NewAddItemController(usecase.NewAddItemUseCase(redisRepo, grpcProductClient)),
			RemoveItem:    controller.NewRemoveItemController(usecase.NewRemoveItemUseCase(redisRepo)),
			DecrementItem: controller.NewDecrementItemController(usecase.NewDecrementItemUseCase(redisRepo)),
			IncrementItem: controller.NewIncrementItemController(usecase.NewIncrementItemUseCase(redisRepo, grpcProductClient)),
			ClearBasket:   controller.NewClearBasketController(usecase.NewClearBasketUseCase(redisRepo)),
			GetBasket:     controller.NewGetBasketController(usecase.NewGetBasketUseCase(redisRepo)),
			Count:         controller.NewBasketCountController(usecase.NewBasketCountUseCase(redisRepo)),
		},
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "hello basket service",
		"info":    "Fiber handler connected to domain layer",
	})
}
