// internal/notification-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/notification-service/domain"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	messaging domain.Messaging
}

func NewHandlers(messaging domain.Messaging) *Handlers {
	return &Handlers{messaging: messaging}
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
