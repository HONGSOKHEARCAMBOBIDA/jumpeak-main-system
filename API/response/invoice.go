package response

import "time"

type InvoiceItemResponse struct {
	ID             uint64  `json:"id"`
	ProductID      *uint64 `json:"product_id,omitempty"`
	Description    string  `json:"description"`
	Quantity       float64 `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountAmount float64 `json:"discount_amount"`
	Subtotal       float64 `json:"subtotal"`
}

type InvoiceResponse struct {
	ID                 uint64    `json:"id"`
	CompanyID          uint64    `json:"company_id"`
	BranchID           uint64    `json:"branch_id"`
	CustomerID         uint64    `json:"customer_id"`
	CustomerName       string    `json:"customer_name"`
	InvoiceNumber      string    `json:"invoice_number"`
	InvoiceDate        time.Time `json:"invoice_date"`
	DueDate            time.Time `json:"due_date"`
	CurrencyCode       string    `json:"currency_code"`
	ExchangeRateToBase float64   `json:"exchange_rate_to_base"`
	TotalAmount        float64   `json:"total_amount"`
	PaidAmount         float64   `json:"paid_amount"`
	OutstandingAmount  float64   `json:"outstanding_amount"`
	Status             string    `json:"status"`
}
