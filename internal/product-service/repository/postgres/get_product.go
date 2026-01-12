package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/product-service/domain"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func parseVector(s string) ([]float32, error) {
	// Başındaki ve sonundaki [] karakterlerini atıyoruz
	s = strings.Trim(s, "[]")
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	vec := make([]float32, len(parts))
	for i, p := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return nil, err
		}
		vec[i] = float32(val)
	}
	return vec, nil
}
func (r *Repository) GetProduct(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*domain.Product, error) {
	p := &domain.Product{}
	var imagesJSON []byte
	var attributesJSON []byte
	var embeddingStr string

	const query = `
        SELECT 
            p.id, p.seller_id, p.category_id, COALESCE(c.name, ''),
            p.name, p.description, p.price, p.stock_count, p.status, 
            p.attributes, p.embedding::text, 
            COALESCE(json_agg(pi.*) FILTER (WHERE pi.id IS NOT NULL), '[]'),
            p.created_at, p.updated_at,
            -- Favori kontrolü: Eğer userID varsa favorites tablosuna bak
            CASE WHEN $2::uuid IS NOT NULL THEN
                EXISTS (SELECT 1 FROM favorites WHERE user_id = $2 AND product_id = p.id)
            ELSE false END as is_favorited
            
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.deleted_at IS NULL
        WHERE p.id = $1
        GROUP BY p.id, c.name;
    `

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.CategoryName,
		&p.Name, &p.Description, &p.Price, &p.StockCount, &p.Status,
		&attributesJSON, &embeddingStr, &imagesJSON, &p.CreatedAt, &p.UpdatedAt,
		&p.IsFavorited,
	)
	if err != nil {
		return nil, err
	}

	if embeddingStr != "" {
		p.Embedding, err = parseVector(embeddingStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vector: %w", err)
		}
	}

	p.Attributes = make(map[string]interface{})
	if len(attributesJSON) > 0 {
		_ = json.Unmarshal(attributesJSON, &p.Attributes)
	}

	if len(imagesJSON) > 0 {
		_ = json.Unmarshal(imagesJSON, &p.Images)
	}

	return p, nil
}
