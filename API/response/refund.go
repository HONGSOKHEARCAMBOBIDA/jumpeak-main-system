package response

type RefundResponse struct {
	ID            uint64  `json:"id"`
	CompanyName   string  `json:"company_Name"`
	BranchID      int     `json:"branch_id"`
	BranchName    string  `json:"branch_Name"`
	BranchCode    string  `json:"branch_Code"`
	CustomerName  string  `json:"customer_name"`
	Currency      string  `json:"currency"`
	PaymentNumber string  `gorm:"type:varchar(30);not null;uniqueIndex:uq_payment_number" json:"payment_number"`
	PaymentID     uint64  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason"`
	RefundedAt    string  `json:"refunded_at"`
	CreateBy      string  `json:"create_by" gorm:"column:create_by"`
}
