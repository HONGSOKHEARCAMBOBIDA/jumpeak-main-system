package response

import (
	"mysql/model"
	"mysql/model/base"
)

type CompanyResponse struct {
	base.ModelBase
	Name           string              `gorm:"type:varchar(150);not null" json:"name"`
	BaseCurrency   string              `gorm:"type:char(3);not null;default:USD" json:"base_currency"`
	Status         model.CompanyStatus `gorm:"type:enum('ACTIVE','SUSPENDED');not null;default:ACTIVE" json:"status"`
	BranchResponse []BranchResponse    `json:"branches" gorm:"-"`
}
