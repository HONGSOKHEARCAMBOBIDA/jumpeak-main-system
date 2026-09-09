package request

type DebtAdjustmentRequestCreate struct {
	CustomerID uint64  `json:"customer_id"`
	InvoiceID  *uint64 `json:"invoice_id,omitempty"`
	Type       string  `json:"type"` // WRITE_OFF, CORRECTION, DISCOUNT
	Amount     float64 `json:"amount"`
	Reason     string  `json:"reason"`
}
