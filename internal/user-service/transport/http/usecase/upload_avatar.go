package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	"mime/multipart"

	"github.com/google/uuid"
)

type UploadAvatarUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error
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

func (u *uploadAvatarUseCase) Execute(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	uploadRes, err := u.cloudinary.UploadAvatar(ctx, fileHeader, userID.String())

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
