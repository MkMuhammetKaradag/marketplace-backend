// internal/user-service/transport/http/router.go
package http

import (
	"fmt"
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
	signInHandler := r.handlers.SignIn()
	signOutHandler := r.handlers.SignOut()
	allSignOutHandler := r.handlers.AllSignOut()
	changeUserRoleHandler := r.handlers.ChangeUserRole()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/signup", handler.HandleBasic[controller.SignUpRequest, controller.SignUpResponse](siginUpHandler))
	app.Post("/user-activate", handler.HandleBasic[controller.UserActivateRequest, controller.UserActivateResponse](userActivateHandler))
	app.Post("/signin", handler.HandleWithFiber[controller.SignInRequest, controller.SignInResponse](signInHandler))
	app.Post("/signout", handler.HandleWithFiber[controller.SignOutRequest, controller.SignOutResponse](signOutHandler))
	app.Post("/all-signout", handler.HandleWithFiber[controller.AllSignOutRequest, controller.AllSignOutResponse](allSignOutHandler))
	app.Post("/change-user-role/:user_id", handler.HandleWithFiber[controller.ChangeUserRolerRequest, controller.ChangeUserRolerResponse](changeUserRoleHandler))
	app.Get("/profile", func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {

			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		fmt.Println("user id : " + userIDStr)
		return c.SendString("Hello World " + userIDStr)
	})
}
