// internal/basket-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/basket-service/domain"
)

type Handlers struct {
	basketRepository domain.BasketRepository
}

func NewHandlers(repository domain.BasketRepository) *Handlers {
	return &Handlers{basketRepository: repository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello basket service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
