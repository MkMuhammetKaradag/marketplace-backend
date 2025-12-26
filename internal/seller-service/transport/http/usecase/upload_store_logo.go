package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"

	"mime/multipart"

	"github.com/google/uuid"
)

type UploadStoreLogoUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, fileHeader *multipart.FileHeader) error
}
type uploadStoreLogoUseCase struct {
	repo          domain.SellerRepository
	cloudinarySvc domain.ImageService
}

func NewUploadStoreLogoUseCase(repo domain.SellerRepository, cloudinarySvc domain.ImageService) UploadStoreLogoUseCase {
	return &uploadStoreLogoUseCase{
		repo:          repo,
		cloudinarySvc: cloudinarySvc,
	}
}

func (u *uploadStoreLogoUseCase) Execute(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	uploadRes, publicID, err := u.cloudinarySvc.UploadStoreLogo(ctx, fileHeader, userID.String(), sellerID.String())

	if err != nil {
		return err
	}
	err = u.repo.UpdateStoreLogo(ctx, userID, sellerID, uploadRes)
	if err != nil {
		clErr := u.cloudinarySvc.DeleteImage(ctx, publicID)
		fmt.Println(clErr)
		return err
	}
	fmt.Println(uploadRes)
	return nil
}
