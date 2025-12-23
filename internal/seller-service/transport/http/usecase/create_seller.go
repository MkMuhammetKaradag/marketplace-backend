// internal/seller-service/transport/http/usecase/seller_onboard.go
package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/util"

	"github.com/google/uuid"
)

type CreateSellerUseCase interface {
	Execute(ctx context.Context, seller *domain.Seller) (string, error)
}
type createSellerUseCase struct {
	sellerRepository domain.SellerRepository
}

func NewCreateSellerUseCase(repository domain.SellerRepository) CreateSellerUseCase {
	return &createSellerUseCase{
		sellerRepository: repository,
	}
}

func (u *createSellerUseCase) Execute(ctx context.Context, seller *domain.Seller) (string, error) {
	parsedUserID, err := uuid.Parse(seller.UserID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	existingSeller, err := u.sellerRepository.GetSellerByUserID(ctx, parsedUserID)

	// Eğer veritabanı hatası varsa (ama 'kayıt yok' hatası değilse) hata dön
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("database check error: %w", err)
	}
	// 2. Eğer kayıt varsa durumuna göre davran
	if existingSeller != nil {
		switch existingSeller.Status {
		case "approved":
			return "", fmt.Errorf("zaten onaylanmış bir mağazanız bulunmaktadır")
		case "pending":
			return "", fmt.Errorf("devam eden bir başvurunuz var, lütfen sonuçlanmasını bekleyin")
		case "rejected":
			// Reddildiyse bilgileri güncelle ve durumu tekrar 'pending'e çek
			seller.ID = existingSeller.ID // Mevcut ID üzerinden güncelleme yapması için
			seller.StoreSlug = util.Slugify(seller.StoreName)

			err := u.sellerRepository.UpdateForReapplication(ctx, seller)
			if err != nil {
				return "", err
			}
			return seller.ID, nil
		}
	}
	seller.StoreSlug = util.Slugify(seller.StoreName)
	if seller.StoreSlug == "" {
		return "", fmt.Errorf("The store name cannot create a valid slug.")
	}

	id, err := u.sellerRepository.Create(ctx, seller)
	if err != nil {
		return "", err
	}

	return id, nil
}
