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
	h := r.handlers

	// Swagger & Health Check
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/hello", h.Hello)

	// --- PRODUCT GROUP ---
	// Poliçe listendeki "/products/..." rotalarıyla uyumlu
	products := app.Group("/")
	{
		products.Post("/create", handler.HandleWithFiber[controller.CreateProductRequest, controller.CreateProductResponse](h.Product.Create))
		products.Post("/upload/:product_id", handler.HandleWithFiber[controller.UploadProductImagesRequest, controller.UploadProductImagesResponse](h.Product.UploadImage))
		products.Post("/category", handler.HandleWithFiber[controller.CreateCategoryRequest, controller.CreateCategoryResponse](h.Category.Create))
		products.Put("/update/:product_id", handler.HandleWithFiber[controller.UpdateProductRequest, controller.UpdateProductResponse](h.Product.Update))
		products.Delete("/delete/:product_id", handler.HandleWithFiber[controller.DeleteProductRequest, controller.DeleteProductResponse](h.Product.Delete))

		// Tekil ürün görüntüleme
		products.Get("/product/:product_id", handler.HandleWithFiber[controller.GetProductRequest, controller.GetProductResponse](h.Product.Get))
	}

	// --- SEARCH & DISCOVERY ---
	// Arama ve öneri rotaları
	search := app.Group("/") // Genelde search işlemleri de ürünlerin altındadır
	{
		search.Get("/recommended", handler.HandleWithFiber[controller.GetRecommendationsRequest, controller.GetRecommendationsResponse](h.Search.Recommended))
		search.Get("/search", handler.HandleWithFiber[controller.SearchProductsRequest, controller.SearchProductsResponse](h.Search.Search))
	}

	// --- FAVORITES ---
	// Kullanıcı bazlı favori işlemleri
	favorites := app.Group("/")
	{
		favorites.Post("/toggle-favorite/:product_id", handler.HandleWithFiber[controller.ToggleFavoriteRequest, controller.ToggleFavoriteResponse](h.Favorite.Toggle))
		favorites.Get("/favorites", handler.HandleWithFiber[controller.GetFavoritesRequest, controller.GetFavoritesResponse](h.Favorite.Get))
	}
}
