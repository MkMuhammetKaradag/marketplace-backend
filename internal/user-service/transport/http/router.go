package http

import (
	"marketplace/internal/user-service/handler"
	"marketplace/internal/user-service/transport/http/controller"

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
	siginUpHandler := r.handlers.SignUp()
	userActivateHandler := r.handlers.UserActivate()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/signup", handler.HandleBasic[controller.SignUpRequest, controller.SignUpResponse](siginUpHandler))
	app.Post("/user-activate", handler.HandleBasic[controller.UserActivateRequest, controller.UserActivateResponse](userActivateHandler))
}
