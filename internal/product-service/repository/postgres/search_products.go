package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"
	"strings"
)

const (
	SEARCH_PRODUCTS = `
    SELECT 
        p.id, p.seller_id, p.category_id, p.name, p.description, p.price, p.stock_count,
        (1 - (p.embedding <=> $1)) as score
    FROM products p
    WHERE 
        (p.status = 'active' OR p.status = 'inactive') 
        AND p.stock_count >= 0
        AND (

            (1 - (p.embedding <=> $1)) > 0.4
            OR 
            p.name ILIKE '%' || $2 || '%' 
            OR 
            p.description ILIKE '%' || $2 || '%'
        )
    ORDER BY score DESC
    LIMIT $3;
`
)

func (r *Repository) SearchProducts(ctx context.Context, queryVector []float32, queryText string, limit int) ([]*domain.Product, error) {
	var strElements []string
	for _, v := range queryVector {
		strElements = append(strElements, fmt.Sprintf("%f", v))
	}
	vectorStr := "[" + strings.Join(strElements, ",") + "]"

	rows, err := r.db.QueryContext(ctx, SEARCH_PRODUCTS, vectorStr, queryText, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.Product, 0)
	for rows.Next() {
		p := &domain.Product{}
		var score float64
		err := rows.Scan(&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.StockCount, &score)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
