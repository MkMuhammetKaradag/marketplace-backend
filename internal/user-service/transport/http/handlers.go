// internal/user-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/http/controller"
	"marketplace/internal/user-service/transport/http/usecase"
)

type Handlers struct {
	userService       domain.UserService
	userRepository    domain.UserRepository
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
func (h *Handlers) SignOut() *controller.SignOutController {
	logoutUseCase := usecase.NewSignOutUseCase(h.sessionRepository)
	return controller.NewSignOutController(logoutUseCase)
}
func (h *Handlers) AllSignOut() *controller.AllSignOutController {
	logoutUseCase := usecase.NewAllSignOutUseCase(h.sessionRepository)
	return controller.NewAllSignOutController(logoutUseCase)
}

func (h *Handlers) AddUserRole() *controller.AddUserRolerController {
	addUserRolerUseCase := usecase.NewAddUserRolerUseCase(h.userRepository)
	return controller.NewAddUserRolerController(addUserRolerUseCase)
}
func (h *Handlers) CreateRole() *controller.CreateRoleController {
	createRoleUseCase := usecase.NewCreateRoleUseCase(h.userRepository)
	return controller.NewCreateRoleController(createRoleUseCase)
}

func (h *Handlers) ForgotPassword() *controller.ForgotPasswordController {
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(h.userRepository)
	return controller.NewForgotPasswordController(forgotPasswordUseCase)
}

func (h *Handlers) ResetPassword() *controller.ResetPasswordController {
	resetPasswordUseCase := usecase.NewResetPasswordUseCase(h.userRepository)
	return controller.NewResetPasswordController(resetPasswordUseCase)
}

func (h *Handlers) ChangePassword() *controller.ChangePasswordController {
	changePasswordUseCase := usecase.NewChangePasswordUseCase(h.userRepository)
	return controller.NewChangePasswordController(changePasswordUseCase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
