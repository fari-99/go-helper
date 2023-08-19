package payment_gateways

import (
	"fmt"

	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers"
	ipay88Constant "github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/constants"
	"github.com/fari-99/go-helper/payment_gateways/models"
	"github.com/fari-99/go-helper/payment_gateways/xendit_helpers"
	xenditConstant "github.com/fari-99/go-helper/payment_gateways/xendit_helpers/constants"
)

const (
	XenditID = iota + 1
	Ipay88ID
	FlipID
)

func GetPaymentGateways() map[int]string {
	return map[int]string{
		XenditID: "Xendit",
		Ipay88ID: "Ipay88",
		FlipID:   "Flip",
	}
}

func GetPaymentGatewayMethod(paymentGatewayID int) (interface{}, error) {
	switch paymentGatewayID {
	case XenditID:
		return xenditConstant.GetPaymentTypes(), nil
	case Ipay88ID:
		return ipay88Constant.GetPaymentTypes(), nil
	case FlipID:
		message := "flip didn't have options to select payment method by merchant"
		return map[string]interface{}{
			"message": message,
		}, nil
	default:
		return nil, fmt.Errorf("payment gateway id [%d] is not found", paymentGatewayID)
	}
}

func CreateInvoice(transactionModel models.Transactions) (*models.Invoices, error) {
	switch transactionModel.PaymentGatewayID {
	case XenditID:
		xenditHelpers := xendit_helpers.NewXenditHelpers(transactionModel.TransactionUuid)
		xenditHelpers.SetTransactionDetails(transactionModel)
		xenditHelpers.SetTransactionAddress(*transactionModel.TransactionBillingAddress)
		xenditHelpers.SetTransactionUser(*transactionModel.TransactionUsers)
		xenditHelpers.SetTransactionItems(transactionModel.TransactionItems)
		xenditHelpers.SetPaymentMethods(transactionModel.PaymentMethods)

		xenditInvoice := xendit_helpers.NewInvoices(xenditHelpers)
		invoiceModel, err := xenditInvoice.CreateInvoice()
		if err != nil {
			return nil, err
		}

		return invoiceModel, nil
	case Ipay88ID:
		ipay88Helper := ipay88_helpers.NewIpay88Helper()
		ipay88Helper.SetTransactionModel(transactionModel)
		ipay88Helper.SetTransactionUser(*transactionModel.TransactionUsers)
		ipay88Helper.SetTransactionItems(transactionModel.TransactionItems)
		ipay88Helper.SetBillingAddress(*transactionModel.TransactionBillingAddress)
		ipay88Helper.SetShippingAddress(*transactionModel.TransactionShippingAddress)
		ipay88Helper.SetTransactionCompanies(transactionModel.TransactionCompanies)

		invoiceModel, err := ipay88Helper.CreatPaymentRequest()
		if err != nil {
			return nil, err
		}

		return invoiceModel, nil
	default:
		return nil, fmt.Errorf("payment gateway not found")
	}
}
