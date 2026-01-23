package domain

import (
	"context"

	pp "marketplace/pkg/proto/Product"
	pb "marketplace/pkg/proto/basket"
	pPayment "marketplace/pkg/proto/payment"
)

type BasketClient interface {
	GetBasket(ctx context.Context, id string) (*pb.BasketResponse, error)
	Close() error
}

type ProductClient interface {
	GetProductsByIds(ctx context.Context, ids []string) (*pp.GetProductsByIdsResponse, error)
	Close() error
}
type PaymentClient interface {
	CreatePaymentSession(ctx context.Context, orderID, userID, email string, amount float64) (*pPayment.CreatePaymentResponse, error)
	Close() error
}
