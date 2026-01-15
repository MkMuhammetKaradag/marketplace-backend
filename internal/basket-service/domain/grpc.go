package domain

import (
	"context"
	pb "marketplace/pkg/proto/Product"
)

type ProductClient interface {
	GetProductForBasket(ctx context.Context, id string) (*pb.ProductResponse, error)
	Close() error
}
