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

type PaymentService interface {
	Create(ctx context.Context, userID int, input request.PaymentRequestCreate) error
	Void(ctx context.Context, id int, userID int, input request.PaymentRequestVoid) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.PaymentResponse, *model.PaginationMetadata, error)
}

type paymentservice struct {
	db *gorm.DB
}

func NewPaymentService() PaymentService {
	return &paymentservice{
		db: config.DB,
	}
}

func (s *paymentservice) Create(ctx context.Context, userID int, input request.PaymentRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if input.Amount <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "payment amount must be greater than zero", nil)
	}
	var allocatedTotal float64
	for _, a := range input.Allocations {
		allocatedTotal += a.Amount
	}
	if allocatedTotal > input.Amount {
		return apperror.New(apperror.CodeInvalidInput, "allocations exceed payment amount", nil)
	}

	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		paymentdate, err := time.Parse("2006-01-02", input.PaymentDate)
		if err != nil {
			return err
		}
		newdata := model.Payment{
			CompanyID:          user.CompanyID,
			BranchID:           user.BranchID,
			CustomerID:         input.CustomerID,
			PaymentDate:        paymentdate,
			CurrencyCode:       input.CurrencyCode,
			ExchangeRateToBase: input.ExchangeRateToBase,
			Amount:             input.Amount,
			Method:             model.PaymentMethod(input.Method),
			ReferenceNumber:    input.ReferenceNumber,
			Note:               input.Note,
			Status:             model.PaymentStatusCompleted,
			CreatedBy:          uint64Ptr(uint64(userID)),
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}
		newdata.PaymentNumber = helper.GenerateCode("PAY-", uint(newdata.ID))
		if err := tx.Model(&newdata).Update("payment_number", newdata.PaymentNumber).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to assign payment number", nil)
		}

		// Apply each allocation against its invoice, locking the invoice row so two
		// payments can't both allocate against the same stale outstanding balance.
		for _, alloc := range input.Allocations {
			if alloc.Amount <= 0 {
				return apperror.New(apperror.CodeInvalidInput, "allocation amount must be greater than zero", nil)
			}
			var invoice model.Invoice
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", alloc.InvoiceID).First(&invoice).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperror.New(apperror.CodeNotFound, "invoice not found", nil)
				}
				return apperror.New(apperror.CodeInternal, "failed to fetch invoice", nil)
			}
			if invoice.CustomerID != input.CustomerID {
				return apperror.New(apperror.CodeInvalidInput, "invoice does not belong to this customer", nil)
			}
			if alloc.Amount > invoice.OutstandingAmount {
				return apperror.New(apperror.CodeInvalidInput, fmt.Sprintf("allocation exceeds outstanding balance of invoice %s", invoice.InvoiceNumber), nil)
			}

			allocation := model.PaymentAllocation{
				PaymentID: newdata.ID,
				InvoiceID: uint64(invoice.ID),
				Amount:    alloc.Amount,
			}
			if err := tx.Create(&allocation).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create payment allocation", nil)
			}

			invoice.PaidAmount += alloc.Amount
			invoice.OutstandingAmount -= alloc.Amount
			if invoice.OutstandingAmount <= 0 {
				invoice.OutstandingAmount = 0
				invoice.Status = model.InvoiceStatusPaid
				now := time.Now()
				invoice.PaidAt = &now
			} else {
				invoice.Status = model.InvoiceStatusPartiallyPaid
			}
			if err := tx.Save(&invoice).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to update invoice", nil)
			}
		}

		if err := helper.AppendLedgerEntry(
			tx, user.CompanyID, input.CustomerID, paymentdate,
			model.CustomerLedgerReferencePayment, newdata.ID,
			fmt.Sprintf("Payment %s", newdata.PaymentNumber),
			0, input.Amount,
		); err != nil {
			return err
		}

		if err := tx.Model(&model.Customer{}).
			Where("id = ?", input.CustomerID).
			Update("current_outstanding", gorm.Expr("current_outstanding - ?", input.Amount)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update customer balance", nil)
		}

		return nil
	})
	return err
}

// Void reverses a payment: every allocation is unwound so the affected invoices
// re-open, and the customer's balance/ledger are restored to their pre-payment state.
func (s *paymentservice) Void(ctx context.Context, id int, userID int, input request.PaymentRequestVoid) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment model.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).First(&payment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "payment not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch payment", nil)
		}
		if payment.Status == model.PaymentStatusVoided {
			return apperror.New(apperror.CodeInvalidInput, "payment is already voided", nil)
		}

		var allocations []model.PaymentAllocation
		if err := tx.Where("payment_id = ?", payment.ID).Find(&allocations).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to fetch payment allocations", nil)
		}

		for _, alloc := range allocations {
			var invoice model.Invoice
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", alloc.InvoiceID).First(&invoice).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to fetch invoice", nil)
			}
			invoice.PaidAmount -= alloc.Amount
			invoice.OutstandingAmount += alloc.Amount
			invoice.PaidAt = nil
			if invoice.PaidAmount <= 0 {
				invoice.PaidAmount = 0
				invoice.Status = model.InvoiceStatusOpen
			} else {
				invoice.Status = model.InvoiceStatusPartiallyPaid
			}
			if err := tx.Save(&invoice).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to update invoice", nil)
			}
		}

		payment.Status = model.PaymentStatusVoided
		payment.Note = &input.Reason
		if err := tx.Save(&payment).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to void payment", nil)
		}

		if err := helper.AppendLedgerEntry(
			tx, payment.CompanyID, payment.CustomerID, time.Now(),
			model.CustomerLedgerReferencePayment, payment.ID,
			fmt.Sprintf("Voided payment %s: %s", payment.PaymentNumber, input.Reason),
			payment.Amount, 0,
		); err != nil {
			return err
		}

		if err := tx.Model(&model.Customer{}).
			Where("id = ?", payment.CustomerID).
			Update("current_outstanding", gorm.Expr("current_outstanding + ?", payment.Amount)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update customer balance", nil)
		}

		return nil
	})
	return err
}

func (s *paymentservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.PaymentResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var data []response.PaymentResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("payments p").
			Joins("LEFT JOIN customers c ON c.id = p.customer_id").
			Where("p.company_id = ?", user.CompanyID)
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["customer_id"]; ok && v != "" {
			tx = tx.Where("p.customer_id = ?", v)
		}
		if v, ok := filter["status"]; ok && v != "" {
			tx = tx.Where("p.status = ?", v)
		}
		if v, ok := filter["method"]; ok && v != "" {
			tx = tx.Where("p.method = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count payment: %w", err)
	}
	if total == 0 {
		return []response.PaymentResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		p.id AS id,
		p.company_id AS company_id,
		p.customer_id AS customer_id,
		c.name AS customer_name,
		p.payment_number AS payment_number,
		p.payment_date AS payment_date,
		p.currency_code AS currency_code,
		p.exchange_rate_to_base AS exchange_rate_to_base,
		p.amount AS amount,
		p.method AS method,
		p.reference_number AS reference_number,
		p.note AS note,
		p.status AS status
	`)

	if err := dataQuery.Order("p.id DESC").Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch payments: %w", err)
	}

	for i := range data {
		data[i].PaymentDate = helper.FormatDate(data[i].PaymentDate)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
