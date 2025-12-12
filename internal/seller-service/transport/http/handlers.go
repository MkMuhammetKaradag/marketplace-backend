// internal/user-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/seller-service/domain"
)

type Handlers struct {
	sellerRepository domain.SellerRepository
}

func NewHandlers(repository domain.SellerRepository) *Handlers {
	return &Handlers{sellerRepository: repository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hhelu seller service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
