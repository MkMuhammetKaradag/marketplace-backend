package domain

import "github.com/google/uuid"

type BasketItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"image_url"`
}

type Basket struct {
	UserID uuid.UUID    `json:"user_id"`
	Items  []BasketItem `json:"items"`
}

type BasketItemResponse struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	SubTotal  float64 `json:"sub_total"`
}

type BasketResponse struct {
	UserID     string               `json:"user_id"`
	Items      []BasketItemResponse `json:"items"`
	TotalPrice float64              `json:"total_price"`
}
