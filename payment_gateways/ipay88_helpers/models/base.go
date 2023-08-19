package models

type PaymentRequestData struct {
    Items           []ItemTransactions
    ShippingAddress TransactionAddress
    BillingAddress  TransactionAddress
    Sellers         []Sellers
    settingFields   SettingFields
    AdditionalFee   []ItemTransactions

    TotalItemPrice     float64
    TotalAdditionalFee float64
}
