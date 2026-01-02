package postgres

import (
	"context"

	"github.com/google/uuid"
)

const (
    ADD_SELLER = `
        INSERT INTO local_sellers (seller_id, user_id, status, updated_at) 
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (seller_id) 
        DO UPDATE SET 
            status = EXCLUDED.status, 
            updated_at = NOW()`
)

func (r *Repository) AddSeller(ctx context.Context, sellerID uuid.UUID, userID uuid.UUID, status string) error {
    _, err := r.db.ExecContext(ctx, ADD_SELLER, sellerID, userID, status)
    if err != nil {
        return err
    }
    return nil
}