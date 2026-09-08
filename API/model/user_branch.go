package model

import "mysql/model/base"

type UserBranch struct {
	base.ModelBase
	UserID   uint64 `gorm:"not null;index" json:"user_id"`
	BranchID uint64 `gorm:"not null;index" json:"branch_id"`
}

func (UserBranch) TableName() string {
	return "user_branches"
}
