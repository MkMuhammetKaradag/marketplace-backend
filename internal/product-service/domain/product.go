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
	ID             uuid.UUID              `json:"id"`
	SellerID       uuid.UUID              `json:"seller_id"`
	CategoryID     uuid.UUID              `json:"category_id"`
	CategoryName   string                 `json:"category_name"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Price          float64                `json:"price"`
	StockCount     int                    `json:"stock_count"`
	Status         string                 `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	Attributes     map[string]interface{} `json:"attributes"`
	UpdatedAt      time.Time              `json:"updated_at"`
	IsFavorited    bool                   `json:"is_favorited"`
	AvailableStock int                    `json:"available_stock"`

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
type FavoriteItem struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	StockCount   int            `json:"stock_count"`
	Images       []ProductImage `json:"images"`
	FavoritedAt  time.Time      `json:"favorited_at"`
	CategoryName string         `json:"category_name"`
}
type UpdateProduct struct {
	ProductID   uuid.UUID
	Name        *string
	Description *string
	Price       *float64
	StockCount  *int
	CategoryID  *uuid.UUID
	Attributes  map[string]interface{}
}
