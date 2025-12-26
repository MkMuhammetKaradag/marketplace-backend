package usecase

import (
	"context"
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

	uploadRes, publicID, err := u.cloudinarySvc.UploadImage(ctx, fileHeader, domain.UploadOptions{
		Folder:         "store_banners",
		Width:          1200,
		Height:         630,
		PublicID:       userID.String() + "_" + sellerID.String(),
		Transformation: "c_fill,g_auto,w_1200,h_630,q_auto,f_auto",
	})

	if err != nil {
		return err
	}
	err = u.repo.UpdateStoreBanner(ctx, userID, sellerID, uploadRes)
	if err != nil {
		_ = u.cloudinarySvc.DeleteImage(ctx, publicID)
		return err
	}
	return nil
}
