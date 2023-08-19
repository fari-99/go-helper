package models

type PaymentRequestResponse struct {
    // Unique merchant transaction number / Order ID
    RefNo string `json:"RefNo"`

    // Signature generated with sha256 (see 4.2)
    Signature string `json:"Signature"`

    // Virtual Account number
    VirtualAccountAssigned string `json:"VirtualAccountAssigned"`

    // Expired date for Virtual Account (DD-MM-YYYY HH:MM)
    TransactionExpiryDate string `json:"TransactionExpiryDate"`

    // iPay88 OPSG Transaction ID
    CheckoutID string `json:"CheckoutID"`

    // Transaction Status
    Code string `json:"Code"`

    // Transaction Message
    Message string `json:"Message"`
}
