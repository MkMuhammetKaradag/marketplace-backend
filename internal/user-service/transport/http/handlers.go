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
	sessionRepository domain.SessionRepository
}

func NewHandlers(userService domain.UserService, repository domain.UserRepository, sessionRepo domain.SessionRepository) *Handlers {
	return &Handlers{userService: userService, userRepository: repository, sessionRepository: sessionRepo}
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
	userActivateUseCase := usecase.NewSignInUseCase(h.userRepository, h.sessionRepository)
	return controller.NewSignInController(userActivateUseCase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
