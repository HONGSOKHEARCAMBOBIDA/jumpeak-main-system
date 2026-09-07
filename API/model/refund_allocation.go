package model

import (
	"mysql/model/base"
	"time"
)

type RefundAllocation struct {
	base.ModelBase
	RefundID  uint64    `gorm:"not null;index;uniqueIndex:uq_refund_invoice" json:"refund_id"`
	InvoiceID uint64    `gorm:"not null;index;uniqueIndex:uq_refund_invoice" json:"invoice_id"`
	Amount    float64   `gorm:"type:decimal(18,2);not null" json:"amount"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
