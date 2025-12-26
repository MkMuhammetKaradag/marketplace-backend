package domain

import (
	"context"

	"github.com/google/uuid"
)

type SellerRepository interface {
	Close() error
	Create(ctx context.Context, seller *Seller) (string, error)
	ApproveSeller(ctx context.Context, sellerId, approvedBy string) (string, error)
	RejectSeller(ctx context.Context, sellerId string, rejectedBy string, rejectionReason string) (string, error)
	GetSellerByUserID(ctx context.Context, userID uuid.UUID) (*Seller, error)
	UpdateForReapplication(ctx context.Context, seller *Seller) error
	UpdateStoreLogo(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, storeLogo string) error
	UpdateStoreBanner(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, storeBanner string) error
}
