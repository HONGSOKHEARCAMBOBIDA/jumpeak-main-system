package model

import (
	"mysql/model/base"
	"time"

	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusOpen          InvoiceStatus = "OPEN"
	InvoiceStatusPartiallyPaid InvoiceStatus = "PARTIALLY_PAID"
	InvoiceStatusPaid          InvoiceStatus = "PAID"
	InvoiceStatusCancelled     InvoiceStatus = "CANCELLED"
	InvoiceStatusWrittenOff    InvoiceStatus = "WRITTEN_OFF"
)

type Invoice struct {
	base.ModelBase
	CompanyID          uint64         `gorm:"not null;index;uniqueIndex:uq_invoice_number" json:"company_id"`
	BranchID           uint64         `gorm:"not null;index" json:"branch_id"`
	CustomerID         uint64         `gorm:"not null;index" json:"customer_id"`
	InvoiceNumber      string         `gorm:"type:varchar(30);not null;uniqueIndex:uq_invoice_number" json:"invoice_number"`
	InvoiceDate        time.Time      `gorm:"type:date;not null" json:"invoice_date"`
	DueDate            time.Time      `gorm:"type:date;not null" json:"due_date"`
	CurrencyCode       string         `gorm:"type:char(3);not null" json:"currency_code"`
	ExchangeRateToBase float64        `gorm:"type:decimal(9,6);not null;default:1.000000" json:"exchange_rate_to_base"`
	TotalAmount        float64        `gorm:"type:decimal(18,2);not null" json:"total_amount"`
	PaidAmount         float64        `gorm:"type:decimal(18,2);not null;default:0.00" json:"paid_amount"`
	OutstandingAmount  float64        `gorm:"type:decimal(18,2);not null" json:"outstanding_amount"`
	Status             InvoiceStatus  `gorm:"type:enum('OPEN','PARTIALLY_PAID','PAID','CANCELLED','WRITTEN_OFF');not null;default:OPEN" json:"status"`
	CancelReason       *string        `gorm:"type:varchar(255)" json:"cancel_reason,omitempty"`
	CancelledAt        *time.Time     `json:"cancelled_at,omitempty"`
	CancelledBy        *uint64        `gorm:"index" json:"cancelled_by,omitempty"`
	PaidAt             *time.Time     `json:"paid_at,omitempty"`
	CreatedBy          *uint64        `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy          *uint64        `gorm:"index" json:"updated_by,omitempty"`
	CreatedAt          time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
