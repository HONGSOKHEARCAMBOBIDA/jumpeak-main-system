package request

import "mysql/model"

type CompanyRequestCreate struct {
	Name         string              `gorm:"type:varchar(150);not null" json:"name"`
	BaseCurrency string              `gorm:"type:char(3);not null;default:USD" json:"base_currency"`
	Status       model.CompanyStatus `gorm:"type:enum('ACTIVE','SUSPENDED');not null;default:ACTIVE" json:"status"`
}

type CompanyRequestUpdate struct {
	Name         string              `gorm:"type:varchar(150);not null" json:"name"`
	BaseCurrency string              `gorm:"type:char(3);not null;default:USD" json:"base_currency"`
	Status       model.CompanyStatus `gorm:"type:enum('ACTIVE','SUSPENDED');not null;default:ACTIVE" json:"status"`
}
