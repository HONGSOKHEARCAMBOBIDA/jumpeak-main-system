package model

import "time"

type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "CASH"
	PaymentMethodBank   PaymentMethod = "BANK"
	PaymentMethodABA    PaymentMethod = "ABA"
	PaymentMethodACLEDA PaymentMethod = "ACLEDA"
	PaymentMethodOther  PaymentMethod = "OTHER"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusVoided    PaymentStatus = "VOIDED"
	PaymentStatusRefund    PaymentStatus = "REFUND"
)

type Payment struct {
	ID                 uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID          uint64        `gorm:"not null;index;uniqueIndex:uq_payment_number" json:"company_id"`
	BranchID           *uint64       `gorm:"index" json:"branch_id,omitempty"`
	CustomerID         uint64        `gorm:"not null;index" json:"customer_id"`
	PaymentNumber      string        `gorm:"type:varchar(30);not null;uniqueIndex:uq_payment_number" json:"payment_number"`
	PaymentDate        time.Time     `gorm:"type:date;not null" json:"payment_date"`
	CurrencyCode       string        `gorm:"type:char(3);not null" json:"currency_code"`
	ExchangeRateToBase float64       `gorm:"type:decimal(9,6);not null;default:1.000000" json:"exchange_rate_to_base"`
	Amount             float64       `gorm:"type:decimal(18,2);not null" json:"amount"`
	Method             PaymentMethod `gorm:"type:enum('CASH','BANK','ABA','ACLEDA','OTHER');not null" json:"method"`
	ReferenceNumber    *string       `gorm:"type:varchar(60)" json:"reference_number,omitempty"`
	Note               *string       `gorm:"type:varchar(255)" json:"note,omitempty"`
	Status             PaymentStatus `gorm:"type:enum('COMPLETED','VOIDED','REFUND');not null;default:COMPLETED" json:"status"`
	CreatedBy          *uint64       `gorm:"index" json:"created_by,omitempty"`
	CreatedAt          time.Time     `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time     `gorm:"not null;autoUpdateTime" json:"updated_at"`
}
