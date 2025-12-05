// internal/user-service/transport/http/usecase/signIn.go
package usecase

import (
	"fmt"
	"marketplace/internal/user-service/domain"

	"github.com/gofiber/fiber/v2"
)

type SignInUseCase interface {
	Execute(fiberCtx *fiber.Ctx, identifier, password string) error
}
type signInUseCase struct {
	userRepository domain.UserRepository
}

func NewSignInUseCase(repository domain.UserRepository) SignInUseCase {
	return &signInUseCase{
		userRepository: repository,
	}
}

func (u *signInUseCase) Execute(fiberCtx *fiber.Ctx, identifier, password string) error {

	user, err := u.userRepository.SignIn(fiberCtx.UserContext(), identifier, password)
	if err != nil {
		return err
	}
	fiberCtx.Cookie(&fiber.Cookie{
		Name:     "Session",
		Value:    user.ID,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	fmt.Printf("user signed in: %+v\n", user)

	return nil
}
