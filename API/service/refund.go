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
	"gorm.io/gorm/clause"
)

type RefundService interface {
	Create(ctx context.Context, userID int, input request.RefundRequestCreate) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.RefundResponse, *model.PaginationMetadata, error)
}

type refundservice struct {
	db *gorm.DB
}

func NewRefundService() RefundService {
	return &refundservice{
		db: config.DB,
	}
}

// Create refunds money that was previously received via a Payment. Every
// RefundAllocation re-opens the corresponding amount on its invoice (mirror image
// of what the original PaymentAllocation did) and the customer owes that much again.
func (s *refundservice) Create(ctx context.Context, userID int, input request.RefundRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if input.Amount <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "refund amount must be greater than zero", nil)
	}
	var allocatedTotal float64
	for _, a := range input.Allocations {
		allocatedTotal += a.Amount
	}
	if allocatedTotal != input.Amount {
		return apperror.New(apperror.CodeInvalidInput, "allocations must add up to the refund amount", nil)
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment model.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", input.PaymentID).First(&payment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.New(apperror.CodeNotFound, "payment not found", nil)
			}
			return apperror.New(apperror.CodeInternal, "failed to fetch payment", nil)
		}
		if payment.Status != model.PaymentStatusCompleted {
			return apperror.New(apperror.CodeInvalidInput, "only completed payments can be refunded", nil)
		}

		newdata := model.Refund{
			PaymentID:  payment.ID,
			Amount:     input.Amount,
			Reason:     input.Reason,
			RefundedAt: input.RefundedAt,
			CreatedBy:  uint64Ptr(uint64(userID)),
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}

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
			if invoice.CustomerID != payment.CustomerID {
				return apperror.New(apperror.CodeInvalidInput, "invoice does not belong to the payment's customer", nil)
			}
			if alloc.Amount > invoice.PaidAmount {
				return apperror.New(apperror.CodeInvalidInput, fmt.Sprintf("refund exceeds paid amount on invoice %s", invoice.InvoiceNumber), nil)
			}

			refundAlloc := model.RefundAllocation{
				RefundID:  uint64(newdata.ID),
				InvoiceID: uint64(invoice.ID),
				Amount:    alloc.Amount,
			}
			if err := tx.Create(&refundAlloc).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to create refund allocation", nil)
			}

			invoice.PaidAmount -= alloc.Amount
			invoice.OutstandingAmount += alloc.Amount
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

		if err := helper.AppendLedgerEntry(
			tx, payment.CompanyID, payment.CustomerID, input.RefundedAt,
			model.CustomerLedgerReferenceRefund, uint64(newdata.ID),
			fmt.Sprintf("Refund against payment %s: %s", payment.PaymentNumber, input.Reason),
			input.Amount, 0,
		); err != nil {
			return err
		}

		if err := tx.Model(&model.Customer{}).
			Where("id = ?", payment.CustomerID).
			Update("current_outstanding", gorm.Expr("current_outstanding + ?", input.Amount)).Error; err != nil {
			return apperror.New(apperror.CodeInternal, "failed to update customer balance", nil)
		}

		return nil
	})
	return err
}

func (s *refundservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.RefundResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var data []response.RefundResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("refunds r").
			Joins("LEFT JOIN payments p ON p.id = r.payment_id")
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["payment_id"]; ok && v != "" {
			tx = tx.Where("r.payment_id = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count refund: %w", err)
	}
	if total == 0 {
		return []response.RefundResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		r.id AS id,
		r.payment_id AS payment_id,
		r.amount AS amount,
		r.reason AS reason,
		r.refunded_at AS refunded_at
	`)

	if err := dataQuery.Order("r.id DESC").Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch refunds: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
