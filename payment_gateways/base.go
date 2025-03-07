package payment_gateways

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cast"

	"github.com/fari-99/go-helper/payment_gateways/flip_helpers"
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

func GetRandomFutureTime() *time.Time {
	hoursAhead := rand.Intn(48) + 1 // Between 1 and 48 hours
	futureTime := time.Now().Add(time.Duration(hoursAhead) * time.Hour)
	return &futureTime
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
	case FlipID:
		flipHelper := flip_helpers.NewFlipHelpers(transactionModel.TransactionUuid).
			SetTransactionDetails(transactionModel).
			SetTransactionUser(*transactionModel.TransactionUsers).
			SetTransactionItems(transactionModel.TransactionItems)
		flipData, err := flipHelper.GenerateFlipData()
		if err != nil {
			return nil, err
		}

		flipAcceptPayment := flip_helpers.NewAcceptPayments(flipData)
		invoiceModel, err := flipAcceptPayment.CreateBill()
		if err != nil {
			return nil, err
		}

		return invoiceModel, nil
	default:
		return nil, fmt.Errorf("payment gateway not found")
	}
}

func GetDetails(paymentGatewayID int, identifier string) (interface{}, error) {
	switch paymentGatewayID {
	case XenditID:
		xenditInvoice := xendit_helpers.NewInvoices(nil)
		invoice, err := xenditInvoice.GetInvoiceByID(identifier)
		return invoice, err
	case Ipay88ID:
		panic("not yet available, ipay88 didn't use id, but query params to check data")
	case FlipID:
		flipAcceptPayment := flip_helpers.NewAcceptPayments(nil)
		bill, err := flipAcceptPayment.GetBill(cast.ToInt64(identifier))
		return bill, err
	default:
		return nil, fmt.Errorf("payment gateway not found")
	}
}
