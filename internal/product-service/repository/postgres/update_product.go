package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/product-service/domain"
)

const UPDATE_PRODUCT = `
    UPDATE products 
    SET 
        name = $1, 
        description = $2, 
        price = $3, 
        stock_count = $4, 
        status = $5, 
        category_id = $6, 
        attributes = $7, 
        updated_at = NOW()
    WHERE id = $8 AND seller_id = $9
`

func (r *Repository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	// Attributes map'ini JSON formatına çeviriyoruz
	attrJSON, err := json.Marshal(p.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	_, err = r.db.ExecContext(ctx, UPDATE_PRODUCT,
		p.Name,
		p.Description,
		p.Price,
		p.StockCount,
		p.Status,
		p.CategoryID,
		attrJSON,
		p.ID,
		p.SellerID,
	)
	return err
}
