package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) ReleaseStock(ctx context.Context, orderID uuid.UUID) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// const updateStockQuery = `
	//     UPDATE products p
	//     SET stock_count = p.stock_count + pr.quantity
	//     FROM product_reservations pr
	//     WHERE p.id = pr.product_id AND pr.order_id = $1
	// `
	// _, err = tx.ExecContext(ctx, updateStockQuery, orderID)
	// if err != nil {
	// 	return fmt.Errorf("failed to update final stock from reservations: %w", err)
	// }

	const deleteReservationsQuery = `DELETE FROM product_reservations WHERE order_id = $1`
	_, err = tx.ExecContext(ctx, deleteReservationsQuery, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete reservations: %w", err)
	}

	return tx.Commit()
}
