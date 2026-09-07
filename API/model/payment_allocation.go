package model

import (
	"mysql/model/base"
	"time"
)

type PaymentAllocation struct {
	base.ModelBase
	PaymentID uint64    `gorm:"not null;index;uniqueIndex:uq_alloc_payment_invoice" json:"payment_id"`
	InvoiceID uint64    `gorm:"not null;index;uniqueIndex:uq_alloc_payment_invoice" json:"invoice_id"`
	Amount    float64   `gorm:"type:decimal(18,2);not null" json:"amount"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
