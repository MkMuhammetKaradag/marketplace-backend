package http

import (
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/http/controller"
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	// Artık repo veya servisleri değil, direkt hazır Controller'ları tutuyoruz
	Auth    *authHandlers
	User    *userHandlers
	Role    *roleHandlers
	General *generalHandlers
}

type authHandlers struct {
	SignUp         *controller.SignUpController
	SignIn         *controller.SignInController
	SignOut        *controller.SignOutController
	AllSignOut     *controller.AllSignOutController
	UserActivate   *controller.UserActivateController
	ForgotPassword *controller.ForgotPasswordController
	ResetPassword  *controller.ResetPasswordController
	ChangePassword *controller.ChangePasswordController
}

type userHandlers struct {
	UploadAvatar *controller.UploadAvatarController
}

type roleHandlers struct {
	CreateRole *controller.CreateRoleController
	AddRole    *controller.AddUserRolerController
}

type generalHandlers struct {
	userService domain.UserService
}

func NewHandlers(
	userService domain.UserService,
	repository domain.UserRepository,
	sessionRepo domain.SessionRepository,
	messaging domain.Messaging,
	cloudinary domain.ImageService,
) *Handlers {

	return &Handlers{
		Auth: &authHandlers{
			SignUp:         controller.NewSignUpController(usecase.NewSignUpUseCase(repository, messaging)),
			SignIn:         controller.NewSignInController(usecase.NewSignInUseCase(repository, sessionRepo)),
			SignOut:        controller.NewSignOutController(usecase.NewSignOutUseCase(sessionRepo)),
			AllSignOut:     controller.NewAllSignOutController(usecase.NewAllSignOutUseCase(sessionRepo)),
			UserActivate:   controller.NewUserActivateController(usecase.NewUserActivateUseCase(repository, messaging)),
			ForgotPassword: controller.NewForgotPasswordController(usecase.NewForgotPasswordUseCase(repository, messaging)),
			ResetPassword:  controller.NewResetPasswordController(usecase.NewResetPasswordUseCase(repository)),
			ChangePassword: controller.NewChangePasswordController(usecase.NewChangePasswordUseCase(repository, sessionRepo)),
		},
		User: &userHandlers{
			UploadAvatar: controller.NewUploadAvatarController(usecase.NewUploadAvatarUseCase(repository, cloudinary)),
		},
		Role: &roleHandlers{
			CreateRole: controller.NewCreateRoleController(usecase.NewCreateRoleUseCase(repository)),
			AddRole:    controller.NewAddUserRolerController(usecase.NewAddUserRolerUseCase(repository)),
		},
		General: &generalHandlers{
			userService: userService,
		},
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": h.General.userService.Greeting(c.UserContext()),
		"info":    "Fiber handler connected to domain layer",
	})
}
