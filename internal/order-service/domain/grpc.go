package domain

import (
	"context"

	pp "marketplace/pkg/proto/Product"
	pb "marketplace/pkg/proto/basket"
)

type BasketClient interface {
	GetBasket(ctx context.Context, id string) (*pb.BasketResponse, error)
	Close() error
}

type ProductClient interface {
	GetProductsByIds(ctx context.Context, ids []string) (*pp.GetProductsByIdsResponse, error)
	Close() error
}
