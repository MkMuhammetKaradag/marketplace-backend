package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	AddSeller(ctx context.Context, sellerID uuid.UUID, status string) error
	Close() error
}
