package postgres

import (
	"context"
	"marketplace/internal/product-service/domain"
)

const (
	CREATE_PRODUCT = `
        INSERT INTO products (
            seller_id, 
            name, 
            description, 
            price, 
            stock_count
     
        ) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`
)

func (r *Repository) CreateProduct(
	ctx context.Context,
	p *domain.Product,
) error {

	err := r.db.QueryRowContext(
		ctx,
		CREATE_PRODUCT,
		p.SellerID,
		p.Name,
		p.Description,
		p.Price,
		p.StockCount,
		// p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt) // Oluşan ID'yi newID değişkenine aktar

	if err != nil {
		return err
	}

	return nil
}
