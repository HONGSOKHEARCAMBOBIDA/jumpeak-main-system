package service

import (
	"context"
	"errors"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
)

type BranchService interface {
	Create(ctx context.Context, input request.BranchRequestCreate) error
	Update(ctx context.Context, id int, input request.BranchRequestUpdate) error
	GetBranchNoPagination(ctx context.Context, userID int, companyID int) ([]model.Branch, error)
}

type branchservice struct {
	db *gorm.DB
}

func NewBranchService() BranchService {
	return &branchservice{
		db: config.DB,
	}
}

func (s *branchservice) Create(ctx context.Context, input request.BranchRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Branch{
			CompanyID: input.CompanyID,
			Name:      input.Name,
			Address:   input.Address,
			Phone:     input.Phone,
			Status:    model.BranchStatusActive,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}
		newdata.Code = helper.GenerateCode("BRANCH-", uint(newdata.ID))
		if err := tx.Save(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "faild to update code", nil)
		}
		return nil
	})
	return err
}

func (s *branchservice) Update(ctx context.Context, id int, input request.BranchRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Branch
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.CompanyID = input.CompanyID
		data.Name = input.Name
		data.Address = input.Address
		data.Phone = input.Phone
		data.Status = input.Status
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *branchservice) GetBranchNoPagination(ctx context.Context, userID int, companyID int) ([]model.Branch, error) {
	var data []model.Branch
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, err
	}
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("branches b").
			Where("b.company_id = ?", companyID)
	}
	dataQuery := base().Select(`
		b.id AS id,
		b.name AS name
	`)
	dataQuery = helper.ApplyAccessFilter(dataQuery, s.db, user.Role, user)
	if err := dataQuery.Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
