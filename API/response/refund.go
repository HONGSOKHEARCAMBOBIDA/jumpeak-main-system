package response

type RefundResponse struct {
	ID            uint64  `json:"id"`
	CustomerName  string  `json:"customer_name"`
	PaymentNumber string  `gorm:"type:varchar(30);not null;uniqueIndex:uq_payment_number" json:"payment_number"`
	PaymentID     uint64  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason"`
	RefundedAt    string  `json:"refunded_at"`
}
