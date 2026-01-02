package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	AddSeller(ctx context.Context, sellerID uuid.UUID, userID uuid.UUID, status string) error
	GetSellerIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateProduct(ctx context.Context, p *Product) error
	SaveImagesAndUpdateStatus(ctx context.Context, productID uuid.UUID, images []ProductImage) error
	CreateCategory(ctx context.Context, c *Category) error
	Close() error
}
