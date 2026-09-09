package response

import "time"

type DebtAdjustmentResponse struct {
	ID         uint64    `json:"id"`
	CustomerID uint64    `json:"customer_id"`
	InvoiceID  *uint64   `json:"invoice_id,omitempty"`
	Type       string    `json:"type"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
	ApprovedBy uint64    `json:"approved_by"`
	CreatedAt  time.Time `json:"created_at"`
}
