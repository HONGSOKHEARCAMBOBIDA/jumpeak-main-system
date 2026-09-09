package response

import (
	"mysql/model"
	"mysql/model/base"
)

type ProductResponse struct {
	base.ModelBase
	CompanyID       uint64              `gorm:"not null;index;uniqueIndex:uq_customer_code" json:"company_id"`
	CompanyName     string              `json:"company_name"`
	CompanyCurrency string              `json:"company_currency"`
	BranchName      string              `json:"branch_name"`
	Name            string              `gorm:"type:varchar(150);not null" json:"name"`
	Status          model.ProductStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
}
