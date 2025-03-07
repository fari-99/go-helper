package models

type Invoices struct {
	TransactionUuid   string  `json:"transaction_uuid"`
	InvoiceNo         string  `json:"invoice_no"`
	TotalPrice        float64 `json:"total_price"`
	PaymentGatewayID  int8    `json:"payment_gateway_id"`
	PaymentMethodType int8    `json:"payment_gateway_type"`
	PaymentMethodCode string  `json:"payment_method_id"`

	ExpiredAt      string `json:"expired_at"`
	Identifier     string `json:"identifier"`
	RedirectUrl    string `json:"redirect_url"`
	RedirectParams string `json:"redirect_params"`
	ResponseJson   string `json:"response_json"`
}
