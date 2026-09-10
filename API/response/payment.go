package response

type PaymentResponse struct {
	ID                 uint64  `json:"id"`
	CompanyID          uint64  `json:"company_id"`
	CustomerID         uint64  `json:"customer_id"`
	CustomerName       string  `json:"customer_name"`
	PaymentNumber      string  `json:"payment_number"`
	PaymentDate        string  `json:"payment_date"`
	CurrencyCode       string  `json:"currency_code"`
	ExchangeRateToBase float64 `json:"exchange_rate_to_base"`
	Amount             float64 `json:"amount"`
	Method             string  `json:"method"`
	ReferenceNumber    *string `json:"reference_number,omitempty"`
	Note               string  `json:"note"`
	Status             string  `json:"status"`
}
