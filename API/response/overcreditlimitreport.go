package response

import "mysql/model"

type OverCreditLimitReportResponse struct {
	ID                  uint64               `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID           uint64               `gorm:"not null;index;uniqueIndex:uq_customer_code" json:"company_id"`
	CompanyName         string               `json:"company_name"`
	CompanyCurrency     string               `json:"company_currency"`
	BranchID            *uint64              `gorm:"index" json:"branch_id,omitempty"`
	BranchCode          string               `json:"branch_code"`
	BranchName          string               `json:"branch_name"`
	CustomerCode        string               `gorm:"type:varchar(30);not null;uniqueIndex:uq_customer_code" json:"customer_code"`
	Name                string               `gorm:"type:varchar(150);not null" json:"name"`
	Phone               *string              `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Address             *string              `gorm:"type:varchar(255)" json:"address,omitempty"`
	Notes               *string              `json:"notes" gorm:"column:notes"`
	CreditLimit         float64              `gorm:"type:decimal(18,2);not null;default:0.00" json:"credit_limit"`
	CreditLimitEnforced bool                 `gorm:"not null;default:false" json:"credit_limit_enforced"`
	CurrentOutstanding  float64              `gorm:"type:decimal(18,2);not null;default:0.00" json:"current_outstanding"`
	Status              model.CustomerStatus `gorm:"type:enum('ACTIVE','INACTIVE','BLACKLISTED');not null;default:ACTIVE" json:"status"`
}
