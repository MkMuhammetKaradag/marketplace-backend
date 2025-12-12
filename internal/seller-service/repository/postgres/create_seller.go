package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
)

// createSellerQuery, INSERT işlemi için SQL sorgusudur
// NOT: Yalnızca kullanıcının request ile gönderdiği ve NOT NULL olan alanları ekliyoruz.
// created_at, updated_at, id, is_approved, rating gibi alanlar DB tarafından yönetilir.
const createSellerQuery = `
    INSERT INTO sellers (
        user_id, store_name, store_slug, 
        legal_business_name, tax_number, tax_office, 
        phone_number, email, 
        address_line_1, city, country, 
        bank_account_iban, bank_account_holder_name, bank_account_bic
    ) 
    VALUES (
        $1, $2, $3, 
        $4, $5, $6, 
        $7, $8, 
        $9, $10, $11, 
        $12, $13, $14
    ) 
    RETURNING id;
`

// Create, yeni bir satıcı kaydını veritabanına ekler
func (r *Repository) Create(ctx context.Context, seller *domain.Seller) (string, error) {
	// 1. Store Slug kontrolü (varsayımsal olarak usecase'de oluşturuldu)


	var sellerID string

	// 2. Sorguyu çalıştır
	err := r.db.QueryRowContext(ctx, createSellerQuery,
		seller.UserID,
		seller.StoreName,
		seller.StoreSlug, // Slug'ın domain.Seller yapısında olduğunu varsayıyoruz.
		seller.LegalBusinessName,
		seller.TaxNumber,
		seller.TaxOffice,
		seller.PhoneNumber,
		seller.Email,
		// domain.Seller'da AddressLine olarak gelen alan, DB'de address_line_1'e gider
		seller.AddressLine,
		seller.City,
		seller.Country,
		seller.BankAccountIban,
		seller.BankAccountHolderName,
		seller.BankAccountBic,
	).Scan(&sellerID) // 3. Oluşturulan ID'yi al

	if err != nil {
		// PostgreSQL'den gelen benzersizlik (UNIQUE constraint) hatalarını burada yakalamak önemlidir.
		// Örn: store_name, tax_number, user_id zaten kullanılmış olabilir.
		// Gerçek projede 'github.com/lib/pq' veya benzeri kütüphane ile hata kodları kontrol edilir.
		fmt.Printf("PostgreSQL kayıt hatası: %v\n", err)
		return "", fmt.Errorf("satıcı kaydı yapılamadı: %w", err)
	}

	// 4. Başarıyla oluşturulan ID'yi döndür
	return sellerID, nil
}
