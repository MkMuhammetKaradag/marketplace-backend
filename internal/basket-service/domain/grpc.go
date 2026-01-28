package domain

import (
	"context"
	pb "marketplace/pkg/proto/product"
)

type ProductClient interface {
	GetProductForBasket(ctx context.Context, id string) (*pb.ProductResponse, error)
	Close() error
}
