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

type ProductService interface {
	Create(ctx context.Context, userID int, input request.ProductRequestCreate) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.ProductResponse, *model.PaginationMetadata, error)
	Update(ctx context.Context, id int, input request.ProductRequestUpdate) error
}

type productservice struct {
	db *gorm.DB
}

func NewProductService() ProductService {
	return &productservice{
		db: config.DB,
	}
}

func (s *productservice) Create(ctx context.Context, userID int, input request.ProductRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(input.ProductRequest) > 0 {
			newdata := make([]model.Product, 0, len(input.ProductRequest))
			for _, n := range input.ProductRequest {
				newdata = append(newdata, model.Product{
					CompanyID: user.CompanyID,
					Name:      n.Name,
					Status:    model.ProductStatusActive,
				})
			}
			if err := tx.Create(&newdata).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create proudct", nil)
			}
		}
		return nil
	})
	return err
}

func (s *productservice) Update(ctx context.Context, id int, input request.ProductRequestUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var data model.Product
		if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "classcurriculumn not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch classcurriculumn", nil)
		}
		data.Name = input.Name
		data.Status = input.Status
		if err := tx.Save(&data).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update product", nil)
		}
		return nil
	})
	return err
}

func (s *productservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.ProductResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.ProductResponse
	var total int64
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}
	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("products p").
			Joins("LEFT JOIN companies c ON c.id = p.company_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["name"]; ok && v != "" {
			tx = tx.Where("p.name LIKE ?", "%"+v+"%")
		}
		if v, ok := filter["company_id"]; ok && v != "" {
			tx = tx.Where("c.id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count product: %w", err)
	}

	if total == 0 {
		return []response.ProductResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := applyFilters(base()).Select(`
		c.id AS company_id,
		c.name AS company_name,
		c.base_currency AS company_currency,
		p.id AS id,
		p.name AS name,
		p.status AS status
	`)

	dataQuery = helper.CompanyFilter(dataQuery, s.db, user.Role, user)
	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch product: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil

}
