package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus int
type OrderItemStatus int

const (
	OrderPending   OrderStatus = 1
	OrderPaid      OrderStatus = 2
	OrderShipped   OrderStatus = 3
	OrderCancelled OrderStatus = 4
	OrderCompleted OrderStatus = 5

	OrderItemPending   OrderItemStatus = 1
	OrderItemPaid      OrderItemStatus = 2
	OrderItemShipped   OrderItemStatus = 3
	OrderItemCancelled OrderItemStatus = 4
	OrderItemCompleted OrderItemStatus = 5
)

type Order struct {
	ID              uuid.UUID   `json:"id" gorm:"primaryKey"`
	UserID          uuid.UUID   `json:"user_id"`
	TotalPrice      float64     `json:"total_price"`
	Status          OrderStatus `json:"status"`
	ShippingAddress string      `json:"shipping_address"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID              uuid.UUID       `json:"id" gorm:"primaryKey"`
	OrderID         uuid.UUID       `json:"order_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	SellerID        uuid.UUID       `json:"seller_id"`
	Quantity        int             `json:"quantity"`
	ProductName     string          `json:"product_name"`
	ProductImageUrl string          `json:"product_image_url"`
	UnitPrice       float64         `json:"unit_price"`
	Status          OrderItemStatus `json:"status"`
}
