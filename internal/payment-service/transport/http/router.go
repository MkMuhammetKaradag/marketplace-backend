// internal/payment-service/transport/http/router.go
package http

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {

	app.Get("/hello", r.handlers.Hello)
}
