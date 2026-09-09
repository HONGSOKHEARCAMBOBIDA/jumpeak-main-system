package model

import (
	"mysql/model/base"
	"time"
)

type CustomerLedgerReferenceType string

const (
	CustomerLedgerReferenceInvoice    CustomerLedgerReferenceType = "INVOICE"
	CustomerLedgerReferencePayment    CustomerLedgerReferenceType = "PAYMENT"
	CustomerLedgerReferenceRefund     CustomerLedgerReferenceType = "REFUND"
	CustomerLedgerReferenceAdjustment CustomerLedgerReferenceType = "ADJUSTMENT"
)

type CustomerLedger struct {
	base.ModelBase
	CompanyID      uint64                      `gorm:"not null;index" json:"company_id"`
	CustomerID     uint64                      `gorm:"not null;index:idx_ledger_customer_date,priority:1" json:"customer_id"`
	EntryDate      time.Time                   `gorm:"type:date;not null;index:idx_ledger_customer_date,priority:2" json:"entry_date"`
	ReferenceType  CustomerLedgerReferenceType `gorm:"type:enum('INVOICE','PAYMENT','REFUND','ADJUSTMENT');not null" json:"reference_type"`
	ReferenceID    uint64                      `gorm:"not null" json:"reference_id"`
	Description    string                      `gorm:"type:varchar(255);not null" json:"description"`
	Debit          float64                     `gorm:"type:decimal(18,2);not null;default:0.00" json:"debit"`
	Credit         float64                     `gorm:"type:decimal(18,2);not null;default:0.00" json:"credit"`
	RunningBalance float64                     `gorm:"type:decimal(18,2);not null" json:"running_balance"`
	CreatedAt      time.Time                   `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (CustomerLedger) TableName() string {
	return "customer_ledger"
}
