// internal/user-service/transport/http/router.go
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
	//api := app.Group("/api/v1")
	createSellerHandler := r.handlers.CreateSeller()
	approveSellerHandler := r.handlers.ApproveSeller()
	rejectSellerHandler := r.handlers.RejectSeller()
	getSellerByUserIDHandler := r.handlers.GetSellerByUserID()
	app.Get("/hello", r.handlers.Hello)
	app.Post("/onboard", handler.HandleWithFiber[controller.CreateSellerRequest, controller.CreateSellerResponse](createSellerHandler))
	app.Post("/approve/:seller_id", handler.HandleWithFiber[controller.ApproveSellerRequest, controller.ApproveSellerResponse](approveSellerHandler))
	app.Post("/reject/:seller_id", handler.HandleWithFiber[controller.RejectSellerRequest, controller.RejectSellerResponse](rejectSellerHandler))
	app.Get("/me", handler.HandleWithFiber[controller.GetSellerByUserIDRequest, controller.GetSellerByUserIDResponse](getSellerByUserIDHandler))

}
