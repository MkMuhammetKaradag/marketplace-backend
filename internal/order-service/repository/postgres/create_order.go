package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"
)

func (r *Repository) CreateOrder(ctx context.Context, order *domain.Order) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	
	defer tx.Rollback()


	const orderQuery = `
        INSERT INTO orders (id, user_id, total_price, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err = tx.ExecContext(ctx, orderQuery,
		order.ID, order.UserID, order.TotalPrice, order.Status)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	
	const itemQuery = `
        INSERT INTO order_items (
            id, order_id, product_id, seller_id, quantity, 
            product_name, product_image_url, unit_price, status
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, item.OrderID, item.ProductID, item.SellerID, item.Quantity,
			item.ProductName, item.ProductImageUrl, item.UnitPrice, item.Status)

		if err != nil {
			return fmt.Errorf("failed to insert order item %s: %w", item.ProductID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
