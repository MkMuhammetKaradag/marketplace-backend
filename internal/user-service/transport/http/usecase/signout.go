// internal/user-service/transport/http/usecase/signout.go
package usecase

import (
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
)

type SignOutUseCase interface {
	Execute(ctx *fiber.Ctx) error
}
type signOutUseCase struct {
	sessionRepository domain.SessionRepository
}

func NewSignOutUseCase(repository domain.SessionRepository) SignOutUseCase {
	return &signOutUseCase{
		sessionRepository: repository,
	}
}

func (u *signOutUseCase) Execute(ctx *fiber.Ctx) error {
	token := ctx.Cookies("Session")
	if err := u.sessionRepository.DeleteSession(ctx.UserContext(), token); err != nil {
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
