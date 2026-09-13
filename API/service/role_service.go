package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type RoleService interface {
	CreateRoleHasPermission(ctx context.Context, input request.CreateRolePermissionInput) error
	DeleteRoleHasPermission(ctx context.Context, input request.DeleteRolePermissionsInput) error
	GetRolePermission(ctx context.Context, id int, userID int, pf request.Pagination) ([]response.PermissionWithAssignedRole, *model.PaginationMetadata, error)
	UpdateRole(ctx context.Context, id int, input request.RoleRequestUpdate) error
}

type roleservice struct {
	db *gorm.DB
}

func NewRoleService() RoleService {
	return &roleservice{
		db: config.DB,
	}
}

func (s *roleservice) CreateRoleHasPermission(ctx context.Context, input request.CreateRolePermissionInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var rolePermissions []model.RoleHasPermission
	for _, permissionID := range input.PermissionIDs {
		rolePermissions = append(rolePermissions, model.RoleHasPermission{
			RoleID:       uint(input.RoleID),
			PermissionID: uint(permissionID),
		})
	}

	if err := tx.Create(&rolePermissions).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *roleservice) DeleteRoleHasPermission(ctx context.Context, input request.DeleteRolePermissionsInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.
		Where("role_id = ? AND permission_id IN ?", input.RoleID, input.PermissionIDs).
		Delete(&model.RoleHasPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *roleservice) GetRolePermission(ctx context.Context, id int, userID int, pf request.Pagination) ([]response.PermissionWithAssignedRole, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var permissions []response.PermissionWithAssignedRole
	var total int64
	if id <= 0 {
		return nil, nil, fmt.Errorf("invalid role id: %d", id)
	}

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("permission p").
			Joins("LEFT JOIN role_permission rp ON p.id = rp.permission_id").
			Where("rp.role_id = ?", id)
	}

	if err := base().Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count permission: %w", err)
	}

	if total == 0 {
		return []response.PermissionWithAssignedRole{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := base().Select(`
		p.id AS id,
		p.name AS name,
		p.display_name AS display_name,
		CASE 
			WHEN rp.permission_id IS NULL THEN false
			ELSE true
		END AS assigned
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&permissions).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch permission: %w", err)
	}
	return permissions, helper.BuildPaginationMeta(pf, total), nil
}

func (s *roleservice) UpdateRole(ctx context.Context, id int, input request.RoleRequestUpdate) error {
	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.DisPlayName != nil {
		updates["display_name"] = *input.DisPlayName
	}
	if len(updates) == 0 {
		return errors.New(" no field to update")
	}
	result := s.db.WithContext(ctx).Model(&model.Role{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
