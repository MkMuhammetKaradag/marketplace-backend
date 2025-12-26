package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"

	"mime/multipart"

	"github.com/google/uuid"
)

type UploadStoreBannerUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, fileHeader *multipart.FileHeader) error
}
type uploadStoreBannerUseCase struct {
	repo          domain.SellerRepository
	cloudinarySvc domain.ImageService
}

func NewUploadStoreBannerUseCase(repo domain.SellerRepository, cloudinarySvc domain.ImageService) UploadStoreBannerUseCase {
	return &uploadStoreBannerUseCase{
		repo:          repo,
		cloudinarySvc: cloudinarySvc,
	}
}

func (u *uploadStoreBannerUseCase) Execute(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	uploadRes, publicID, err := u.cloudinarySvc.UploadStoreBanner(ctx, fileHeader, userID.String(), sellerID.String())

	if err != nil {
		return err
	}
	err = u.repo.UpdateStoreBanner(ctx, userID, sellerID, uploadRes)
	if err != nil {
		clErr := u.cloudinarySvc.DeleteImage(ctx, publicID)
		fmt.Println(clErr)
		return err
	}
	fmt.Println(uploadRes)
	return nil
}
