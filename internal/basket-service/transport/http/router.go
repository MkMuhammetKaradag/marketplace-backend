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

	addItem := r.handlers.AddItem()
	removeItem := r.handlers.RemoveItem()
	decrementItem := r.handlers.DecrementItem()
	incrementItem := r.handlers.IncrementItem()
	clearBasket := r.handlers.ClearBasket()
	getBasket := r.handlers.GetBasket()
	basketCount := r.handlers.BasketCount()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/add-item", handler.HandleWithFiber[controller.AddItemRequest, controller.AddItemResponse](addItem))
	app.Delete("/remove-item/:product_id", handler.HandleWithFiber[controller.RemoveItemRequest, controller.RemoveItemResponse](removeItem))
	app.Patch("/decrement-item/:product_id", handler.HandleWithFiber[controller.DecrementItemRequest, controller.DecrementItemResponse](decrementItem))
	app.Patch("/increment-item/:product_id", handler.HandleWithFiber[controller.IncrementItemRequest, controller.IncrementItemResponse](incrementItem))
	app.Delete("/clear-basket", handler.HandleWithFiber[controller.ClearBasketRequest, controller.ClearBasketResponse](clearBasket))
	app.Get("/basket", handler.HandleWithFiber[controller.GetBasketRequest, controller.GetBasketResponse](getBasket))
	app.Get("/count", handler.HandleWithFiber[controller.BasketCountRequest, controller.BasketCountResponse](basketCount))
}
