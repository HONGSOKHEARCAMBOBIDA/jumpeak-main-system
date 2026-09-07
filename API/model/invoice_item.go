package model

import (
	"mysql/model/base"
	"time"
)

type InvoiceItem struct {
	base.ModelBase
	InvoiceID      uint64    `gorm:"not null;index" json:"invoice_id"`
	ProductID      *uint64   `gorm:"index" json:"product_id,omitempty"`
	Description    string    `gorm:"type:varchar(255);not null" json:"description"`
	Quantity       float64   `gorm:"type:decimal(12,2);not null" json:"quantity"`
	UnitPrice      float64   `gorm:"type:decimal(18,2);not null" json:"unit_price"`
	DiscountAmount float64   `gorm:"type:decimal(18,2);not null;default:0.00" json:"discount_amount"`
	Subtotal       float64   `gorm:"type:decimal(18,2);not null" json:"subtotal"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}
