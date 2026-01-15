// internal/basket-service/transport/http/router.go
package http

import (
	"marketplace/internal/basket-service/handler"
	"marketplace/internal/basket-service/transport/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	//api := app.Group("/api/v1")

	addItem := r.handlers.AddItem()
	removeItem := r.handlers.RemoveItem()
	decrementItem := r.handlers.DecrementItem()
	incrementItem := r.handlers.IncrementItem()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/add-item", handler.HandleWithFiber[controller.AddItemRequest, controller.AddItemResponse](addItem))
	app.Delete("/remove-item/:product_id", handler.HandleWithFiber[controller.RemoveItemRequest, controller.RemoveItemResponse](removeItem))
	app.Patch("/decrement-item/:product_id", handler.HandleWithFiber[controller.DecrementItemRequest, controller.DecrementItemResponse](decrementItem))
	app.Patch("/increment-item/:product_id", handler.HandleWithFiber[controller.IncrementItemRequest, controller.IncrementItemResponse](incrementItem))
}
