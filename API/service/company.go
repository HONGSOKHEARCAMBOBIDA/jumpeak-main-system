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

type CompanyService interface {
	Create(ctx context.Context, input request.CompanyRequestCreate) error
	Get(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.CompanyResponse, *model.PaginationMetadata, error)
	Update(ctx context.Context, id int, input request.CompanyRequestUpdate) error
}

type companyservice struct {
	db *gorm.DB
}

func NewCompanyService() CompanyService {
	return &companyservice{
		db: config.DB,
	}
}

func (s *companyservice) Create(ctx context.Context, input request.CompanyRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Company{
			Name:         input.Name,
			BaseCurrency: input.BaseCurrency,
			Status:       input.Status,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}
		return nil
	})
	return err
}

func (s *companyservice) Update(ctx context.Context, id int, input request.CompanyRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Company
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Name = input.Name
		data.BaseCurrency = input.BaseCurrency
		data.Status = input.Status
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update student", nil)
		}
		return nil
	})
	return err
}

func (s *companyservice) Get(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.CompanyResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var companies []response.CompanyResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("companies c")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("c.name LIKE ?", "%"+v+"%")
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count companies: %w", err)
	}

	if total == 0 {
		return []response.CompanyResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		c.id AS id,
		c.name AS name,
		c.base_currency AS base_currency,
		c.status AS status
	`)

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&companies).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch companies: %w", err)
	}

	companyIDs := make([]uint64, 0, len(companies))
	for _, c := range companies {
		companyIDs = append(companyIDs, uint64(c.ID))
	}

	var branches []response.BranchResponse
	if err := s.db.WithContext(ctx).Table("branches b").
		Where("b.company_id IN ?", companyIDs).
		Select(`
			b.id AS id,
			b.company_id AS company_id,
			b.name AS name,
			b.code AS code,
			b.address AS address,
			b.status AS status
		`).Scan(&branches).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch branches: %w", err)
	}

	branchesByCompany := make(map[uint64][]response.BranchResponse, len(companies))
	for _, b := range branches {
		branchesByCompany[b.CompanyID] = append(branchesByCompany[b.CompanyID], b)
	}

	for i := range companies {
		companies[i].BranchResponse = branchesByCompany[uint64(companies[i].ID)]
	}

	return companies, helper.BuildPaginationMeta(pf, total), nil
}
