package response

type DebtAdjustmentResponse struct {
	ID            uint64  `json:"id"`
	CustomerID    uint64  `json:"customer_id"`
	CustomerName  string  `json:"customer_name"`
	InvoiceID     *uint64 `json:"invoice_id,omitempty"`
	InvoiceNumber string  `gorm:"type:varchar(30);not null;uniqueIndex:uq_invoice_number" json:"invoice_number"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason"`
	ApprovedBy    uint64  `json:"approved_by"`
	CreatedAt     string  `json:"created_at"`
}
