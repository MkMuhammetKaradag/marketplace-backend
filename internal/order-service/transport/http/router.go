// internal/order-service/transport/http/router.go
package http

import (
	"marketplace/internal/order-service/handler"
	"marketplace/internal/order-service/transport/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {

	createOrder := r.handlers.CreateOrder()
	getOrdersByUser := r.handlers.GetOrdersByUser()
	app.Post("/order", handler.HandleWithFiber[controller.CreateOrderRequest, controller.CreateOrderResponse](createOrder))
	app.Get("/user", handler.HandleWithFiber[controller.GetOrdersByUserRequest, controller.GetOrdersByUserResponse](getOrdersByUser))
	app.Get("/hello", r.handlers.Hello)
}
