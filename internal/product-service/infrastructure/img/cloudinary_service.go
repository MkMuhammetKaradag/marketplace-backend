// internal/product-service/infrastructure/img/cloudinary_service.go
package img

import (
	"bytes"
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

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

func (s *CloudinaryService) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, opts domain.UploadOptions) (string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// Eğer özel bir transformation gelmediyse varsayılanı kullan
	if opts.Transformation == "" {
		opts.Transformation = fmt.Sprintf("c_fill,g_auto,w_%d,h_%d,q_auto,f_auto", opts.Width, opts.Height)
	}

	uploadRes, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         opts.Folder,
		PublicID:       opts.PublicID,
		Overwrite:      api.Bool(true),
		Invalidate:     api.Bool(true),
		Transformation: opts.Transformation,
	})

	if err != nil {
		return "", "", err
	}
	return uploadRes.SecureURL, uploadRes.PublicID, nil
}
func (s *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
func (s *CloudinaryService) UploadImageFromBytes(ctx context.Context, data []byte, opts domain.UploadOptions) (string, error) {

	reader := bytes.NewReader(data)

	uploadRes, err := s.client.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:         opts.Folder,
		PublicID:       opts.PublicID,
		Transformation: opts.Transformation,
	})

	if err != nil {
		return "", err
	}
	return uploadRes.SecureURL, nil
}
