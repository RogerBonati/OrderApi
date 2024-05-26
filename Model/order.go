package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderId     uint64     `json:"order_id"`
	CustomerId  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	itemID   uuid.UUID `json:"item_id,omitempty"`
	Quantity uint      `json:"quantity,omitempty"`
	Price    float32   `json:"price,omitempty"`
}
