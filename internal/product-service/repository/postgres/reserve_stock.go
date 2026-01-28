package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) ReserveStocks(ctx context.Context, orderID uuid.UUID, items []domain.OrderItemReserve) ([]domain.ProductInfo, error) {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("could not start transaction: %w", err)
    }
    defer tx.Rollback()

    reservedProducts := make([]domain.ProductInfo, 0, len(items))

    for _, item := range items {
        var pInfo domain.ProductInfo
        var availableStock int

        // Senin orijinal sorguna benzer ama rezervasyon için optimize edilmiş versiyon
        const checkAndGetQuery = `
            SELECT 
                p.name, 
                p.price, 
                p.seller_id,
                COALESCE(pi.image_url, '') as image_url, 
                (p.stock_count - COALESCE(res.total_reserved, 0)) as available_stock
            FROM products p
            -- İlk görseli çekmek için LATERAL JOIN
            LEFT JOIN LATERAL (
                SELECT image_url FROM product_images 
                WHERE product_id = p.id AND deleted_at IS NULL 
                ORDER BY created_at ASC LIMIT 1
            ) pi ON true
            -- Rezervasyon toplamlarını hesapla
            LEFT JOIN (
                SELECT product_id, SUM(quantity) as total_reserved 
                FROM product_reservations 
                WHERE expires_at > NOW() 
                GROUP BY product_id
            ) res ON p.id = res.product_id
            WHERE p.id = $1
            FOR UPDATE OF p
        `
        
        err := tx.QueryRowContext(ctx, checkAndGetQuery, item.ProductID).Scan(
            &pInfo.Name, 
            &pInfo.Price, 
            &pInfo.SellerID, // Order service için seller_id önemli
            &pInfo.ImageURL, 
            &availableStock,
        )
        if err != nil {
            return nil, fmt.Errorf("product %s not found or query error: %w", item.ProductID, err)
        }

        // Stok kontrolü
        if availableStock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for %s: available %d, requested %d",
                pInfo.Name, availableStock, item.Quantity)
        }

        // Rezervasyon kaydını ekle
        const reserveQuery = `
            INSERT INTO product_reservations (product_id, order_id, quantity, expires_at)
            VALUES ($1, $2, $3, NOW() + INTERVAL '30 minutes')
        `
        _, err = tx.ExecContext(ctx, reserveQuery, item.ProductID, orderID, item.Quantity)
        if err != nil {
            return nil, fmt.Errorf("failed to insert reservation for %s: %w", pInfo.Name, err)
        }

        pInfo.ID = item.ProductID
        reservedProducts = append(reservedProducts, pInfo)
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return reservedProducts, nil
}
