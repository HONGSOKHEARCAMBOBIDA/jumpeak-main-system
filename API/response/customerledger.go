package response

type CustomerLedgerResponse struct {
	ID             uint64  `json:"id"`
	CompanyName    string  `json:"company_Name"`
	BranchName     string  `json:"branch_Name"`
	BranchCode     string  `json:"branch_Code"`
	Currency       string  `json:"currency"`
	CustomerID     uint64  `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	EntryDate      string  `json:"entry_date"`
	ReferenceType  string  `json:"reference_type"`
	ReferenceID    uint64  `json:"reference_id"`
	Description    string  `json:"description"`
	Debit          float64 `json:"debit"`
	Credit         float64 `json:"credit"`
	RunningBalance float64 `json:"running_balance"`
}
