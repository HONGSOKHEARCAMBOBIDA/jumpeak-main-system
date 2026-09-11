package response

// type InvoiceItemResponse struct {
// 	ID             uint64  `json:"id"`
// 	ProductID      *uint64 `json:"product_id,omitempty"`
// 	Description    string  `json:"description"`
// 	Quantity       float64 `json:"quantity"`
// 	UnitPrice      float64 `json:"unit_price"`
// 	DiscountAmount float64 `json:"discount_amount"`
// 	Subtotal       float64 `json:"subtotal"`
// }

type InvoiceResponse struct {
	ID                  uint64                `json:"id"`
	CompanyID           uint64                `json:"company_id"`
	CompanyName         string                `json:"company_Name"`
	BranchID            uint64                `json:"branch_id" gorm:"column:branch_id"`
	BranchName          string                `json:"branch_Name"`
	BranchPhone         string                `json:"branch_phone"`
	CustomerID          uint64                `json:"customer_id"`
	CustomerName        string                `json:"customer_name"`
	InvoiceNumber       string                `json:"invoice_number"`
	InvoiceDate         string                `json:"invoice_date"`
	DueDate             string                `json:"due_date"`
	CurrencyCode        string                `json:"currency_code"`
	ExchangeRateToBase  float64               `json:"exchange_rate_to_base"`
	TotalAmount         float64               `json:"total_amount"`
	PaidAmount          float64               `json:"paid_amount"`
	OutstandingAmount   float64               `json:"outstanding_amount"`
	Status              string                `json:"status"`
	CancelReason        *string               `gorm:"type:varchar(255)" json:"cancel_reason,omitempty"`
	InvoiceItemResponse []InvoiceItemResponse `json:"invoice_item" gorm:"-"`
}
