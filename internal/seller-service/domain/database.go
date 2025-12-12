package domain

import "context"

type SellerRepository interface {
	Close() error
	Create(ctx context.Context, seller *Seller) (string, error)
}
