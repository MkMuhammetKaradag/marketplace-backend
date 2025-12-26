package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const UPDATE_STORE_BANNER = "UPDATE sellers SET store_banner_url = $1 WHERE id = $2 AND user_id = $3"

func (p *Repository) UpdateStoreBanner(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, storeBanner string) error {

	result, err := p.db.ExecContext(ctx, UPDATE_STORE_BANNER, storeBanner, sellerID, userID)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("seller not found or unauthorized: seller_id %s, user_id %s", sellerID, userID)
	}

	fmt.Printf("Banner updated: %s for seller: %s\n", storeBanner, sellerID)
	return nil
}
