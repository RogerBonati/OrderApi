package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderId     uint64     `json:"order_id,omitempty"`
	CustomerId  uuid.UUID  `json:"customer_id,omitempty"`
	LineItems   []LineItem `json:"line_items,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type LineItem struct {
	itemID   uuid.UUID `json:"item_id,omitempty"`
	Quantity uint      `json:"quantity,omitempty"`
	Price    float32   `json:"price,omitempty"`
}
