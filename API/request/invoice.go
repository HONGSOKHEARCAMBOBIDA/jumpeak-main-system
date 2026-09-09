package request

import "time"

type InvoiceItemInput struct {
	ProductID      *uint64 `json:"product_id,omitempty"`
	Description    string  `json:"description"`
	Quantity       float64 `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountAmount float64 `json:"discount_amount"`
}

type InvoiceRequestCreate struct {
	CustomerID         uint64             `json:"customer_id"`
	BranchID           uint64             `json:"branch_id"`
	InvoiceDate        time.Time          `json:"invoice_date"`
	DueDate            time.Time          `json:"due_date"`
	CurrencyCode       string             `json:"currency_code"`
	ExchangeRateToBase float64            `json:"exchange_rate_to_base"`
	Items              []InvoiceItemInput `json:"items"`
}

type InvoiceRequestCancel struct {
	Reason string `json:"reason"`
}
