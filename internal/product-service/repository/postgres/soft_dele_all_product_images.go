package postgres

import (
	"context"

	"github.com/google/uuid"
)

const SOFT_DELETE_OLD_IMAGES = `
    UPDATE product_images 
    SET deleted_at = NOW(), is_main = false 
    WHERE product_id = $1 AND deleted_at IS NULL
`

func (r *Repository) SoftDeleteAllProductImages(ctx context.Context, productID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, SOFT_DELETE_OLD_IMAGES, productID)
	return err
}
