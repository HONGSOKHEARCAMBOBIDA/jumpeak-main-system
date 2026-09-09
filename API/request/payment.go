package request

import "time"

type PaymentAllocationInput struct {
	InvoiceID uint64  `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type PaymentRequestCreate struct {
	CustomerID         uint64                   `json:"customer_id"`
	BranchID           *uint64                  `json:"branch_id,omitempty"`
	PaymentDate        time.Time                `json:"payment_date"`
	CurrencyCode       string                   `json:"currency_code"`
	ExchangeRateToBase float64                  `json:"exchange_rate_to_base"`
	Amount             float64                  `json:"amount"`
	Method             string                   `json:"method"`
	ReferenceNumber    *string                  `json:"reference_number,omitempty"`
	Note               *string                  `json:"note,omitempty"`
	Allocations        []PaymentAllocationInput `json:"allocations"`
}

type PaymentRequestVoid struct {
	Reason string `json:"reason"`
}
