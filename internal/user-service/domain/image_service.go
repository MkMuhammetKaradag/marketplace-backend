package domain

import (
	"context"
	"mime/multipart"
)

type ImageService interface {
	UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error)
}
