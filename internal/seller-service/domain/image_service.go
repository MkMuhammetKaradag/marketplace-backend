package domain

import (
	"context"
	"mime/multipart"
)

type ImageService interface {
	UploadStoreLogo(ctx context.Context, file *multipart.FileHeader, userID string, sellerID string) (string, string, error)
	DeleteImage(ctx context.Context, publicID string) error
}
