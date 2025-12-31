package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	AddSeller(ctx context.Context, sellerID uuid.UUID, status string) error
	CreateProduct(ctx context.Context, p *Product) error
	SaveImagesAndUpdateStatus(ctx context.Context, productID uuid.UUID, images []ProductImage) error
	Close() error
}
