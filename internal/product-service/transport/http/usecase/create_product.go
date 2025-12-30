package usecase

import (
	"marketplace/internal/product-service/domain"

	"github.com/gofiber/fiber/v2"
)

type CreateProductUseCase interface {
	Execute(fiberCtx *fiber.Ctx, req *domain.Product) error
}

type createProductUseCase struct {
	productRepository domain.ProductRepository
}

func NewCreateProductUseCase(productRepository domain.ProductRepository) CreateProductUseCase {
	return &createProductUseCase{
		productRepository: productRepository,
	}
}

func (c *createProductUseCase) Execute(fiberCtx *fiber.Ctx, req *domain.Product) error {
	return c.productRepository.CreateProduct(fiberCtx.UserContext(), req)
}
