package postgres

import (
	"context"
	"marketplace/internal/seller-service/domain"
)

const updateForReapplicationQuery = `
    UPDATE sellers 
    SET 
        store_name = $1, 
        store_slug = $2, 
        legal_business_name = $3, 
        tax_number = $4, 
        tax_office = $5, 
        phone_number = $6, 
        email = $7, 
        address_line_1 = $8, 
        city = $9, 
        country = $10, 
        bank_account_iban = $11, 
        bank_account_holder_name = $12, 
        bank_account_bic = $13,
        status = 'pending',       -- Durumu sıfırla
        rejection_reason = NULL,  -- Eski sebebi sil
        updated_at = NOW()
    WHERE id = $14 AND status = 'rejected';
`

func (r *Repository) UpdateForReapplication(ctx context.Context, seller *domain.Seller) error {
	_, err := r.db.ExecContext(ctx, updateForReapplicationQuery,
		seller.StoreName,
		seller.StoreSlug,
		seller.LegalBusinessName,
		seller.TaxNumber,
		seller.TaxOffice,
		seller.PhoneNumber,
		seller.Email,
		seller.AddressLine,
		seller.City,
		seller.Country,
		seller.BankAccountIban,
		seller.BankAccountHolderName,
		seller.BankAccountBic,
		seller.ID, // Where koşulu için
	)
	return err
}
