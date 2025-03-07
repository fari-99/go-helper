package models

type TransactionAddress struct {
    // Shipping FirstName
    FirstName string `json:"FirstName"`

    // Shipping LastName
    LastName string `json:"LastName"`

    // Shipping Address
    Address string `json:"Address"`

    // Shipping City
    City string `json:"City"`

    // Shipping State
    State string `json:"State"`

    // Shipping Postal Code
    PostalCode string `json:"PostalCode"`

    // Shipping Phone
    Phone string `json:"Phone"`

    // Shipping Country Code
    CountryCode string `json:"CountryCode"`
}
