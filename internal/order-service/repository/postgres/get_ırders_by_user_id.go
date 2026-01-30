package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	const query = `
        SELECT 
            o.id, o.user_id, o.total_price, o.status, o.shipping_address, o.created_at,
            oi.id, oi.product_id, oi.seller_id, oi.quantity, oi.product_name, oi.product_image_url, oi.unit_price, oi.status
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user orders: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[uuid.UUID]*domain.Order)
	var sortedIDs []uuid.UUID

	for rows.Next() {
		//var oID uuid.UUID
		var oStatus, oiStatus int
		var item domain.OrderItem
		var order domain.Order
		var shippingAddr sql.NullString

		err := rows.Scan(
			&order.ID, &order.UserID, &order.TotalPrice, &oStatus, &shippingAddr, &order.CreatedAt,
			&item.ID, &item.ProductID, &item.SellerID, &item.Quantity, &item.ProductName, &item.ProductImageUrl, &item.UnitPrice, &oiStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		if shippingAddr.Valid {
			order.ShippingAddress = shippingAddr.String
		} else {
			order.ShippingAddress = ""
		}

		order.Status = domain.OrderStatus(oStatus)
		item.Status = domain.OrderItemStatus(oiStatus)

		if existingOrder, ok := ordersMap[order.ID]; ok {
			existingOrder.Items = append(existingOrder.Items, item)
		} else {
			order.Items = []domain.OrderItem{item}
			ordersMap[order.ID] = &order
			sortedIDs = append(sortedIDs, order.ID)
		}
	}

	var result []domain.Order
	for _, id := range sortedIDs {
		result = append(result, *ordersMap[id])
	}

	return result, nil
}
