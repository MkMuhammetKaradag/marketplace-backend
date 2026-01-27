package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) ReserveStocks(ctx context.Context, orderID uuid.UUID, items []domain.OrderItemReserve) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer tx.Rollback()

	for _, item := range items {

		var availableStock int
		const checkQuery = `
			SELECT (p.stock_count - COALESCE(SUM(r.quantity), 0)) as available
			FROM products p
			LEFT JOIN product_reservations r ON p.id = r.product_id AND r.expires_at > NOW()
			WHERE p.id = $1
			GROUP BY p.id, p.stock_count
		`
		err := tx.QueryRowContext(ctx, checkQuery, item.ProductID).Scan(&availableStock)
		if err != nil {
			return fmt.Errorf("failed to check stock for product %s: %w", item.ProductID, err)
		}

		if availableStock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s: available %d, requested %d",
				item.ProductID, availableStock, item.Quantity)
		}

		const reserveQuery = `
			INSERT INTO product_reservations (product_id, order_id, quantity, expires_at)
			VALUES ($1, $2, $3, NOW() + INTERVAL '3 minutes')
		`
		_, err = tx.ExecContext(ctx, reserveQuery, item.ProductID, orderID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert reservation for product %s: %w", item.ProductID, err)
		}
	}

	return tx.Commit()
}
