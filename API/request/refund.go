package request

type RefundAllocationInput struct {
	InvoiceID uint64  `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type RefundRequestCreate struct {
	PaymentID   uint64                  `json:"payment_id"`
	Amount      float64                 `json:"amount"`
	Reason      string                  `json:"reason"`
	RefundedAt  string                  `json:"refunded_at"`
	Allocations []RefundAllocationInput `json:"allocations"`
}
