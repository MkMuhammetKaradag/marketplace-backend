package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"

	"github.com/google/uuid"
)

const getSellerByUserIDQuery = `
    SELECT 
        id, user_id, store_name, store_slug, store_description, 
        rating, total_sales, legal_business_name, tax_number, tax_office,
        phone_number, email, address_line_1, city, country,
        bank_account_iban, bank_account_holder_name, status, 
        rejection_reason, created_at
    FROM sellers 
    WHERE user_id = $1`

func (r *Repository) GetSellerByUserID(ctx context.Context, userID uuid.UUID) (*domain.Seller, error) {
	var s domain.Seller
	err := r.db.QueryRowContext(ctx, getSellerByUserIDQuery, userID).Scan(
		&s.ID,                    // 1. id
		&s.UserID,                // 2. user_id
		&s.StoreName,             // 3. store_name
		&s.StoreSlug,             // 4. store_slug
		&s.StoreDescription,      // 5. store_description (*string NULL destekler)
		&s.Rating,                // 6. rating
		&s.TotalSales,            // 7. total_sales
		&s.LegalBusinessName,     // 8. legal_business_name
		&s.TaxNumber,             // 9. tax_number
		&s.TaxOffice,             // 10. tax_office
		&s.PhoneNumber,           // 11. phone_number
		&s.Email,                 // 12. email
		&s.AddressLine,           // 13. address_line_1
		&s.City,                  // 14. city
		&s.Country,               // 15. country
		&s.BankAccountIban,       // 16. bank_account_iban
		&s.BankAccountHolderName, // 17. bank_account_holder_name
		&s.Status,                // 18. status 
		&s.RejectionReason,       // 19. rejection_reason (*string NULL destekler)
		&s.CreatedAt,             // 20. created_at
	)

	if err != nil {
		return nil, err
	}

	const historyQuery = `
        SELECT status, reason, created_at 
        FROM seller_status_history 
        WHERE seller_id = $1 
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, historyQuery, s.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h domain.SellerStatusHistory
		if err := rows.Scan(&h.Status, &h.Reason, &h.CreatedAt); err != nil {
			return nil, err
		}
		s.History = append(s.History, h)
	}
	return &s, nil
}
