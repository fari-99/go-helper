package models

type PaymentMethods struct {
	PaymentMethodTypeID int    `json:"payment_method_type_id"` // e-wallet, virtual-account, etc
	PaymentMethodID     int    `json:"payment_method_id"`      // ovo, dana, bni-va, bca-va, etc
	Code                string `json:"code"`                   // "OVO", "DANA", etc
}
