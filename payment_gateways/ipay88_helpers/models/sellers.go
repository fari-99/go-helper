package models

type Sellers struct {
    // Seller’s ID
    ID string `json:"Id"`

    // Seller’s Name
    Name string `json:"Name"`

    LegalID string `json:"LegalId"`

    // Seller's identifier number (KTP / SIM / etc)
    SellerIDNumber *string `json:"SellerIdNumber,omitempty"`

    // Seller’s Email
    Email string `json:"Email"`

    // Seller’s Website
    URL string `json:"Url"`

    // Seller’s Shop Address
    Address TransactionAddress `json:"address"`
}
