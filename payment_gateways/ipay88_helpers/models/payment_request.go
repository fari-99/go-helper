package models

type PaymentRequests struct {
    // Version API
    // ex: 2.0
    APIVersion string `json:"APIVersion"`

    // Merchant code assigned by iPay88
    // ex: ID00001
    MerchantCode string `json:"MerchantCode"`

    // Payment method identifier, see doc 8.1 for the list.
    // ex: 70
    PaymentID string `json:"PaymentId"`

    // Support to IDR, USD, THB
    // ex: IDR
    Currency string `json:"Currency"`

    // Unique merchant transaction number / Order ID
    // ex: 920201930925AM
    RefNo string `json:"RefNo"`

    // The amount must only contain the exact amount without 2 digits after.
    // For example:
    // 100.50 is expressed as 100,50
    // 10 is expressed as 10
    // 0.50 is expressed as 0.50
    // Amount = Sum of Each transaction + Amount Of Shipping + Amount Of Discount
    Amount string `json:"Amount"`

    // Product description
    // ex: Alat Elektronik
    ProdDesc string `json:"ProdDesc"`

    // Customer name
    // ex: Thoriq
    UserName string `json:"UserName"`

    // Customer email for receiving receipt
    // ex: thoriq@ipay88.co.id
    UserEmail string `json:"UserEmail"`

    // User contact number
    // ex: 08123123123
    UserContact string `json:"UserContact"`

    // Request type identifier, see doc 8.2 for the list.
    // ex: REDIRECT / SEAMLESS
    RequestType *string `json:"RequestType,omitempty"`

    // Merchant remarks
    Remark *string `json:"Remark,omitempty"`

    // Encoding type
    // ISO-8859-1   | English
    // UTF-8        | English
    // GB2312       | Chinese Simplified
    // GD18030      | Chinese Simplified
    // BIG5         | Chinese Traditional
    // ex: ISO-8859-1
    Lang *string `json:"Lang,omitempty"` // ex:

    // Payment response page
    // ex: https://store.co.id/resp.asp
    ResponseURL string `json:"ResponseURL"`

    // Backend response page
    // ex: https://store.co.id/backend.asp
    BackendURL string `json:"BackendURL"`

    // SHA256 signature
    Signature string `json:"Signature"`

    // for BCA KlikPay
    // mandatory for BCA KlikPay
    // ex: 10000000
    FullTransactionAmount *string `json:"FullTransactionAmount,omitempty"`

    // for BCA KlikPay
    // ex: 10000000
    MiscFee *string `json:"MiscFee,omitempty"`

    // Only for installment detail with multiple product
    ItemTransactions []ItemTransactions `json:"ItemTransactions"`

    // Shipping Address Detail
    ShippingAddress TransactionAddress `json:"ShippingAddress"`

    // Billing Address detail
    BillingAddress TransactionAddress `json:"BillingAddress"`

    // Only For Seller Detail
    Sellers []Sellers `json:"Sellers"`

    SettingField []SettingFields `json:"SettingField"`
}
