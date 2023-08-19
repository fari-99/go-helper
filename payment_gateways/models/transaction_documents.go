package models

type TransactionDocuments struct {
	TransactionUuid string `json:"transaction_uuid"`
}

func (model TransactionDocuments) TableName() string {
	return "transaction_documents"
}
