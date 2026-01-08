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

func (r *Repository) SearchProducts(ctx context.Context, queryVector []float32, req domain.SearchProductsParams) ([]*domain.Product, error) {

	var args []interface{}
	argID := 1

	// Vektör hazırlığı
	var vectorStr string
	if len(queryVector) > 0 {
		var strElements []string
		for _, v := range queryVector {
			strElements = append(strElements, fmt.Sprintf("%f", v))
		}
		vectorStr = "[" + strings.Join(strElements, ",") + "]"
	}

	// Temel sorgu (score hesaplama kısmı)
	selectClause := `SELECT p.id, p.seller_id, p.category_id, p.name, p.description, p.price, p.stock_count`
	if vectorStr != "" {
		selectClause += fmt.Sprintf(", (1 - (p.embedding <=> $%d)) as score", argID)
		args = append(args, vectorStr)
		argID++
	} else {
		selectClause += ", 0 as score"
	}

	baseQuery := selectClause + " FROM products p WHERE (p.status = 'active' OR p.status = 'inactive')  AND p.stock_count > 0"

	// --- Dinamik Filtreler ---

	// 1. Text Search (Eğer query varsa)
	if req.Query != "" {
		baseQuery += fmt.Sprintf(` AND ((1 - (p.embedding <=> $1)) > 0.4 OR p.name ILIKE $%d OR p.description ILIKE $%d)`, argID, argID)
		args = append(args, "%"+req.Query+"%")
		argID++
	}

	// 2. Fiyat Filtreleri
	if req.MinPrice != nil {
		baseQuery += fmt.Sprintf(" AND p.price >= $%d", argID)
		args = append(args, *req.MinPrice)
		argID++
	}
	if req.MaxPrice != nil {
		baseQuery += fmt.Sprintf(" AND p.price <= $%d", argID)
		args = append(args, *req.MaxPrice)
		argID++
	}

	// 3. Kategori Filtresi
	if req.CategoryID != nil && *req.CategoryID != "" {
		baseQuery += fmt.Sprintf(" AND p.category_id = $%d", argID)
		args = append(args, *req.CategoryID)
		argID++
	}

	// Sıralama ve Limit
	baseQuery += " ORDER BY score DESC LIMIT " + fmt.Sprintf("$%d", argID)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	args = append(args, req.Limit)

	// Sorguyu çalıştır
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
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
