package model

import (
	"mysql/model/base"
	"time"
)

type DebtAdjustmentType string

const (
	DebtAdjustmentTypeWriteOff   DebtAdjustmentType = "WRITE_OFF"
	DebtAdjustmentTypeCorrection DebtAdjustmentType = "CORRECTION"
	DebtAdjustmentTypeDiscount   DebtAdjustmentType = "DISCOUNT"
)

type DebtAdjustment struct {
	base.ModelBase
	CompanyID  uint64             `gorm:"not null;index" json:"company_id"`
	CustomerID uint64             `gorm:"not null;index" json:"customer_id"`
	InvoiceID  *uint64            `gorm:"index" json:"invoice_id,omitempty"`
	Type       DebtAdjustmentType `gorm:"type:enum('WRITE_OFF','CORRECTION','DISCOUNT');not null" json:"type"`
	Amount     float64            `gorm:"type:decimal(18,2);not null" json:"amount"`
	Reason     string             `gorm:"type:varchar(255);not null" json:"reason"`
	ApprovedBy uint64             `gorm:"not null;index" json:"approved_by"`
	CreatedAt  time.Time          `gorm:"not null;autoCreateTime" json:"created_at"`
}
