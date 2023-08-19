package flip_helpers

import (
	"encoding/json"
	"fmt"

	"github.com/fari-99/go-flip"
	flipConstants "github.com/fari-99/go-flip/constants"
	flipModel "github.com/fari-99/go-flip/models"

	"github.com/fari-99/go-helper/payment_gateways/models"
)

type AcceptPayments interface {
	CreateBill() (*models.Invoices, error)
	GetBill() (*flipModel.GetBillingResponse, error)
	UpdateBill() (*flipModel.EditBillingResponse, error)
	ConfirmBill() (*flipModel.ConfirmBillPaymentResponse, error)
}

type acceptPayments struct {
	baseHelpers *FlipHelpers
}

func NewAcceptPayments(flipHelpers *FlipHelpers) AcceptPayments {
	return acceptPayments{flipHelpers}
}

func (repo acceptPayments) CreateBill() (*models.Invoices, error) {
	createBillParams := flipModel.CreateBillRequest{
		Title: repo.baseHelpers.TransactionModel.ReferenceNo,
		Type:  flipConstants.BillTypeSingle,
		// Amount:                base.TransactionModel, // TODO: calculate total amount
		ExpiredDate:           repo.baseHelpers.TransactionModel.ExpiredAt.Format(flipConstants.TimeFormatExpiredDate),
		RedirectUrl:           repo.baseHelpers.TransactionModel.RedirectUrl,
		IsAddressRequired:     flipConstants.SelfieFlagTrue,
		IsPhoneNumberRequired: flipConstants.SelfieFlagTrue,
		Step:                  flipConstants.BillStepTwo,
		SenderName:            fmt.Sprintf("%s %s", repo.baseHelpers.TransactionUser.FirstName, repo.baseHelpers.TransactionUser.LastName),
		SenderEmail:           repo.baseHelpers.TransactionUser.EmailAddress,
		SenderPhoneNumber:     repo.baseHelpers.TransactionUser.Phone,
		SenderAddress:         repo.baseHelpers.TransactionUser.Address,
		// SenderBank:            base.TransactionModel.PaymentMethodType, // TODO: get payment gateway type
		SenderBankType: repo.baseHelpers.TransactionModel.PaymentMethodCode,
	}

	flipHelper := flip.NewBaseFlip()
	bill, idemKey, err := flipHelper.CreateBill(createBillParams)
	if err != nil {
		return nil, err
	}

	billMarshal, _ := json.Marshal(bill)

	invoice := models.Invoices{
		TransactionUuid: repo.baseHelpers.TransactionModel.TransactionUuid,
		// TotalPrice:        0, // TODO: calculate total amount
		PaymentGatewayID:  repo.baseHelpers.TransactionModel.PaymentGatewayID,
		PaymentMethodType: repo.baseHelpers.TransactionModel.PaymentMethodType,
		PaymentMethodCode: repo.baseHelpers.TransactionModel.PaymentMethodCode,
		ExpiredAt:         *bill.ExpiredDate,
		Identifier:        idemKey,
		RedirectUrl:       bill.RedirectUrl,
		ResponseJson:      string(billMarshal),
	}

	return &invoice, nil
}

func (repo acceptPayments) GetBill() (*flipModel.GetBillingResponse, error) {
	return nil, nil
}

func (repo acceptPayments) UpdateBill() (*flipModel.EditBillingResponse, error) {
	return nil, nil
}

func (repo acceptPayments) ConfirmBill() (*flipModel.ConfirmBillPaymentResponse, error) {
	return nil, nil
}
