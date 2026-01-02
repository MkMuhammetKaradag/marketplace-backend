package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/product-service/domain"
)

const (
	CREATE_PRODUCT = `
        INSERT INTO products (
            seller_id, 
            category_id,
            name, 
            description, 
            price, 
            stock_count,
            attributes
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id, created_at, updated_at`
)

func (r *Repository) CreateProduct(ctx context.Context, p *domain.Product) error {
	// 1. Map yapısını JSONB için byte slice'a çeviriyoruz
	attrBytes, err := json.Marshal(p.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}


	err = r.db.QueryRowContext(
		ctx,
		CREATE_PRODUCT,
		p.SellerID,    // $1
		p.CategoryID,  // $2 
		p.Name,        // $3
		p.Description, // $4
		p.Price,       // $5
		p.StockCount,  // $6
		attrBytes,     // $7
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}
