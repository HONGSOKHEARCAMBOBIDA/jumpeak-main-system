package response

import (
	"mysql/model"
	"mysql/model/base"
)

type UserResponse struct {
	base.ModelBase
	BranchID        *uint64                `gorm:"index" json:"branch_id,omitempty"`
	Name            string                 `gorm:"type:varchar(120);not null" json:"name"`
	Email           string                 `gorm:"type:varchar(190);not null;uniqueIndex:uq_user_email" json:"email"`
	RoleID          int64                  `gorm:"not null;index" json:"role_id"`
	RoleName        string                 `json:"role_name"`
	RoleDisplayName string                 `json:"role_display_name"`
	ManageBranch    model.UserManageBranch `gorm:"type:enum('ONE','MULTIPLE','ALL');not null;default:ONE" json:"manage_branch"`
	Status          model.UserStatus       `gorm:"type:enum('ACTIVE','DISABLED');not null;default:ACTIVE" json:"status"`
	ManageBranchID  []ManageBranchID       `json:"branch_ids" gorm:"-"`
}

type ManageBranchID struct {
	base.ModelBase
	UserID   int `json:"user_id"`
	BranchID int `json:"branch_id"`
}
