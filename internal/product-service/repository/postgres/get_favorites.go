package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/product-service/domain"
	"time"

	"github.com/google/uuid"
)

const (
	GET_USER_FAVORITES = `
		SELECT 
			p.id, p.name, p.price, p.stock_count, COALESCE(c.name, ''),
			COALESCE(json_agg(pi.*) FILTER (WHERE pi.id IS NOT NULL), '[]') as images,
			f.created_at as favorited_at
		FROM products p
		JOIN favorites f ON p.id = f.product_id
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN product_images pi ON p.id = pi.product_id
		WHERE f.user_id = $1
		GROUP BY p.id, f.created_at,c.name
		ORDER BY f.created_at DESC;
	`
)

func (r *Repository) GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*domain.FavoriteItem, error) {
	rows, err := r.db.QueryContext(ctx, GET_USER_FAVORITES, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favorites: %w", err)
	}
	defer rows.Close()

	var products []*domain.FavoriteItem
	for rows.Next() {
		p := &domain.FavoriteItem{}
		var imagesJSON []byte
		var favoritedAt time.Time

		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.StockCount, &p.CategoryName,
			&imagesJSON, &favoritedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(imagesJSON) > 0 {
			_ = json.Unmarshal(imagesJSON, &p.Images)
		}

		products = append(products, p)
	}

	return products, nil
}
