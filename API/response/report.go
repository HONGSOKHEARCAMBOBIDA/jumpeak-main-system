package response

type CustomerOutstandingReport struct {
	CompanyID               int                       `json:"company_id"`
	CompanyName             string                    `json:"company_name"`
	TotalAmount             float64                   `json:"total_amount"`
	Currency                string                    `json:"currency"`
	BranchOutStandingReport []BranchOutStandingReport `json:"branch_outstanding" gorm:"-"`
}

type BranchOutStandingReport struct {
	CompanyID   int     `json:"company_id"`
	BranchName  string  `json:"Branch_name"`
	TotalAmount float64 `json:"total_amount"`
	Currency    string  `json:"currency"`
}
