package controller

import (
	"marketplace/internal/product-service/transport/http/usecase"
	"marketplace/internal/user-service/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DeleteProductRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type DeleteProductResponse struct {
	Message string `json:"message"`
}
type DeleteProductController struct {
	usecase usecase.DeleteProductUseCase
}

func NewDeleteProductController(usecase usecase.DeleteProductUseCase) *DeleteProductController {
	return &DeleteProductController{
		usecase: usecase,
	}
}

// Handle godoc
// @Summary delete product
// @Description delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} DeleteProductResponse
// @Router /products/delete/{product_id} [delete]
func (c *DeleteProductController) Handle(fiberCtx *fiber.Ctx, req *DeleteProductRequest) (*DeleteProductResponse, error) {
	parsedUserID, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}
	userPermissions, err := strconv.ParseInt(fiberCtx.Get("X-User-Permissions"), 10, 64)
	if err != nil {
		return nil, err
	}

	isAdmin := (userPermissions & domain.PermissionAdministrator) != 0
	err = c.usecase.Execute(fiberCtx.UserContext(), parsedUserID, req.ProductID, isAdmin)
	if err != nil {
		return nil, err
	}

	return &DeleteProductResponse{Message: "Product deleted successfully"}, nil

}
