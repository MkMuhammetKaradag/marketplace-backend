// internal/user-service/transport/http/router.go
package http

import (
	"marketplace/internal/seller-service/handler"
	"marketplace/internal/seller-service/transport/http/controller"

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
	createSellerHandler := r.handlers.CreateSeller()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/onboard", handler.HandleWithFiber[controller.CreateSellerRequest, controller.CreateSellerResponse](createSellerHandler))

}
