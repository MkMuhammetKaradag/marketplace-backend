// internal/user-service/transport/http/usecase/all_signout.go
package usecase

import (
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
)

type AllSignOutUseCase interface {
	Execute(ctx *fiber.Ctx) error
}
type allSignOutUseCase struct {
	sessionRepository domain.SessionRepository
}

func NewAllSignOutUseCase(repository domain.SessionRepository) AllSignOutUseCase {
	return &allSignOutUseCase{
		sessionRepository: repository,
	}
}

func (u *allSignOutUseCase) Execute(ctx *fiber.Ctx) error {
	token := ctx.Cookies("Session")
	if err := u.sessionRepository.DeleteUserAllSession(ctx.UserContext(), token); err != nil {
		return err

	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	return nil
}
