package request

import "mysql/model"

type BranchRequestCreate struct {
	CompanyID uint64  `gorm:"not null;index" json:"company_id"`
	Name      string  `gorm:"type:varchar(120);not null" json:"name"`
	Address   *string `gorm:"type:varchar(255)" json:"address,omitempty"`
	Phone     string  `json:"phone" gorm:"column:phone"`
}

type BranchRequestUpdate struct {
	CompanyID uint64             `gorm:"not null;index" json:"company_id"`
	Name      string             `gorm:"type:varchar(120);not null" json:"name"`
	Address   *string            `gorm:"type:varchar(255)" json:"address,omitempty"`
	Phone     string             `json:"phone" gorm:"column:phone"`
	Status    model.BranchStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
}
