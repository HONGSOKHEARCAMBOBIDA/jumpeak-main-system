package model

import (
	"mysql/model/base"
	"time"

	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusInactive ProductStatus = "INACTIVE"
)

type Product struct {
	base.ModelBase
	CompanyID    uint64         `gorm:"not null;index" json:"company_id"`
	Name         string         `gorm:"type:varchar(150);not null" json:"name"`
	DefaultPrice float64        `json:"default_price" gorm:"column:default_price"`
	Status       ProductStatus  `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
	CreatedAt    time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
