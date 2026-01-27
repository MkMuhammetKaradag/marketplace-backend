package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	const updateOrderQuery = `
        UPDATE orders 
        SET status = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $2`

	_, err = tx.ExecContext(ctx, updateOrderQuery, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}


	const updateItemsQuery = `
        UPDATE order_items 
        SET status = $1 
        WHERE order_id = $2`

	_, err = tx.ExecContext(ctx, updateItemsQuery, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order items status: %w", err)
	}

	return tx.Commit()
}
