package flip_helpers

import (
	"os"

	"github.com/fari-99/go-helper/payment_gateways/models"
)

type FlipHelpers struct {
	TransactionUuid         string
	TransactionModel        *models.Transactions
	TransactionAddressModel *models.TransactionAddress
	TransactionItemModels   []models.TransactionItems
	TransactionUser         *models.TransactionUsers
	PaymentMethods          []models.PaymentMethods
}

func NewFlipHelpers(transactionUuid string) *FlipHelpers {
	if os.Getenv("FLIP_ENVIRONMENT") == "" || os.Getenv("FLIP_ENVIRONMENT") != "prod" {
		_ = os.Setenv("FLIP_ENVIRONMENT", "dev")
	}

	if os.Getenv("FLIP_SECRET_TOKEN") == "" {
		panic("flip secret token is empty")
	}

	if os.Getenv("FLIP_VALIDATION_TOKEN") == "" {
		panic("flip validation token is empty")
	}

	base := &FlipHelpers{
		TransactionUuid: transactionUuid,
	}

	return base
}
