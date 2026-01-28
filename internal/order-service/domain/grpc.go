package domain

import (
	"context"

	pb "marketplace/pkg/proto/basket"
	pPayment "marketplace/pkg/proto/payment"
	pp "marketplace/pkg/proto/product"
	cp "marketplace/pkg/proto/common"
)

type BasketClient interface {
	GetBasket(ctx context.Context, id string) (*pb.BasketResponse, error)
	Close() error
}

type ProductClient interface {
	GetProductsByIds(ctx context.Context, ids []string) (*pp.GetProductsByIdsResponse, error)
	ReserveStock(ctx context.Context, orderID string, items []*cp.OrderItemData) (*pp.ReserveStockResponse, error)
	Close() error
}
type PaymentClient interface {
	CreatePaymentSession(ctx context.Context, orderID, userID, email string, amount float64) (*pPayment.CreatePaymentResponse, error)
	Close() error
}
