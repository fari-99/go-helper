package models

type PaymentRequery struct {
    // Merchant Code assigned by iPay88
    MerchantCode string `json:"MerchantCode"`

    // Unique merchant transaction number / Order ID
    RefNo string `json:"RefNo"`

    // The amount must only contain the exact amount without 2 digits after.
    // For example:
    // 100.50 is expressed as 100,50
    // 10 is expressed as 10
    // 0.50 is expressed as 0.50
    // Amount = Sum of Each transaction + Amount Of Shipping + Amount Of Discount
    Amount string `json:"Amount"`
}
