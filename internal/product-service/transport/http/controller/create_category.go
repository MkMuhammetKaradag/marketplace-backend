package controller

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"
	"marketplace/internal/product-service/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	ParentID    uuid.UUID `json:"parent_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type CreateCategoryResponse struct {
	Message string `json:"message"`
}
type CreateCategoryController struct {
	usecase usecase.CreateCategoryUseCase
}

func NewCreateCategoryController(usecase usecase.CreateCategoryUseCase) *CreateCategoryController {
	return &CreateCategoryController{
		usecase: usecase,
	}
}

func (c *CreateCategoryController) Handle(fiberCtx *fiber.Ctx, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {

	slug := util.Slugify(req.Name)
	p := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		Slug:        slug,
		ParentID:    req.ParentID,
	}
	err := c.usecase.Execute(fiberCtx.UserContext(), p)
	if err != nil {
		return nil, err
	}

	return &CreateCategoryResponse{Message: "Category created successfully"}, nil

}
