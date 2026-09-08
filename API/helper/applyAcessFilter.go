package helper

import (
	"errors"
	"mysql/model"

	"gorm.io/gorm"
)

func CompanyFilter(query, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	if role.Level < 7 {
		return query.Where("c.id = ?", user.CompanyID)
	}
	return query
}

func UserFilter(query, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	if role.Level <= 1 {
		return query.Where("u.id = ?", user.ID)
	}
	return query
}

func ApplyAccessFilter(query, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	if role.Level > 1 && role.Level < 7 {
		switch user.ManageBranch {
		case model.UserManageBranchOne:
			return query.Where("b.id =?", user.BranchID)
		case model.UserManageBranchMultiple:
			var branchIDs []int
			if err := db.Model(&model.UserBranch{}).Where("user_id = ?", user.ID).Pluck("branch_id", &branchIDs).Error; err != nil {
				return query.Where("1 = 0") // or propagate the error, don't fail open
			}
			return query.Where("b.id IN ?", branchIDs)
		}
		return query
	} else if role.Level <= 1 {
		return query.Where("b.id =?", user.BranchID)
	}
	return query
}

func ApplyAccessGetRole(query, db *gorm.DB, role model.Role, user model.User) *gorm.DB {
	return query.Where("r.level <= ?", user.Role.Level)
}

func CanManageUser(actorRole model.Role, targetRole model.Role) error {
	if targetRole.Level >= actorRole.Level {
		return errors.New("you do not have permission to manage a user with a higher role level")
	}
	return nil
}
