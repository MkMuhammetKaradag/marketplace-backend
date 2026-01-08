package postgres

import (
	"context"

	"github.com/google/uuid"
)

const (
	ADD_FAVORITE = `INSERT INTO favorites (user_id, product_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	REMOVE_FAVORITE = `DELETE FROM favorites WHERE user_id = $1 AND product_id = $2`

	IS_FAVORITED = `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND product_id = $2)`
)

func (r *Repository) ToggleFavorite(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, IS_FAVORITED, userID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {

		_, err = r.db.ExecContext(ctx, REMOVE_FAVORITE, userID, productID)
		return false, err
	}

	_, err = r.db.ExecContext(ctx, ADD_FAVORITE, userID, productID)
	return true, err
}
