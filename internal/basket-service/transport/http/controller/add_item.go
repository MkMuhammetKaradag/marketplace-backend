package controller

import (
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"image_url"`
}

type AddItemResponse struct {
	Message string `json:"message"`
}
type AddItemController struct {
	usecase usecase.AddItemUseCase
}

func NewAddItemController(usecase usecase.AddItemUseCase) *AddItemController {
	return &AddItemController{
		usecase: usecase,
	}
}

func (c *AddItemController) Handle(fiberCtx *fiber.Ctx, req *AddItemRequest) (*AddItemResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	p := &domain.BasketItem{
		ProductID: req.ProductID,
		Name:      req.Name,
		Price:     req.Price,
		Quantity:  req.Quantity,
		ImageURL:  req.ImageURL,
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId, p)
	if err != nil {
		return nil, err
	}

	return &AddItemResponse{Message: "Item added to basket successfully"}, nil

}
