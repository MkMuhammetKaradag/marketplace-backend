// internal/order-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/transport/http/controller"
	"marketplace/internal/order-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	orderRepository   domain.OrderRepository
	grpcProductClient domain.ProductClient
	grpcBasketClient  domain.BasketClient
	grpcPaymentClient domain.PaymentClient
}

func NewHandlers(repo domain.OrderRepository, grpcProductClient domain.ProductClient, grpcBasketClient domain.BasketClient, grpcPaymentClient domain.PaymentClient) *Handlers {
	return &Handlers{
		orderRepository:   repo,
		grpcProductClient: grpcProductClient,
		grpcBasketClient:  grpcBasketClient,
		grpcPaymentClient: grpcPaymentClient,
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello order service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) CreateOrder() *controller.CreateOrderController {
	usecase := usecase.NewCreateOrderUseCase(h.orderRepository, h.grpcProductClient, h.grpcBasketClient, h.grpcPaymentClient)
	return controller.NewCreateOrderController(usecase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
