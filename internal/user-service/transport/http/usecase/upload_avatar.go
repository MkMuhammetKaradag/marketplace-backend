package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/google/uuid"
)

type UploadAvatarUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error
}
type uploadAvatarUseCase struct {
	repo       domain.UserRepository
	cloudinary *cloudinary.Cloudinary
}

func NewUploadAvatarUseCase(repo domain.UserRepository, cloudinary *cloudinary.Cloudinary) UploadAvatarUseCase {
	return &uploadAvatarUseCase{
		repo:       repo,
		cloudinary: cloudinary,
	}
}

func (u *uploadAvatarUseCase) Execute(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	uploadRes, err := u.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "profile_pictures",
		PublicID:       userID.String(),
		Overwrite:      api.Bool(true),                            // Eklemeyi unutma
		Invalidate:     api.Bool(true),                            // Eski resim CDN'den de temizlensin                      // Kullanıcı ID'sini isim yaparsak, her yeni yüklemede eski resmin üzerine yazar
		Transformation: "c_fill,g_face,h_500,w_500,q_auto,f_auto", // Otomatik yüz odaklı 500x500 kare yap
	})
	if err != nil {
		return err
	}
	err = u.repo.UpdateAvatar(ctx, userID, uploadRes.SecureURL)
	if err != nil {
		return err
	}
	fmt.Println(uploadRes)
	return nil
}
