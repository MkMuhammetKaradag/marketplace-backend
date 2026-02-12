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
	h := r.handlers

	// Public / Test Route
	app.Get("/hello", h.Hello)

	// Basket Group

	basket := app.Group("/")
	{
		// Ekleme ve Güncelleme
		basket.Post("/add-item", handler.HandleWithFiber[controller.AddItemRequest, controller.AddItemResponse](h.Basket.AddItem))
		basket.Patch("/increment-item/:product_id", handler.HandleWithFiber[controller.IncrementItemRequest, controller.IncrementItemResponse](h.Basket.IncrementItem))
		basket.Patch("/decrement-item/:product_id", handler.HandleWithFiber[controller.DecrementItemRequest, controller.DecrementItemResponse](h.Basket.DecrementItem))

		// Görüntüleme
		basket.Get("/basket", handler.HandleWithFiber[controller.GetBasketRequest, controller.GetBasketResponse](h.Basket.GetBasket))
		basket.Get("/count", handler.HandleWithFiber[controller.BasketCountRequest, controller.BasketCountResponse](h.Basket.Count))

		// Silme
		basket.Delete("/remove-item/:product_id", handler.HandleWithFiber[controller.RemoveItemRequest, controller.RemoveItemResponse](h.Basket.RemoveItem))
		basket.Delete("/clear-basket", handler.HandleWithFiber[controller.ClearBasketRequest, controller.ClearBasketResponse](h.Basket.ClearBasket))
	}
}
