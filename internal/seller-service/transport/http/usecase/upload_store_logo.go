package usecase

import (
	"context"
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

	uploadRes, publicID, err := u.cloudinarySvc.UploadImage(ctx, fileHeader, domain.UploadOptions{
		Folder:         "store_logos",
		Width:          250,
		Height:         250,
		PublicID:       userID.String() + "_" + sellerID.String(),
		Transformation: "c_fill,g_auto,w_250,h_250,q_auto,f_auto",
	})

	if err != nil {
		return err
	}
	err = u.repo.UpdateStoreLogo(ctx, userID, sellerID, uploadRes)
	if err != nil {
		_ = u.cloudinarySvc.DeleteImage(ctx, publicID)
		return err
	}
	return nil
}
