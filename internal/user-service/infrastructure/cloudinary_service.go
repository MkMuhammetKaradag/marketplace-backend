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

func (s *CloudinaryService) UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error) {

	uploadRes, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "profile_pictures",
		PublicID:       userID,
		Overwrite:      api.Bool(true), 
		Invalidate:     api.Bool(true),
		Transformation: "c_fill,g_face,h_500,w_500,q_auto,f_auto",
	})

	if err != nil {
		return "", fmt.Errorf("cloudinary upload: %w", err)
	}
	return uploadRes.SecureURL, nil
}
