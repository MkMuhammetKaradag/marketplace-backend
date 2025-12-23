// internal/seller-service/transport/http/usecase/seller_onboard.go
package usecase

import (
	"context"
	"database/sql"
	"errors"
	"marketplace/internal/seller-service/domain"

	"github.com/google/uuid"
)

type GetSellerByUserIDUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (*domain.Seller, error)
}
type getSellerByUserIDUseCase struct {
	sellerRepository domain.SellerRepository
}

func NewGetSellerByUserIDUseCase(repository domain.SellerRepository) GetSellerByUserIDUseCase {
	return &getSellerByUserIDUseCase{
		sellerRepository: repository,
	}
}

func (u *getSellerByUserIDUseCase) Execute(ctx context.Context, userID uuid.UUID) (*domain.Seller, error) {

	seller, err := u.sellerRepository.GetSellerByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, nil
		}
		return nil, err
	}

	return seller, nil
}
