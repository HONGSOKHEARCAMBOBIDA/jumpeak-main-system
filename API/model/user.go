package model

import (
	"mysql/model/base"
	"time"

	"gorm.io/gorm"
)

type UserManageBranch string

const (
	UserManageBranchOne      UserManageBranch = "ONE"
	UserManageBranchMultiple UserManageBranch = "MULTIPLE"
	UserManageBranchAll      UserManageBranch = "ALL"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusDisabled UserStatus = "DISABLED"
)

type User struct {
	base.ModelBase
	CompanyID    uint64           `gorm:"not null;index;uniqueIndex:uq_user_email" json:"company_id"`
	BranchID     *uint64          `gorm:"index" json:"branch_id,omitempty"`
	Name         string           `gorm:"type:varchar(120);not null" json:"name"`
	Email        string           `gorm:"type:varchar(190);not null;uniqueIndex:uq_user_email" json:"email"`
	PasswordHash string           `gorm:"type:varchar(255);not null" json:"password_hash"`
	RoleID       int64            `gorm:"not null;index" json:"role_id"`
	ManageBranch UserManageBranch `gorm:"type:enum('ONE','MULTIPLE','ALL');not null;default:ONE" json:"manage_branch"`
	Status       UserStatus       `gorm:"type:enum('ACTIVE','DISABLED');not null;default:ACTIVE" json:"status"`
	CreatedAt    time.Time        `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty"`
	Role         Role
}
