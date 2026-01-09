package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	IsMain    bool      `json:"is_main"`
	SortOrder int       `json:"sort_order"`
}
type Product struct {
	ID           uuid.UUID              `json:"id"`
	SellerID     uuid.UUID              `json:"seller_id"`
	CategoryID   uuid.UUID              `json:"category_id"`
	CategoryName string                 `json:"category_name"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Price        float64                `json:"price"`
	StockCount   int                    `json:"stock_count"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	Attributes   map[string]interface{} `json:"attributes"`
	UpdatedAt    time.Time              `json:"updated_at"`
	IsFavorited  bool                   `json:"is_favorited"`

	Images    []ProductImage `json:"images"`
	Embedding []float32      `json:"-"`
}
type SearchProductsParams struct {
	Limit      int      `json:"limit"`
	Query      string   `json:"query"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	CategoryID *string  `json:"category_id"`
}
