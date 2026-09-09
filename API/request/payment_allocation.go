package request

type PaymentAllocation struct {
	InvoiceID uint64  `gorm:"not null;index;uniqueIndex:uq_alloc_payment_invoice" json:"invoice_id"`
	Amount    float64 `gorm:"type:decimal(18,2);not null" json:"amount"`
}
