package response

import (
	"mysql/model"
	"mysql/model/base"
)

type BranchResponse struct {
	base.ModelBase
	CompanyID    uint64             `gorm:"not null;index" json:"company_id"`
	Name         string             `gorm:"type:varchar(120);not null" json:"name"`
	Code         string             `gorm:"type:varchar(30);not null" json:"code"`
	Address      *string            `gorm:"type:varchar(255)" json:"address,omitempty"`
	Phone        string             `json:"phone" gorm:"column:phone"`
	Status       model.BranchStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
	UserResponse []UserResponse     `json:"users" gorm:"-"`
}
