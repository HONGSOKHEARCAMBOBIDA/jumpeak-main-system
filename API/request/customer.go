package request

import "mysql/model"

type CustomerRequestCreate struct {
	Name        string  `gorm:"type:varchar(150);not null" json:"name"`
	Phone       *string `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Address     *string `gorm:"type:varchar(255)" json:"address,omitempty"`
	Notes       *string `gorm:"type:text" json:"notes,omitempty"`
	CreditLimit float64 `gorm:"type:decimal(18,2);not null;default:0.00" json:"credit_limit"`
}

type CustomerRequestUpdate struct {
	Name                string               `gorm:"type:varchar(150);not null" json:"name"`
	Phone               *string              `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Address             *string              `gorm:"type:varchar(255)" json:"address,omitempty"`
	Notes               *string              `gorm:"type:text" json:"notes,omitempty"`
	CreditLimit         float64              `gorm:"type:decimal(18,2);not null;default:0.00" json:"credit_limit"`
	CreditLimitEnforced bool                 `gorm:"not null;default:false" json:"credit_limit_enforced"`
	Status              model.CustomerStatus `gorm:"type:enum('ACTIVE','INACTIVE','BLACKLISTED');not null;default:ACTIVE" json:"status"`
}
