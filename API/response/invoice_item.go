package response

import "mysql/model/base"

type InvoiceItemResponse struct {
	base.ModelBase
	InvoiceID      uint64  `gorm:"not null;index" json:"invoice_id"`
	ProductID      *uint64 `gorm:"index" json:"product_id,omitempty"`
	ProductName    string  `gorm:"index" json:"product_Name,omitempty"`
	Description    string  `gorm:"type:varchar(255);not null" json:"description"`
	Quantity       float64 `gorm:"type:decimal(12,2);not null" json:"quantity"`
	UnitPrice      float64 `gorm:"type:decimal(18,2);not null" json:"unit_price"`
	DiscountAmount float64 `gorm:"type:decimal(18,2);not null;default:0.00" json:"discount_amount"`
	Subtotal       float64 `gorm:"type:decimal(18,2);not null" json:"subtotal"`
	CurrencyCode   string  `gorm:"type:char(3);not null" json:"currency_code"`
}
