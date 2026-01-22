// internal/payment-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello payment service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
