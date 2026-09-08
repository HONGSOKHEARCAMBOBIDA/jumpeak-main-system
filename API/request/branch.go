package request

import "mysql/model"

type BranchRequestCreate struct {
	CompanyID uint64  `gorm:"not null;index" json:"company_id"`
	Name      string  `gorm:"type:varchar(120);not null" json:"name"`
	Address   *string `gorm:"type:varchar(255)" json:"address,omitempty"`
}

type BranchRequestUpdate struct {
	CompanyID uint64             `gorm:"not null;index" json:"company_id"`
	Name      string             `gorm:"type:varchar(120);not null" json:"name"`
	Address   *string            `gorm:"type:varchar(255)" json:"address,omitempty"`
	Status    model.BranchStatus `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
}
