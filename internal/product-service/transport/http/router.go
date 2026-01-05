// internal/user-service/transport/http/router.go
package http

import (
	"fmt"
	"marketplace/internal/product-service/handler"
	"marketplace/internal/product-service/transport/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	//api := app.Group("/api/v1")

	createProduct := r.handlers.CreateProduct()
	uploadProductImages := r.handlers.UploadProductImages()
	createCategory := r.handlers.CreateCategory()
	getRecommendedProducts := r.handlers.GetRecommendedProducts()
	getProduct := r.handlers.GetProduct()

	app.Get("/hello", r.handlers.Hello)
	app.Post("/create", handler.HandleWithFiber[controller.CreateProductRequest, controller.CreateProductResponse](createProduct))
	app.Post("/upload/:product_id", handler.HandleWithFiber[controller.UploadProductImagesRequest, controller.UploadProductImagesResponse](uploadProductImages))
	app.Post("/category", handler.HandleWithFiber[controller.CreateCategoryRequest, controller.CreateCategoryResponse](createCategory))
	app.Get("/recommended", handler.HandleWithFiber[controller.GetRecommendationsRequest, controller.GetRecommendationsResponse](getRecommendedProducts))
	app.Get("/product/:product_id", handler.HandleWithFiber[controller.GetProductRequest, controller.GetProductResponse](getProduct))
	app.Get("/profile", func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {

			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		fmt.Println("user id : " + userIDStr)
		return c.SendString("Hello World " + userIDStr)
	})
}
