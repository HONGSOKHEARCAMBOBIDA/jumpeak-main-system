package service

import (
	"context"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"

	"gorm.io/gorm"
)

// CustomerLedgerService is intentionally read-only: every ledger row is written as a
// side effect of Invoice/Payment/Refund/DebtAdjustment services, never edited directly.
type CustomerLedgerService interface {
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerLedgerResponse, *model.PaginationMetadata, error)
}

type customerledgerservice struct {
	db *gorm.DB
}

func NewCustomerLedgerService() CustomerLedgerService {
	return &customerledgerservice{
		db: config.DB,
	}
}

func (s *customerledgerservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.CustomerLedgerResponse, *model.PaginationMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	helper.NormalizePagination(&pf)

	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var data []response.CustomerLedgerResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("customer_ledger l").
			Joins("LEFT JOIN customers c ON c.id = l.customer_id").
			Joins("LEFT JOIN companies cp ON cp.id = c.company_id").
			Joins("LEFT JOIN branches b ON b.id = c.branch_id").
			Where("l.company_id = ?", user.CompanyID)
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["customer_id"]; ok && v != "" {
			tx = tx.Where("l.customer_id = ?", v)
		}
		if v, ok := filter["reference_type"]; ok && v != "" {
			tx = tx.Where("l.reference_type = ?", v)
		}
		if v, ok := filter["date_from"]; ok && v != "" {
			tx = tx.Where("l.entry_date >= ?", v)
		}
		if v, ok := filter["date_to"]; ok && v != "" {
			tx = tx.Where("l.entry_date <= ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count customer ledger: %w", err)
	}
	if total == 0 {
		return []response.CustomerLedgerResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		l.id AS id,
		l.customer_id AS customer_id,
		c.name AS customer_name,
		l.entry_date AS entry_date,
		l.reference_type AS reference_type,
		l.reference_id AS reference_id,
		l.description AS description,
		l.debit AS debit,
		l.credit AS credit,
		l.running_balance AS running_balance,
		cp.name AS company_Name,
		b.name AS branch_Name,
		b.code AS branch_Code
	`)

	dataQuery = helper.ApplyAccessFilter(dataQuery, s.db, user.Role, user)
	// Chronological order (oldest first) so the running balance reads top-to-bottom
	// like a statement of account; reverse in the UI if a "latest first" view is wanted.
	if err := dataQuery.Order("l.entry_date ASC").Order("l.id ASC").Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch customer ledger: %w", err)
	}

	for i := range data {
		data[i].EntryDate = helper.FormatDate(data[i].EntryDate)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
