package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
)

const createSellerQuery = `
    INSERT INTO sellers (
        user_id, store_name, store_slug, 
        legal_business_name, tax_number, tax_office, 
        phone_number, email, 
        address_line_1, city, country, 
        bank_account_iban, bank_account_holder_name,bank_account_bic
    ) 
    VALUES (
        $1, $2, $3, 
        $4, $5, $6, 
        $7, $8, 
        $9, $10, $11, 
        $12, $13,$14
    ) 
    RETURNING id;
`

func (r *Repository) Create(ctx context.Context, seller *domain.Seller) (string, error) {

	var sellerID string

	err := r.db.QueryRowContext(ctx, createSellerQuery,
		seller.UserID,
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
	).Scan(&sellerID)

	if err != nil {

		fmt.Printf("PostgreSQL  save seller error: %v\n", err)
		return "", fmt.Errorf("seller save error: %w", err)
	}

	return sellerID, nil
}
