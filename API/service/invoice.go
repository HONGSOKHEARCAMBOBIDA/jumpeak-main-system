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
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvoiceService interface {
	Create(ctx context.Context, userID int, input request.InvoiceRequestCreate) error
	Cancel(ctx context.Context, id int, userID int, input request.InvoiceRequestCancel) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.InvoiceResponse, *model.PaginationMetadata, error)
}

type invoiceservice struct {
	db *gorm.DB
}

func NewInvoiceService() InvoiceService {
	return &invoiceservice{
		db: config.DB,
	}
}

func (s *invoiceservice) Create(ctx context.Context, userID int, input request.InvoiceRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if len(input.Items) == 0 {
		return apperror.New(apperror.CodeInvalidInput, "invoice must have at least one item", nil)
	}

	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var total float64
	items := make([]model.InvoiceItem, 0, len(input.Items))
	for _, it := range input.Items {
		subtotal := (it.Quantity * it.UnitPrice) - it.DiscountAmount
		total += subtotal
		items = append(items, model.InvoiceItem{
			ProductID:      it.ProductID,
			Description:    it.Description,
			Quantity:       it.Quantity,
			UnitPrice:      it.UnitPrice,
			DiscountAmount: it.DiscountAmount,
			Subtotal:       subtotal,
		})
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var customer model.Customer
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", input.CustomerID).
			First(&customer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "customer not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch customer", nil)
		}
		if customer.Status == model.CustomerStatusInactive {
			return apperror.New(apperror.CodeInvalidInput, "customer is not active", nil)
		}
		if customer.Status == model.CustomerStatusBlacklisted {
			return apperror.New(apperror.CodeForbidden, "customer in black list", nil)
		}
		if customer.CreditLimitEnforced && customer.CurrentOutstanding+total > customer.CreditLimit {
			return apperror.New(apperror.CodeInvalidInput, "invoice exceeds customer credit limit", nil)
		}

		invoiceDate, err := time.Parse("2006-01-02", input.InvoiceDate)
		if err != nil {
			return err
		}

		dueDate, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return err
		}

		newdata := model.Invoice{
			CompanyID:          user.CompanyID,
			BranchID:           *user.BranchID,
			CustomerID:         input.CustomerID,
			InvoiceDate:        invoiceDate,
			DueDate:            dueDate,
			CurrencyCode:       input.CurrencyCode,
			ExchangeRateToBase: input.ExchangeRateToBase,
			TotalAmount:        total,
			PaidAmount:         0,
			OutstandingAmount:  total,
			Status:             model.InvoiceStatusOpen,
			CreatedBy:          uintPtr(userID),
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}

		newdata.InvoiceNumber = helper.GenerateCode("INV-", uint(newdata.ID))
		if err := tx.Model(&newdata).Update("invoice_number", newdata.InvoiceNumber).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to assign invoice number", nil)
		}

		for i := range items {
			items[i].InvoiceID = uint64(newdata.ID)
		}
		if err := tx.Create(&items).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to create invoice items", nil)
		}

		if err := helper.AppendLedgerEntry(
			tx, user.CompanyID, input.CustomerID, invoiceDate,
			model.CustomerLedgerReferenceInvoice, uint64(newdata.ID),
			fmt.Sprintf("វិក្កយបត្រ %s", newdata.InvoiceNumber),
			total, 0,
		); err != nil {
			return err
		}

		if err := tx.Model(&model.Customer{}).
			Where("id = ?", input.CustomerID).
			Update("current_outstanding", gorm.Expr("current_outstanding + ?", total)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update customer balance", nil)
		}

		return nil
	})
	return err
}

// Cancel voids an invoice that has not received any payment yet and reverses its
// effect on the customer's outstanding balance and ledger.
func (s *invoiceservice) Cancel(ctx context.Context, id int, userID int, input request.InvoiceRequestCancel) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var invoice model.Invoice
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).First(&invoice).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "invoice not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch invoice", nil)
		}
		if invoice.Status == model.InvoiceStatusCancelled {
			return apperror.New(apperror.CodeInvalidInput, "invoice is already cancelled", nil)
		}
		if invoice.PaidAmount > 0 {
			return apperror.New(apperror.CodeInvalidInput, "cannot cancel an invoice that has received payment; issue a refund or adjustment instead", nil)
		}

		outstanding := invoice.OutstandingAmount
		now := time.Now()
		invoice.Status = model.InvoiceStatusCancelled
		invoice.CancelReason = &input.Reason
		invoice.CancelledAt = &now
		invoice.CancelledBy = uint64Ptr(uint64(userID))
		invoice.OutstandingAmount = 0
		invoice.UpdatedBy = uint64Ptr(uint64(userID))
		if err := tx.Save(&invoice).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to cancel invoice", nil)
		}

		if err := helper.AppendLedgerEntry(
			tx, invoice.CompanyID, invoice.CustomerID, now,
			model.CustomerLedgerReferenceInvoice, uint64(invoice.ID),
			fmt.Sprintf("លុបវិក្កយបត្រ %s", invoice.InvoiceNumber),
			0, outstanding,
		); err != nil {
			return err
		}

		if err := tx.Model(&model.Customer{}).
			Where("id = ?", invoice.CustomerID).
			Update("current_outstanding", gorm.Expr("current_outstanding - ?", outstanding)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update customer balance", nil)
		}

		return nil
	})
	return err
}

func (s *invoiceservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.InvoiceResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var data []response.InvoiceResponse
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
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count invoice: %w", err)
	}

	if total == 0 {
		return []response.InvoiceResponse{}, helper.BuildPaginationMeta(pf, total), nil
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

func uintPtr(v int) *uint64 {
	u := uint64(v)
	return &u
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}
