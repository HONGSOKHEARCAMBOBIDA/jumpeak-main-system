package model

import (
	"mysql/model/base"
	"time"
)

type Refund struct {
	base.ModelBase
	PaymentID  uint64    `gorm:"not null;index" json:"payment_id"`
	Amount     float64   `gorm:"type:decimal(18,2);not null" json:"amount"`
	Reason     string    `gorm:"type:varchar(255);not null" json:"reason"`
	RefundedAt time.Time `gorm:"not null" json:"refunded_at"`
	CreatedBy  *uint64   `gorm:"index" json:"created_by,omitempty"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
