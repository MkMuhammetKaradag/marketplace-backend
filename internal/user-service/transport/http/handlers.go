// internal/user-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/http/controller"
	"marketplace/internal/user-service/transport/http/usecase"
)

type Handlers struct {
	userService    domain.UserService
	userRepository domain.UserRepository
}

func NewHandlers(userService domain.UserService, repository domain.UserRepository) *Handlers {
	return &Handlers{userService: userService, userRepository: repository}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: h.userService.Greeting(c.UserContext()),
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}
func (h *Handlers) SignUp() *controller.SignUpController {
	signUpUseCase := usecase.NewSignUpUseCase(h.userRepository)
	return controller.NewSignUpController(signUpUseCase)
}
func (h *Handlers) UserActivate() *controller.UserActivateController {
	userActivateUseCase := usecase.NewUserActivateUseCase(h.userRepository)
	return controller.NewUserActivateController(userActivateUseCase)
}
func (h *Handlers) SignIn() *controller.SignInController {
	userActivateUseCase := usecase.NewSignInUseCase(h.userRepository)
	return controller.NewSignInController(userActivateUseCase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
