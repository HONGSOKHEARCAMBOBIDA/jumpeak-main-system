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
	GetOverCreditLimitReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.OverCreditLimitReportResponse, *model.PaginationMetadata, error)
	GetOverDueDateReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.OverDueDateReportResponse, *model.PaginationMetadata, error)
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

func (s *reportservice) GetOverCreditLimitReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.OverCreditLimitReportResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)
	var data []response.OverCreditLimitReportResponse
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
		tx = tx.Where("c.current_outstanding > c.credit_limit")
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count customer: %w", err)
	}

	if total == 0 {
		return []response.OverCreditLimitReportResponse{}, helper.BuildPaginationMeta(pf, total), nil
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

	for i := range data {
		data[i].CompanyCurrency = helper.Currency(data[i].CompanyCurrency)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil

}

func (s *reportservice) GetOverDueDateReport(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.OverDueDateReportResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var data []response.OverDueDateReportResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("invoices i").
			Joins("LEFT JOIN customers c ON c.id = i.customer_id").
			Joins("LEFT JOIN companies cp ON cp.id = i.company_id").
			Joins("LEFT JOIN branches b ON b.id = i.branch_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["customer_id"]; ok && v != "" {
			tx = tx.Where("i.customer_id = ?", v)
		}
		if v, ok := filter["status"]; ok && v != "" {
			tx = tx.Where("i.status = ?", v)
		}
		if v, ok := filter["invoice_number"]; ok && v != "" {
			tx = tx.Where("i.invoice_number LIKE ?", "%"+v+"%")
		}
		tx = tx.Where("i.due_date < CURRENT_DATE")
		tx = tx.Where("i.status = ? OR i.status = ?", model.InvoiceStatusOpen, model.InvoiceStatusPartiallyPaid)
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count invoice: %w", err)
	}

	if total == 0 {
		return []response.OverDueDateReportResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		i.id AS id,
		i.company_id AS company_id,
		i.customer_id AS customer_id,
		c.name AS customer_name,
		i.invoice_number AS invoice_number,
		i.invoice_date AS invoice_date,
		i.due_date AS due_date,
		CASE
		WHEN i.due_date < CURRENT_DATE
		THEN DATEDIFF(CURRENT_DATE, i.due_date)
		ELSE 0
		END AS late_count,
		i.currency_code AS currency_code,
		i.exchange_rate_to_base AS exchange_rate_to_base,
		i.total_amount AS total_amount,
		i.paid_amount AS paid_amount,
		i.outstanding_amount AS outstanding_amount,
		i.status AS status,
		i.cancel_reason AS cancel_reason,
		cp.name AS company_Name,
		b.id AS branch_id,
		b.name AS branch_Name,
		b.phone AS branch_phone
	`)

	dataQuery = helper.ApplyAccessFilter(dataQuery, s.db, user.Role, user)

	if err := dataQuery.Order("i.id DESC").Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch invoices: %w", err)
	}

	for i := range data {
		data[i].InvoiceDate = helper.FormatDate(data[i].InvoiceDate)
		data[i].DueDate = helper.FormatDate(data[i].DueDate)
		data[i].CurrencyCode = helper.Currency(data[i].CurrencyCode)
	}

	invoiceIDs := make([]int, len(data))
	for i, a := range data {
		invoiceIDs[i] = int(a.ID)
	}

	var invoiceitems []response.InvoiceItemResponse
	if err := s.db.WithContext(ctx).Table("invoice_items ii").
		Joins("LEFT JOIN products p ON p.id = ii.product_id").
		Joins("LEFT JOIN invoices i ON i.id = ii.invoice_id").
		Where("ii.invoice_id IN ?", invoiceIDs).
		Select(`
		i.currency_code AS currency_code,
		ii.invoice_id AS invoice_id,
		ii.id AS id,
		p.id AS product_id,
		p.name AS product_Name,
		ii.description AS description,
		ii.quantity AS quantity,
		ii.unit_price AS unit_price,
		ii.discount_amount AS discount_amount,
		ii.subtotal AS subtotal
	`).Scan(&invoiceitems).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch invoice items: %w", err)
	}

	for i := range invoiceitems {
		invoiceitems[i].CurrencyCode = helper.Currency(invoiceitems[i].CurrencyCode)
	}

	itembyinvoice := make(map[uint64][]response.InvoiceItemResponse, len(data))
	for _, i := range invoiceitems {
		itembyinvoice[i.InvoiceID] = append(itembyinvoice[i.InvoiceID], i)
	}
	for i := range data {
		data[i].InvoiceItemResponse = itembyinvoice[data[i].ID]
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
