package models

import "time"

type TransactionItems struct {
	TransactionUuid     string     `json:"transaction_uuid"`
	TransactionItemUuid string     `json:"transaction_item_uuid"`
	Qty                 uint64     `json:"qty"`
	TotalPrice          float64    `json:"price"`
	ExpiredAt           *time.Time `json:"expired_at"`
	ProductName         string     `json:"product_name"`
	ProductCategoryName string     `json:"product_category_name"`
	ItemUrl             string     `json:"item_url"`
}
