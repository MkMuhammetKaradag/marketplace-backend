package postgres

import (
	"context"

	"github.com/google/uuid"
)

const SOFT_DELETE_PRODUCT = `
    UPDATE products 
    SET status = 'deleted', updated_at = NOW() 
    WHERE id = $1
`

func (r *Repository) SoftDeleteProduct(ctx context.Context, productID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, SOFT_DELETE_PRODUCT, productID)
	return err
}
