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
	h := r.handlers

	// Public Routes
	app.Get("/hello", h.Hello)

	// Auth Group
	auth := app.Group("/")
	{
		auth.Post("/signup", handler.HandleBasic[controller.SignUpRequest, controller.SignUpResponse](h.Auth.SignUp))
		auth.Post("/user-activate", handler.HandleBasic[controller.UserActivateRequest, controller.UserActivateResponse](h.Auth.UserActivate))
		auth.Post("/signin", handler.HandleWithFiber[controller.SignInRequest, controller.SignInResponse](h.Auth.SignIn))
		auth.Post("/signout", handler.HandleWithFiber[controller.SignOutRequest, controller.SignOutResponse](h.Auth.SignOut))
		auth.Post("/all-signout", handler.HandleWithFiber[controller.AllSignOutRequest, controller.AllSignOutResponse](h.Auth.AllSignOut))
		auth.Post("/forgot-password", handler.HandleBasic[controller.ForgotPasswordRequest, controller.ForgotPasswordResponse](h.Auth.ForgotPassword))
		auth.Post("/reset-password", handler.HandleBasic[controller.ResetPasswordRequest, controller.ResetPasswordResponse](h.Auth.ResetPassword))
		auth.Post("/change-password", handler.HandleWithFiber[controller.ChangePasswordRequest, controller.ChangePasswordResponse](h.Auth.ChangePassword))
	}

	// Role Management Group
	roles := app.Group("/roles")
	{
		roles.Post("/create", handler.HandleWithFiber[controller.CreateRoleRequest, controller.CreateRoleResponse](h.Role.CreateRole))
		roles.Post("/assign/:user_id", handler.HandleWithFiber[controller.AddUserRolerRequest, controller.AddUserRolerResponse](h.Role.AddRole))
	}

	// User Profile Group
	user := app.Group("/user")
	{
		user.Post("/upload-avatar", handler.HandleWithFiber[controller.UploadAvatarRequest, controller.UploadAvatarResponse](h.User.UploadAvatar))
		user.Get("/profile", r.profilePlaceholder)
	}
}

func (r *Router) profilePlaceholder(c *fiber.Ctx) error {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	fmt.Println("user id : " + userIDStr)
	return c.SendString("Hello World " + userIDStr)
}
