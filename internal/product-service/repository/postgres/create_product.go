package postgres

import (
	"context"
	"encoding/json" // JSON dönüşümü için gerekli
	"fmt"
	"marketplace/internal/product-service/domain"
)

const (
	CREATE_PRODUCT = `
		INSERT INTO products (
			seller_id, 
			name, 
			description, 
			price, 
			stock_count,
			attributes
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
)

func (r *Repository) CreateProduct(ctx context.Context, p *domain.Product) error {

	attrBytes, err := json.Marshal(p.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		CREATE_PRODUCT,
		p.SellerID,
		p.Name,
		p.Description,
		p.Price,
		p.StockCount,
		attrBytes,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}
