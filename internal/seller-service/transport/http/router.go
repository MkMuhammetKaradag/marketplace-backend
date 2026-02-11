// internal/seller-service/transport/http/router.go
package http

import (
	"marketplace/internal/seller-service/handler"
	"marketplace/internal/seller-service/transport/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	h := r.handlers

	
	app.Get("/hello", h.Hello)

	
	store := app.Group("/store")
	{
		store.Post("/onboard", handler.HandleWithFiber[controller.CreateSellerRequest, controller.CreateSellerResponse](h.Store.Create))
		store.Get("/me", handler.HandleWithFiber[controller.GetSellerByUserIDRequest, controller.GetSellerByUserIDResponse](h.Store.GetByUserID))
		store.Post("/upload-logo/:seller_id", handler.HandleWithFiber[controller.UploadStoreLogoRequest, controller.UploadStoreLogoResponse](h.Store.UploadLogo))
		store.Post("/upload-banner/:seller_id", handler.HandleWithFiber[controller.UploadStoreBannerRequest, controller.UploadStoreBannerResponse](h.Store.UploadBanner))
	}

	
	admin := app.Group("/admin/sellers")
	{
		admin.Post("/approve/:seller_id", handler.HandleWithFiber[controller.ApproveSellerRequest, controller.ApproveSellerResponse](h.Admin.Approve))
		admin.Post("/reject/:seller_id", handler.HandleWithFiber[controller.RejectSellerRequest, controller.RejectSellerResponse](h.Admin.Reject))
	}
}
