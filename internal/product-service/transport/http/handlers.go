package http

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/controller"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Product  *productHandlers
	Category *categoryHandlers
	Favorite *favoriteHandlers
	Search   *searchHandlers
	General  *generalHandlers
}

type productHandlers struct {
	Create      *controller.CreateProductController
	Update      *controller.UpdateProductController
	Delete      *controller.DeleteProductController
	Get         *controller.GetProductController
	UploadImage *controller.UploadProductImagesController
}

type categoryHandlers struct {
	Create *controller.CreateCategoryController
}

type favoriteHandlers struct {
	Toggle *controller.ToggleFavoriteController
	Get    *controller.GetFavoritesController
}

type searchHandlers struct {
	Search      *controller.SearchProductsController
	Recommended *controller.GetRecommendationsController
}

type generalHandlers struct {
	productService domain.ProductService
}

func NewHandlers(
	ps domain.ProductService,
	repo domain.ProductRepository,
	imgSvc domain.ImageService,
	ai domain.AiProvider,
	wrk domain.Worker,
	msg domain.Messaging,
) *Handlers {
	// Tüm UseCase ve Controller'lar uygulama ayağa kalkarken bir kez oluşturulur.
	return &Handlers{
		Product: &productHandlers{
			Create:      controller.NewCreateProductController(usecase.NewCreateProductUseCase(repo, ai)),
			Update:      controller.NewUpdateProductController(usecase.NewUpdateProductUseCase(repo, ai, msg)),
			Delete:      controller.NewDeleteProductController(usecase.NewDeleteProductUseCase(repo)),
			Get:         controller.NewGetProductController(usecase.NewGetProductUseCase(repo, wrk)),
			UploadImage: controller.NewUploadProductImagesController(usecase.NewUploadProductImagesUseCase(repo, imgSvc, wrk)),
		},
		Category: &categoryHandlers{
			Create: controller.NewCreateCategoryController(usecase.NewCreateCategoryUseCase(repo)),
		},
		Favorite: &favoriteHandlers{
			Toggle: controller.NewToggleFavoriteController(usecase.NewToggleFavoriteUseCase(repo, wrk)),
			Get:    controller.NewGetFavoritesController(usecase.NewGetFavoritesUseCase(repo)),
		},
		Search: &searchHandlers{
			Search:      controller.NewSearchProductsController(usecase.NewSearchProductsUseCase(repo, ai)),
			Recommended: controller.NewGetRecommendationsController(usecase.NewGetRecommendedProductsUseCase(repo)),
		},
		General: &generalHandlers{
			productService: ps,
		},
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": h.General.productService.Greeting(c.UserContext()),
		"info":    "Product Service Handlers connected to domain layer",
	})
}
