package response

import "time"

type CustomerLedgerResponse struct {
	ID             uint64    `json:"id"`
	CustomerID     uint64    `json:"customer_id"`
	CustomerName   string    `json:"customer_name"`
	EntryDate      time.Time `json:"entry_date"`
	ReferenceType  string    `json:"reference_type"`
	ReferenceID    uint64    `json:"reference_id"`
	Description    string    `json:"description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	RunningBalance float64   `json:"running_balance"`
}
