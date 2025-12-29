// internal/user-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/product-service/domain"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	userService    domain.ProductService
	userRepository domain.ProductRepository
}

func NewHandlers(userService domain.ProductService, repository domain.ProductRepository) *Handlers {
	return &Handlers{userService: userService, userRepository: repository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: h.userService.Greeting(c.UserContext()),
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
