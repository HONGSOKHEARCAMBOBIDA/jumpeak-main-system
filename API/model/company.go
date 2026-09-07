package model

import (
	"mysql/model/base"
	"time"
)

type CompanyStatus string

const (
	CompanyStatusActive    CompanyStatus = "ACTIVE"
	CompanyStatusSuspended CompanyStatus = "SUSPENDED"
)

type Company struct {
	base.ModelBase
	Name         string        `gorm:"type:varchar(150);not null" json:"name"`
	BaseCurrency string        `gorm:"type:char(3);not null;default:USD" json:"base_currency"`
	Status       CompanyStatus `gorm:"type:enum('ACTIVE','SUSPENDED');not null;default:ACTIVE" json:"status"`
	CreatedAt    time.Time     `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time    `gorm:"index" json:"deleted_at,omitempty"`
}
