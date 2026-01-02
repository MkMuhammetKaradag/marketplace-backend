package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	GET_SELLERID_BY_USERID = `
		SELECT seller_id FROM local_sellers WHERE user_id = $1 AND status = 'approved' LIMIT 1
	`
)

func (r *Repository) GetSellerIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var sellerID uuid.UUID
	err := r.db.QueryRowContext(ctx, GET_SELLERID_BY_USERID, userID).Scan(&sellerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("seller not found: %w", err)
		}
		return uuid.Nil, err
	}
	return sellerID, nil
}
