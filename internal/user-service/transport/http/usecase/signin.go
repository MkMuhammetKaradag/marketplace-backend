// internal/user-service/transport/http/usecase/signIn.go
package usecase

import (
	"fmt"
	"marketplace/internal/user-service/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SignInUseCase interface {
	Execute(fiberCtx *fiber.Ctx, identifier, password string) error
}
type signInUseCase struct {
	userRepository    domain.UserRepository
	sessionRepository domain.SessionRepository
}

func NewSignInUseCase(repository domain.UserRepository, sessionRepo domain.SessionRepository) SignInUseCase {
	return &signInUseCase{
		userRepository:    repository,
		sessionRepository: sessionRepo,
	}
}

func (u *signInUseCase) Execute(fiberCtx *fiber.Ctx, identifier, password string) error {

	user, err := u.userRepository.SignIn(fiberCtx.UserContext(), identifier, password)
	if err != nil {
		return err
	}

	sessionToken := uuid.New().String()
	device := fiberCtx.Get("User-Agent")
	ip := fiberCtx.IP()

	userData := &domain.SessionData{
		UserID:    user.ID,
		Device:    device,
		Username:  "boş",
		Ip:        ip,
		CreatedAt: time.Now(),
	}
	if err := u.sessionRepository.CreateSession(fiberCtx.UserContext(), sessionToken, 24*time.Hour, userData); err != nil {
		return err
	}
	fiberCtx.Cookie(&fiber.Cookie{
		Name:     "Session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	fmt.Printf("user signed in: %+v\n", user)

	return nil
}
