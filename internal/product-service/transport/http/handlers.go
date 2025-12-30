// internal/user-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/controller"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	userService       domain.ProductService
	productRepository domain.ProductRepository
}

func NewHandlers(userService domain.ProductService, repository domain.ProductRepository) *Handlers {
	return &Handlers{userService: userService, productRepository: repository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: h.userService.Greeting(c.UserContext()),
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) CreateProduct() *controller.CreateProductController {
	usecase := usecase.NewCreateProductUseCase(h.productRepository)
	return controller.NewCreateProductController(usecase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
