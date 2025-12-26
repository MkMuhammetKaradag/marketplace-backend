// internal/user-service/infrastructure/cloudinary_service.go
package infrastructure

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudName, apiKey, apiSecret string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{client: cld}, nil
}

func (s *CloudinaryService) UploadStoreLogo(ctx context.Context, fileHeader *multipart.FileHeader, userID string, sellerID string) (string, string, error) {
	file, _ := fileHeader.Open()
	defer file.Close()

	uploadRes, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "store_logos",
		PublicID:       userID + "_" + sellerID,
		Overwrite:      api.Bool(true),
		Invalidate:     api.Bool(true),
		Transformation: "c_fill,h_250,w_250,q_auto,f_auto",
	})

	if err != nil {
		return "", "", fmt.Errorf("cloudinary upload store logo: %w", err)
	}
	return uploadRes.SecureURL, uploadRes.PublicID, nil
}

func (s *CloudinaryService) UploadStoreBanner(ctx context.Context, fileHeader *multipart.FileHeader, userID string, sellerID string) (string, string, error) {
	file, _ := fileHeader.Open()
	defer file.Close()

	uploadRes, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "store_banners",
		PublicID:       userID + "_" + sellerID,
		Overwrite:      api.Bool(true),
		Invalidate:     api.Bool(true),
		Transformation: "c_fill,h_250,w_250,q_auto,f_auto",
	})

	if err != nil {
		return "", "", fmt.Errorf("cloudinary upload store banner: %w", err)
	}
	return uploadRes.SecureURL, uploadRes.PublicID, nil
}

func (s *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
