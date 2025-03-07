package flip_helpers

import (
	"os"

	"github.com/fari-99/go-helper/payment_gateways/models"
)

type FlipHelpers interface {
	SetTransactionDetails(transactionModel models.Transactions) flipHelpers
	SetTransactionItems(transactionItemModels []models.TransactionItems) flipHelpers
	SetTransactionUser(transactionUser models.TransactionUsers) flipHelpers
	GenerateFlipData() (*FlipData, error)
}

type flipHelpers struct {
	transactionUuid         string
	transactionModel        *models.Transactions
	transactionAddressModel *models.TransactionAddress
	transactionItemModels   []models.TransactionItems
	transactionUser         *models.TransactionUsers
	paymentMethods          []models.PaymentMethods
}

func NewFlipHelpers(transactionUuid string) FlipHelpers {
	if os.Getenv("FLIP_ENVIRONMENT") == "" || os.Getenv("FLIP_ENVIRONMENT") != "prod" {
		_ = os.Setenv("FLIP_ENVIRONMENT", "dev")
	}

	if os.Getenv("FLIP_SECRET_TOKEN") == "" {
		panic("flip secret token is empty")
	}

	if os.Getenv("FLIP_VALIDATION_TOKEN") == "" {
		panic("flip validation token is empty")
	}

	return flipHelpers{
		transactionUuid: transactionUuid,
	}
}

func (base flipHelpers) SetTransactionDetails(transactionModel models.Transactions) flipHelpers {
	base.transactionModel = &transactionModel
	return base
}

func (base flipHelpers) SetTransactionItems(transactionItemModels []models.TransactionItems) flipHelpers {
	base.transactionItemModels = transactionItemModels
	return base
}

func (base flipHelpers) SetTransactionUser(transactionUser models.TransactionUsers) flipHelpers {
	base.transactionUser = &transactionUser
	return base
}

type FlipData struct {
	TransactionModel models.Transactions
	TransactionUser  models.TransactionUsers
	TotalItemFee     float64
}

func (base flipHelpers) GenerateFlipData() (*FlipData, error) {
	return &FlipData{
		TransactionModel: *base.transactionModel,
		TransactionUser:  *base.transactionUser,
		TotalItemFee:     calculateItemFee(base.transactionItemModels),
	}, nil
}

func calculateItemFee(items []models.TransactionItems) float64 {
	var amount float64
	for _, item := range items {
		amount += item.TotalPrice
	}

	return amount
}
