package request

import "mysql/model"

type UserRequestCreate struct {
	BranchID     *uint64                `gorm:"index" json:"branch_id,omitempty"`
	Name         string                 `gorm:"type:varchar(120);not null" json:"name"`
	RoleID       int64                  `gorm:"not null;index" json:"role_id"`
	ManageBranch model.UserManageBranch `gorm:"type:enum('ONE','MULTIPLE','ALL');not null;default:ONE" json:"manage_branch"`
	BranchIDs    *[]uint64              `json:"branch_ids"`
}

type UserRequestUpdate struct {
	BranchID     *uint64                `gorm:"index" json:"branch_id,omitempty"`
	Name         string                 `gorm:"type:varchar(120);not null" json:"name"`
	RoleID       int64                  `gorm:"not null;index" json:"role_id"`
	ManageBranch model.UserManageBranch `gorm:"type:enum('ONE','MULTIPLE','ALL');not null;default:ONE" json:"manage_branch"`
	Status       model.UserStatus       `json:"status"`
	BranchIDs    *[]uint64              `json:"branch_ids"`
}
