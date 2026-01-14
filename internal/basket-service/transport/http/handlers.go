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
}

func NewHandlers(postgresRepository domain.BasketPostgresRepository, redisRepository domain.BasketRedisRepository) *Handlers {
	return &Handlers{basketPostgresRepository: postgresRepository, basketRedisRepository: redisRepository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello basket service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) AddItem() *controller.AddItemController {
	usecase := usecase.NewAddItemUseCase(h.basketRedisRepository)
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

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
