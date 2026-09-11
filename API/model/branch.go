package model

import (
	"mysql/model/base"
	"time"

	"gorm.io/gorm"
)

type BranchStatus string

const (
	BranchStatusActive   BranchStatus = "ACTIVE"
	BranchStatusInactive BranchStatus = "INACTIVE"
)

type Branch struct {
	base.ModelBase
	CompanyID uint64         `gorm:"not null;index" json:"company_id"`
	Name      string         `gorm:"type:varchar(120);not null" json:"name"`
	Code      string         `gorm:"type:varchar(30);not null" json:"code"`
	Address   *string        `gorm:"type:varchar(255)" json:"address,omitempty"`
	Phone     string         `json:"phone" gorm:"column:phone"`
	Status    BranchStatus   `gorm:"type:enum('ACTIVE','INACTIVE');not null;default:ACTIVE" json:"status"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
