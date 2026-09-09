package helper

import (
	"errors"
	"mysql/constant/apperror"
	"mysql/model"
	"time"

	"gorm.io/gorm"
)

// lastLedgerBalance returns the customer's most recent running balance, or 0 if the
// customer has no ledger entries yet. Must be called with the *transaction* handle
// so it reads a consistent view of the row it's about to extend.
func lastLedgerBalance(tx *gorm.DB, customerID uint64) (float64, error) {
	var last model.CustomerLedger
	err := tx.Where("customer_id = ?", customerID).
		Order("id DESC").
		Limit(1).
		First(&last).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return last.RunningBalance, nil
}

// appendLedgerEntry writes one CustomerLedger row chained off the customer's last
// running balance. Debit increases what the customer owes (invoices, refunds paid
// out, reversed payments). Credit decreases it (payments received, write-offs,
// discounts, cancelled invoices).
func appendLedgerEntry(
	tx *gorm.DB,
	companyID, customerID uint64,
	entryDate time.Time,
	refType model.CustomerLedgerReferenceType,
	refID uint64,
	description string,
	debit, credit float64,
) error {
	lastBalance, err := lastLedgerBalance(tx, customerID)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "failed to read customer ledger", nil)
	}

	entry := model.CustomerLedger{
		CompanyID:      companyID,
		CustomerID:     customerID,
		EntryDate:      entryDate,
		ReferenceType:  refType,
		ReferenceID:    refID,
		Description:    description,
		Debit:          debit,
		Credit:         credit,
		RunningBalance: lastBalance + debit - credit,
	}
	if err := tx.Create(&entry).Error; err != nil {
		return apperror.New(apperror.CodeInternal, "failed to write customer ledger", nil)
	}
	return nil
}
