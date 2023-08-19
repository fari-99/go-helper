package models

import "time"

type Transactions struct {
	TransactionUuid   string     `json:"transaction_uuid"`
	PaymentGatewayID  int8       `json:"payment_gateway_id"`
	PaymentMethodType int8       `json:"payment_method_type"`
	PaymentMethodCode string     `json:"payment_method_code"`
	ReferenceNo       string     `json:"title"`
	Descriptions      string     `json:"descriptions"`
	RedirectUrl       string     `json:"redirect_url"`
	ExpiredAt         *time.Time `json:"expired_at"`

	TransactionItems           []TransactionItems     `json:"transaction_items"`
	TransactionShippingAddress *TransactionAddress    `json:"transaction_shipping_address"`
	TransactionBillingAddress  *TransactionAddress    `json:"transaction_billing_address"`
	TransactionUsers           *TransactionUsers      `json:"transaction_users"`
	TransactionCompanies       []TransactionCompanies `json:"transaction_companies"`
	PaymentMethods             []PaymentMethods       `json:"payment_methods"`
}

func (model Transactions) TableName() string {
	return "transactions"
}
