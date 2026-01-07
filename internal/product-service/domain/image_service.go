package domain

import (
	"context"
	"mime/multipart"
)

type UploadOptions struct {
	Folder         string
	Width          int
	Height         int
	PublicID       string
	Transformation string
}
type ImageService interface {
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, opts UploadOptions) (string, string, error)
	UploadImageFromBytes(ctx context.Context, data []byte, opts UploadOptions) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
}
