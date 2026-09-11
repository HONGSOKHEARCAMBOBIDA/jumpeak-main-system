package service

import (
	"context"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type ReportService interface {
	GetCustomerOutstandingReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerOutstandingReport, *model.PaginationMetadata, error)
}

type reportservice struct {
	db *gorm.DB
}

func NewReportService() ReportService {
	return &reportservice{
		db: config.DB,
	}
}

func (s *reportservice) GetCustomerOutstandingReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerOutstandingReport, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.CustomerOutstandingReport

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
		if v, ok := filter["branch_id"]; ok && v != "" {
			tx = tx.Where("b.id = ?", v)
		}
		if v, ok := filter["company_id"]; ok && v != "" {
			tx = tx.Where("cp.id = ?", v)
		}
		return tx
	}

	var total int64

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count companies: %w", err)
	}

	if total == 0 {
		return []response.CustomerOutstandingReport{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize

	dataQuery := helper.ApplyAccessFilter(applyFilters(base()), s.db, user.Role, user).
		Select(`
			cp.id AS company_id,
			cp.name AS company_name,
			cp.base_currency AS currency,
			COALESCE(SUM(c.current_outstanding), 0) AS total_amount
		`).
		Group("cp.id, cp.name")

	if err := dataQuery.Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch companies: %w", err)
	}

	if len(data) == 0 {
		return data, helper.BuildPaginationMeta(pf, total), nil
	}

	for i := range data {
		data[i].Currency = helper.Currency(data[i].Currency)
	}

	companyIDs := make([]uint64, 0, len(data))
	for _, c := range data {
		companyIDs = append(companyIDs, uint64(c.CompanyID))
	}

	var branchOutStandingReport []response.BranchOutStandingReport
	branchQuery := s.db.WithContext(ctx).Table("customers c").
		Joins("LEFT JOIN branches b ON b.id = c.branch_id").
		Joins("LEFT JOIN companies cp ON cp.id = c.company_id").
		Where("cp.id IN (?)", companyIDs).
		Select(`
			cp.id AS company_id,
			cp.base_currency AS currency,
			b.id AS branch_id,
			b.name AS Branch_name,
			COALESCE(SUM(c.current_outstanding), 0) AS total_amount
		`).
		Group("cp.id, b.id, b.name")

	branchQuery = helper.ApplyAccessFilter(branchQuery, s.db, user.Role, user)
	if err := branchQuery.Scan(&branchOutStandingReport).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch branches: %w", err)
	}

	for i := range branchOutStandingReport {
		branchOutStandingReport[i].Currency = helper.Currency(branchOutStandingReport[i].Currency)
	}

	branchesByCompany := make(map[uint64][]response.BranchOutStandingReport, len(data))
	for _, b := range branchOutStandingReport {
		branchesByCompany[uint64(b.CompanyID)] = append(branchesByCompany[uint64(b.CompanyID)], b)
	}

	for i := range data {
		data[i].BranchOutStandingReport = branchesByCompany[uint64(data[i].CompanyID)]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
