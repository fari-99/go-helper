package models

type ItemTransactions struct {
    // Product id
    // ex: 001
    ID string `json:"Id"`

    // Name of item
    // ex: Smartphone
    Name string `json:"Name"`

    // Quantity of item
    // ex: 2
    Quantity string `json:"Quantity"`

    // Amount of item product
    // ex: 5000
    Amount string `json:"Amount"`

    // Type of Product
    // ex:
    Type *string `json:"Type,omitempty"`

    // URL of the product in the merchant site / platform
    // ex:
    URL *string `json:"Url,omitempty"`

    // URL of the image of the product in the merchant site / platform
    // ex:
    ImageURL *string `json:"ImageUrl,omitempty"`

    // For BCA KlikPay
    Tenor *string `json:"Tenor,omitempty"`

    // For BCA KlikPay
    CodePlan *string `json:"CodePlan,omitempty"`

    // For BCA KlikPay
    MerchantId *string `json:"MerchantId,omitempty"`

    // Possible values : SELLER, ITEM
    ParentType *string `json:"ParentType,omitempty"`

    // It will correspond to the SELLER ID if the parentType is SELLER,
    // alternatively it will correspond to the ITEM ID if the parentType is ITEM
    ParentID *string `json:"ParentId,omitempty"`
}
