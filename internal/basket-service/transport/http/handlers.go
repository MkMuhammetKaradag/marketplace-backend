// internal/basket-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/http/controller"
	"marketplace/internal/basket-service/transport/http/usecase"
)

type Handlers struct {
	basketPostgresRepository domain.BasketPostgresRepository
	basketRedisRepository    domain.BasketRedisRepository
	grpcProductClient        domain.ProductClient
}

func NewHandlers(postgresRepository domain.BasketPostgresRepository, redisRepository domain.BasketRedisRepository, grpcProductClient domain.ProductClient) *Handlers {
	return &Handlers{basketPostgresRepository: postgresRepository, basketRedisRepository: redisRepository, grpcProductClient: grpcProductClient}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello basket service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) AddItem() *controller.AddItemController {
	usecase := usecase.NewAddItemUseCase(h.basketRedisRepository, h.grpcProductClient)
	return controller.NewAddItemController(usecase)

}
func (h *Handlers) RemoveItem() *controller.RemoveItemController {
	usecase := usecase.NewRemoveItemUseCase(h.basketRedisRepository)
	return controller.NewRemoveItemController(usecase)
}

func (h *Handlers) DecrementItem() *controller.DecrementItemController {
	usecase := usecase.NewDecrementItemUseCase(h.basketRedisRepository)
	return controller.NewDecrementItemController(usecase)
}

func (h *Handlers) IncrementItem() *controller.IncrementItemController {
	usecase := usecase.NewIncrementItemUseCase(h.basketRedisRepository, h.grpcProductClient)
	return controller.NewIncrementItemController(usecase)
}
func (h *Handlers) ClearBasket() *controller.ClearBasketController {
	usecase := usecase.NewClearBasketUseCase(h.basketRedisRepository)
	return controller.NewClearBasketController(usecase)
}

func (h *Handlers) GetBasket() *controller.GetBasketController {
	usecase := usecase.NewGetBasketUseCase(h.basketRedisRepository)
	return controller.NewGetBasketController(usecase)
}

func (h *Handlers) BasketCount() *controller.BasketCountController {
	usecase := usecase.NewBasketCountUseCase(h.basketRedisRepository)
	return controller.NewBasketCountController(usecase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
