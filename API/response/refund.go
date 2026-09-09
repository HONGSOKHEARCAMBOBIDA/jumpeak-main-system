package response

import "time"

type RefundResponse struct {
	ID         uint64    `json:"id"`
	PaymentID  uint64    `json:"payment_id"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
	RefundedAt time.Time `json:"refunded_at"`
}
