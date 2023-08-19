package models

type PaymentRedirectUrlData struct {
    Url        string `json:"url"`
    CheckoutID string `json:"checkout_id"`
    Signature  string `json:"signature"`
}
