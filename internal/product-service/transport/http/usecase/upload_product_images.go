package usecase

import (
	"fmt"
	"marketplace/internal/product-service/domain"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadProductImagesUseCase interface {
	Execute(fiberCtx *fiber.Ctx, productID uuid.UUID, files []*multipart.FileHeader) error
}

type uploadProductImagesUseCase struct {
	productRepository domain.ProductRepository
	cloudinarySvc     domain.ImageService
}

func NewUploadProductImagesUseCase(productRepository domain.ProductRepository, cloudinarySvc domain.ImageService) UploadProductImagesUseCase {
	return &uploadProductImagesUseCase{
		productRepository: productRepository,
		cloudinarySvc:     cloudinarySvc,
	}
}

func (c *uploadProductImagesUseCase) Execute(fiberCtx *fiber.Ctx, productID uuid.UUID, files []*multipart.FileHeader) error {
	var productImages []domain.ProductImage
	errChan := make(chan error, len(files))
	imgChan := make(chan domain.ProductImage, len(files))

	// Fiber context'inden context'i alıyoruz (fiberCtx bitmeden işimiz bitmeli)
	ctx := fiberCtx.UserContext()

	for i, file := range files {
		// Go routine içinde i ve file değerlerini sabitlemek için kopyalıyoruz
		go func(index int, f *multipart.FileHeader) {
			publicID := fmt.Sprintf("%s_%d", productID.String(), index)

			url, _, err := c.cloudinarySvc.UploadImage(ctx, f, domain.UploadOptions{
				Folder:   "products",
				PublicID: publicID,
				Width:    1200,
				Height:   630,
			})

			if err != nil {
				errChan <- err
				return
			}

			imgChan <- domain.ProductImage{
				ImageURL:  url,
				IsMain:    index == 0,
				SortOrder: index,
			}
			errChan <- nil
		}(i, file)
	}

	// Tüm sonuçları topla
	for i := 0; i < len(files); i++ {
		if err := <-errChan; err != nil {
			return err // Hata varsa hemen dön
		}
		productImages = append(productImages, <-imgChan)
	}

	return c.productRepository.SaveImagesAndUpdateStatus(ctx, productID, productImages)
}
