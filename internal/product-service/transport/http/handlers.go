// internal/user-service/transport/http/handlers.go
package http

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/controller"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	userService       domain.ProductService
	productRepository domain.ProductRepository
	cloudinarySvc     domain.ImageService
	aiProvider        domain.AiProvider
	worker            domain.Worker
}

func NewHandlers(userService domain.ProductService, repository domain.ProductRepository, cloudinarySvc domain.ImageService, aiProvider domain.AiProvider, worker domain.Worker) *Handlers {
	return &Handlers{userService: userService, productRepository: repository, cloudinarySvc: cloudinarySvc, aiProvider: aiProvider, worker: worker}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: h.userService.Greeting(c.UserContext()),
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) CreateProduct() *controller.CreateProductController {
	usecase := usecase.NewCreateProductUseCase(h.productRepository, h.aiProvider)
	return controller.NewCreateProductController(usecase)
}

func (h *Handlers) UploadProductImages() *controller.UploadProductImagesController {
	usecase := usecase.NewUploadProductImagesUseCase(h.productRepository, h.cloudinarySvc, h.worker)
	return controller.NewUploadProductImagesController(usecase)
}

func (h *Handlers) CreateCategory() *controller.CreateCategoryController {
	usecase := usecase.NewCreateCategoryUseCase(h.productRepository)
	return controller.NewCreateCategoryController(usecase)
}

func (h *Handlers) GetRecommendedProducts() *controller.GetRecommendationsController {
	usecase := usecase.NewGetRecommendedProductsUseCase(h.productRepository)
	return controller.NewGetRecommendationsController(usecase)
}

func (h *Handlers) GetProduct() *controller.GetProductController {
	usecase := usecase.NewGetProductUseCase(h.productRepository)
	return controller.NewGetProductController(usecase)
}

func (h *Handlers) SearchProducts() *controller.SearchProductsController {
	usecase := usecase.NewSearchProductsUseCase(h.productRepository, h.aiProvider)
	return controller.NewSearchProductsController(usecase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
