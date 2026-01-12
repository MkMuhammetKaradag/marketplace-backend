package postgres

import (
	"context"
	"encoding/json"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

const (
	GET_PRODUCT_INTERNAL = `
        SELECT 
            id, seller_id, category_id, name, description, 
            price, stock_count, status, attributes, embedding::text
        FROM products 
        WHERE id = $1 
    `
)

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	p := &domain.Product{}
	var attributesJSON []byte
	var embeddingStr string

	err := r.db.QueryRowContext(ctx, GET_PRODUCT_INTERNAL, id).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description,
		&p.Price, &p.StockCount, &p.Status, &attributesJSON, &embeddingStr,
	)
	if err != nil {
		return nil, err
	}

	if len(attributesJSON) > 0 {
		_ = json.Unmarshal(attributesJSON, &p.Attributes)
	}

	return p, nil
}
