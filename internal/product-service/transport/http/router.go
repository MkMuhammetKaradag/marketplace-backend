// internal/product-service/transport/http/router.go
package http

import (
	"marketplace/internal/product-service/handler"
	"marketplace/internal/product-service/transport/http/controller"

	_ "marketplace/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	//api := app.Group("/api/v1")
	app.Get("/swagger/*", swagger.HandlerDefault)
	createProduct := r.handlers.CreateProduct()
	uploadProductImages := r.handlers.UploadProductImages()
	createCategory := r.handlers.CreateCategory()
	getRecommendedProducts := r.handlers.GetRecommendedProducts()
	getProduct := r.handlers.GetProduct()
	searchProducts := r.handlers.SearchProducts()
	toggleFavorite := r.handlers.ToggleFavorite()
	getUserFavorites := r.handlers.GetUserFavorites()
	updateProduct := r.handlers.UpdateProduct()
	deleteProduct := r.handlers.DeleteProduct()

	app.Get("/hello", r.handlers.Hello)
	app.Post("/create", handler.HandleWithFiber[controller.CreateProductRequest, controller.CreateProductResponse](createProduct))
	app.Post("/upload/:product_id", handler.HandleWithFiber[controller.UploadProductImagesRequest, controller.UploadProductImagesResponse](uploadProductImages))
	app.Post("/category", handler.HandleWithFiber[controller.CreateCategoryRequest, controller.CreateCategoryResponse](createCategory))
	app.Get("/recommended", handler.HandleWithFiber[controller.GetRecommendationsRequest, controller.GetRecommendationsResponse](getRecommendedProducts))
	app.Get("/product/:product_id", handler.HandleWithFiber[controller.GetProductRequest, controller.GetProductResponse](getProduct))
	app.Get("/search", handler.HandleWithFiber[controller.SearchProductsRequest, controller.SearchProductsResponse](searchProducts))
	app.Post("/toggle-favorite/:product_id", handler.HandleWithFiber[controller.ToggleFavoriteRequest, controller.ToggleFavoriteResponse](toggleFavorite))
	app.Get("/favorites", handler.HandleWithFiber[controller.GetFavoritesRequest, controller.GetFavoritesResponse](getUserFavorites))
	app.Put("/update/:product_id", handler.HandleWithFiber[controller.UpdateProductRequest, controller.UpdateProductResponse](updateProduct))
	app.Delete("/delete/:product_id", handler.HandleWithFiber[controller.DeleteProductRequest, controller.DeleteProductResponse](deleteProduct))
}
