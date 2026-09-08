package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/constant/apperror"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
)

type CustomerService interface {
	Create(ctx context.Context, userID int, input request.CustomerRequestCreate) error
	Update(ctx context.Context, id int, userID int, input request.CustomerRequestUpdate) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerResponse, *model.PaginationMetadata, error)
}

type customerservice struct {
	db *gorm.DB
}

func NewCustomerService() CustomerService {
	return &customerservice{
		db: config.DB,
	}
}

func (s *customerservice) Create(ctx context.Context, userID int, input request.CustomerRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Customer{
			CompanyID:           user.CompanyID,
			BranchID:            user.BranchID,
			Name:                input.Name,
			Phone:               input.Phone,
			Address:             input.Address,
			Notes:               input.Notes,
			CreditLimit:         input.CreditLimit,
			CreditLimitEnforced: true,
			Status:              model.CustomerStatusActive,
			CreatedBy:           &userID,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}
		newdata.CustomerCode = helper.GenerateCode("CUSTOMER-", uint(newdata.ID))
		if err := tx.Save(&newdata).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "faild to update code", nil)
		}
		return nil
	})
	return err
}

func (s *customerservice) Update(ctx context.Context, id int, userID int, input request.CustomerRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Customer
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Name = input.Name
		data.Phone = input.Phone
		data.Address = input.Address
		data.Notes = input.Notes
		data.CreditLimit = input.CreditLimit
		data.CreditLimitEnforced = input.CreditLimitEnforced
		data.Status = input.Status
		data.UpdatedBy = &userID
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *customerservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.CustomerResponse
	var total int64
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("customers c").
			Joins("LEFT JOIN branches b ON b.id = c.branch_id").
			Joins("LEFT JOIN companies cp ON cp.id = c.company_id")

	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("c.name LIKE ?", "%"+v+"%")
		}
		if v, ok := filter["branch_id"]; ok && v != "" {
			tx = tx.Where("b.id = ?", v)
		}
		if v, ok := filter["company_id"]; ok && v != "" {
			tx = tx.Where("cp.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count customer: %w", err)
	}

	if total == 0 {
		return []response.CustomerResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		c.id AS id,
		cp.id AS company_id,
		cp.name AS company_name,
		cp.base_currency AS company_currency,
		b.id AS branch_id,
		b.code AS branch_code,
		b.name AS branch_name,
		c.customer_code AS customer_code,
		c.name AS name,
		c.phone AS phone,
		c.address AS address,
		c.notes AS notes,
		c.credit_limit AS credit_limit,
		c.credit_limit_enforced AS credit_limit_enforced,
		c.current_outstanding AS current_outstanding,
		c.status AS status
	`)

	dataQuery = helper.ApplyAccessFilter(dataQuery, s.db, user.Role, user)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch companies: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil

}
