package domain

import "context"

type SellerRepository interface {
	Close() error
	Create(ctx context.Context, seller *Seller) (string, error)
	ApproveSeller(ctx context.Context, sellerId, approvedBy string) (string, error)
	RejectSeller(ctx context.Context, sellerId string, rejectedBy string, rejectionReason string) (string, error)
}
