package helper

import "mysql/model"

func CustomerLedger(Type model.CustomerLedgerReferenceType) string {
	switch Type {
	case model.CustomerLedgerReferenceInvoice:
		return "វិក្កយបត្រ"

	case model.CustomerLedgerReferencePayment:
		return "ការសងប្រាក់"

	case model.CustomerLedgerReferenceRefund:
		return "បង់ប្រាក់ទៅអតិថិជនវិញ"

	case model.CustomerLedgerReferenceAdjustment:
		return "កែបំណុល"

	default:
		return ""
	}
}
