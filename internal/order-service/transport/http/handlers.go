// internal/order-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/order-service/domain"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	orderRepository domain.OrderRepository
}

func NewHandlers(repo domain.OrderRepository) *Handlers {
	return &Handlers{
		orderRepository: repo,
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello order service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
