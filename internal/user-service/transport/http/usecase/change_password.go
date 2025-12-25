package usecase

import (
	"context"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string, closeAllSessions bool) error
}
type changePasswordUseCase struct {
	repo domain.UserRepository
	sess domain.SessionRepository
}

func NewChangePasswordUseCase(repo domain.UserRepository, sess domain.SessionRepository) ChangePasswordUseCase {
	return &changePasswordUseCase{
		repo: repo,
		sess: sess,
	}
}

func (u *changePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string, closeAllSessions bool) error {

	err := u.repo.ChangePassword(ctx, userID, oldPassword, newPassword)
	if err != nil {
		return err
	}

	if closeAllSessions {
		err = u.sess.DeleteUserAllSession(ctx, userID.String())
		if err != nil {
			return err
		}
	}
	return nil
}
