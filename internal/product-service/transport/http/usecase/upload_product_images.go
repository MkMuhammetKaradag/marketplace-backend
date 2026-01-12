package usecase

import (
	"bytes"
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
	imageWorker       domain.Worker
}

func NewUploadProductImagesUseCase(productRepository domain.ProductRepository, cloudinarySvc domain.ImageService, imageWorker domain.Worker) UploadProductImagesUseCase {
	return &uploadProductImagesUseCase{
		productRepository: productRepository,
		cloudinarySvc:     cloudinarySvc,
		imageWorker:       imageWorker,
	}
}

func (c *uploadProductImagesUseCase) Execute(fiberCtx *fiber.Ctx, productID uuid.UUID, files []*multipart.FileHeader) error {
	err := c.productRepository.SoftDeleteAllProductImages(fiberCtx.UserContext(), productID)
	if err != nil {
		return fmt.Errorf("failed to clear old images: %w", err)
	}

	for i, fileHeader := range files {

		file, _ := fileHeader.Open()
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		file.Close()

		payload := domain.UploadImageTaskPayload{
			ProductID: productID,
			ImageData: buf.Bytes(),
			FileName:  fileHeader.Filename,
			IsMain:    i == 0,
			SortOrder: i,
		}

		err := c.imageWorker.EnqueueImageUpload(payload)
		if err != nil {
			return fmt.Errorf("kuyruğa atılamadı: %w", err)
		}
	}

	return nil
}
