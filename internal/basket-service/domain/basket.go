package domain

import "github.com/google/uuid"

type BasketItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"image_url"`
}

type Basket struct {
	UserID uuid.UUID    `json:"user_id"`
	Items  []BasketItem `json:"items"`
}
