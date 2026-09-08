package model

import (
	"time"

	"gorm.io/gorm"
)

type CustomerStatus string

const (
	CustomerStatusActive      CustomerStatus = "ACTIVE"
	CustomerStatusInactive    CustomerStatus = "INACTIVE"
	CustomerStatusBlacklisted CustomerStatus = "BLACKLISTED"
)

type Customer struct {
	ID                  uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID           uint64         `gorm:"not null;index;uniqueIndex:uq_customer_code" json:"company_id"`
	BranchID            *uint64        `gorm:"index" json:"branch_id,omitempty"`
	CustomerCode        string         `gorm:"type:varchar(30);not null;uniqueIndex:uq_customer_code" json:"customer_code"`
	Name                string         `gorm:"type:varchar(150);not null" json:"name"`
	Phone               *string        `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Address             *string        `gorm:"type:varchar(255)" json:"address,omitempty"`
	Notes               *string        `gorm:"type:text" json:"notes,omitempty"`
	CreditLimit         float64        `gorm:"type:decimal(18,2);not null;default:0.00" json:"credit_limit"`
	CreditLimitEnforced bool           `gorm:"not null;default:false" json:"credit_limit_enforced"`
	CurrentOutstanding  float64        `gorm:"type:decimal(18,2);not null;default:0.00" json:"current_outstanding"`
	Status              CustomerStatus `gorm:"type:enum('ACTIVE','INACTIVE','BLACKLISTED');not null;default:ACTIVE" json:"status"`
	CreatedBy           *int           `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy           *int           `gorm:"index" json:"updated_by,omitempty"`
	CreatedAt           time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
