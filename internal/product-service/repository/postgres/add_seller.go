package postgres

import (
	"context"

	"github.com/google/uuid"
)

const (
	ADD_SELLER = "INSERT INTO local_sellers (seller_id,status) VALUES ($1,$2)"
)

func (r *Repository) AddSeller(ctx context.Context, sellerID uuid.UUID, status string) error {

	_, err := r.db.ExecContext(ctx, ADD_SELLER, sellerID, status)
	if err != nil {
		return err
	}
	return nil
}
