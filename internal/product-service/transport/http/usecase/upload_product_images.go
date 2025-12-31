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

	for i, file := range files {
		publicID := fmt.Sprintf("%s_%d", productID.String(), i)
		fmt.Println(publicID)
		uploadResult, _, err := c.cloudinarySvc.UploadImage(fiberCtx.UserContext(), file, domain.UploadOptions{
			Folder:         "products",
			PublicID:       publicID,
			Transformation: "c_fill,g_auto,w_1200,h_630,q_auto,f_auto",
			Width:          1200,
			Height:         630,
		})
		fmt.Println(uploadResult)
		if err != nil {
			return fmt.Errorf("image upload failed: %w", err)
		}

		productImages = append(productImages, domain.ProductImage{
			ImageURL:  uploadResult,
			IsMain:    i == 0,
			SortOrder: i,
		})
	}

	return c.productRepository.SaveImagesAndUpdateStatus(fiberCtx.UserContext(), productID, productImages)
}
