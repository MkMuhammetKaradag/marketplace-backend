// internal/user-service/transport/http/router.go
package http

import (
	"fmt"

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

	app.Get("/hello", r.handlers.Hello)

	app.Get("/profile", func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {

			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		fmt.Println("user id : " + userIDStr)
		return c.SendString("Hello World " + userIDStr)
	})
}
