package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	"mime/multipart"

	"github.com/google/uuid"
)

type UploadAvatarUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, file multipart.File) error
}
type uploadAvatarUseCase struct {
	repo       domain.UserRepository
	cloudinary domain.ImageService
}

func NewUploadAvatarUseCase(repo domain.UserRepository, cloudinary domain.ImageService) UploadAvatarUseCase {
	return &uploadAvatarUseCase{
		repo:       repo,
		cloudinary: cloudinary,
	}
}

func (u *uploadAvatarUseCase) Execute(ctx context.Context, userID uuid.UUID, file multipart.File) error {

	uploadRes, err := u.cloudinary.UploadAvatar(ctx, file, userID.String())

	if err != nil {
		return err
	}
	err = u.repo.UpdateAvatar(ctx, userID, uploadRes)
	if err != nil {
		return err
	}
	fmt.Println(uploadRes)
	return nil
}
