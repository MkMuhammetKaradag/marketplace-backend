package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

const GET_PERSONALIZED_FEED = `
    SELECT p.id, p.seller_id, p.category_id, p.name, p.description, p.price, p.stock_count
    FROM products p
    -- Kullanıcının zevk profilini sol tarafa alıyoruz
    LEFT JOIN user_preferences up ON up.user_id = $1
    WHERE (p.status = 'active' OR p.status = 'inactive') AND p.stock_count > 0
    ORDER BY 
        -- Eğer kullanıcının tercihi varsa benzerliğe göre, yoksa en yeniye göre sırala
        CASE 
            WHEN up.interest_vector IS NOT NULL THEN p.embedding <=> up.interest_vector 
            ELSE NULL 
        END ASC, 
        p.created_at DESC
    LIMIT $2;
`

func (r *Repository) GetRecommendedProducts(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, GET_PERSONALIZED_FEED, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		// Kendi struct yapına göre Scan işlemini düzenle
		err := rows.Scan(&p.ID, &p.SellerID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.StockCount)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
