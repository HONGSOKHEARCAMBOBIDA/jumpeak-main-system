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

type DebtAdjustmentService interface {
	Create(ctx context.Context, userID int, input request.DebtAdjustmentRequestCreate) error
	Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.DebtAdjustmentResponse, *model.PaginationMetadata, error)
}

type debtadjustmentservice struct {
	db *gorm.DB
}

func NewDebtAdjustmentService() DebtAdjustmentService {
	return &debtadjustmentservice{
		db: config.DB,
	}
}

// Create records a WRITE_OFF, CORRECTION or DISCOUNT against a customer (optionally
// tied to one invoice) and reduces what the customer owes accordingly.
func (s *debtadjustmentservice) Create(ctx context.Context, userID int, input request.DebtAdjustmentRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	if input.Amount <= 0 {
		return apperror.New(apperror.CodeInvalidInput, "adjustment amount must be greater than zero", nil)
	}

	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	adjType := model.DebtAdjustmentType(input.Type)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.InvoiceID != nil {
			var invoice model.Invoice
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", *input.InvoiceID).First(&invoice).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperror.New(apperror.CodeNotFound, "invoice not found", nil)
				}
				return apperror.New(apperror.CodeInternal, "failed to fetch invoice", nil)
			}
			if invoice.CustomerID != input.CustomerID {
				return apperror.New(apperror.CodeInvalidInput, "invoice does not belong to this customer", nil)
			}
			if input.Amount > invoice.OutstandingAmount {
				return apperror.New(apperror.CodeInvalidInput, "adjustment exceeds invoice outstanding balance", nil)
			}

			invoice.OutstandingAmount -= input.Amount
			switch {
			case invoice.OutstandingAmount > 0:
				invoice.Status = model.InvoiceStatusPartiallyPaid
			case adjType == model.DebtAdjustmentTypeWriteOff:
				invoice.OutstandingAmount = 0
				invoice.Status = model.InvoiceStatusWrittenOff
			default:
				invoice.OutstandingAmount = 0
				invoice.Status = model.InvoiceStatusPaid
			}
			if err := tx.Save(&invoice).Error; err != nil {
				return apperror.New(apperror.CodeInternal, "failed to update invoice", nil)
			}
		}

		newdata := model.DebtAdjustment{
			CompanyID:  user.CompanyID,
			CustomerID: input.CustomerID,
			InvoiceID:  input.InvoiceID,
			Type:       adjType,
			Amount:     input.Amount,
			Reason:     input.Reason,
			ApprovedBy: uint64(userID),
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return helper.MapError(err, "CREATE")
		}

		if err := helper.AppendLedgerEntry(
			tx, user.CompanyID, input.CustomerID, time.Now(),
			model.CustomerLedgerReferenceAdjustment, uint64(newdata.ID),
			fmt.Sprintf("%s: %s", adjType, input.Reason),
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

func (s *debtadjustmentservice) Get(ctx context.Context, userID int, pf request.Pagination, filter map[string]string) ([]response.DebtAdjustmentResponse, *model.PaginationMetadata, error) {
	helper.NormalizePagination(&pf)

	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var data []response.DebtAdjustmentResponse
	var total int64

	base := func() *gorm.DB {
		return s.db.WithContext(ctx).
			Table("debt_adjustments d").
			Where("d.company_id = ?", user.CompanyID)
	}

	applyFilters := func(tx *gorm.DB) *gorm.DB {
		if v, ok := filter["customer_id"]; ok && v != "" {
			tx = tx.Where("d.customer_id = ?", v)
		}
		if v, ok := filter["type"]; ok && v != "" {
			tx = tx.Where("d.type = ?", v)
		}
		return tx
	}

	if err := applyFilters(base()).Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count debt adjustment: %w", err)
	}
	if total == 0 {
		return []response.DebtAdjustmentResponse{}, helper.BuildPaginationMeta(pf, total), nil
	}

	offset := (pf.Page - 1) * pf.PageSize
	dataQuery := applyFilters(base()).Select(`
		d.id AS id,
		d.customer_id AS customer_id,
		d.invoice_id AS invoice_id,
		d.type AS type,
		d.amount AS amount,
		d.reason AS reason,
		d.approved_by AS approved_by,
		d.created_at AS created_at
	`)

	if err := dataQuery.Order("d.id DESC").Offset(offset).Limit(pf.PageSize).Scan(&data).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch debt adjustments: %w", err)
	}

	return data, helper.BuildPaginationMeta(pf, total), nil
}
