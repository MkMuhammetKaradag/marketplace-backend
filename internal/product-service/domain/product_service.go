package domain

import (
	"context"
)

type ProductService interface {
	Greeting(ctx context.Context) string
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Greeting(_ context.Context) string {
	return "Hello from Product Service!"
}
