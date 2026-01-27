package domain

import "github.com/google/uuid"

type OrderItemReserve struct {
	ProductID uuid.UUID
	Quantity  int
}
