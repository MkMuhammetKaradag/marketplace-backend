package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const UPDATE_STORE_LOGO = "UPDATE sellers SET store_logo_url = $1 WHERE id = $2 AND user_id = $3"

func (p *Repository) UpdateStoreLogo(ctx context.Context, userID uuid.UUID, sellerID uuid.UUID, storeLogo string) error {

	result, err := p.db.ExecContext(ctx, UPDATE_STORE_LOGO, storeLogo, sellerID, userID)
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

	fmt.Printf("Logo updated: %s for seller: %s\n", storeLogo, sellerID)
	return nil
}
