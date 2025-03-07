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
	GetBill(billID int64) (*flipModel.EditBillingResponse, error) // TODO: change on go-flip get bill response
	UpdateBill(billID int64, isActive bool) (*flipModel.EditBillingResponse, error)
	ConfirmCallback(token string) (bool, error)
}

type acceptPayments struct {
	flipData *FlipData
}

func NewAcceptPayments(flipData *FlipData) AcceptPayments {
	return acceptPayments{flipData}
}

func (repo acceptPayments) CreateBill() (*models.Invoices, error) {
	transactionModel := repo.flipData.TransactionModel
	transactionUser := repo.flipData.TransactionUser

	createBillParams := flipModel.CreateBillRequest{
		Title:                 repo.flipData.TransactionModel.ReferenceNo,
		Type:                  flipConstants.BillTypeSingle,
		Amount:                fmt.Sprintf("%.2f", repo.flipData.TotalItemFee),
		ExpiredDate:           transactionModel.ExpiredAt.Format(flipConstants.TimeFormatExpiredDate),
		RedirectUrl:           transactionModel.RedirectUrl,
		IsAddressRequired:     flipConstants.SelfieFlagTrue,
		IsPhoneNumberRequired: flipConstants.SelfieFlagTrue,
		Step:                  flipConstants.BillStepTwo, // get flip app payments redirect url
		SenderName:            fmt.Sprintf("%s %s", transactionUser.FirstName, transactionUser.LastName),
		SenderEmail:           transactionUser.EmailAddress,
		SenderPhoneNumber:     transactionUser.Phone,
		SenderAddress:         transactionUser.Address,
	}

	flipBase := flip.NewBaseFlip()
	bill, idemKey, err := flipBase.CreateBill(createBillParams)
	if err != nil {
		return nil, err
	}

	billMarshal, _ := json.Marshal(bill)

	invoice := models.Invoices{
		TransactionUuid:   transactionModel.TransactionUuid,
		TotalPrice:        float64(bill.Amount),
		PaymentGatewayID:  transactionModel.PaymentGatewayID,
		PaymentMethodType: transactionModel.PaymentMethodType,
		PaymentMethodCode: transactionModel.PaymentMethodCode,
		ExpiredAt:         *bill.ExpiredDate,
		Identifier:        idemKey,
		RedirectUrl:       bill.RedirectUrl,
		ResponseJson:      string(billMarshal),
	}

	return &invoice, nil
}

func (repo acceptPayments) GetBill(billID int64) (*flipModel.EditBillingResponse, error) {
	flipBase := flip.NewBaseFlip()
	bill, err := flipBase.GetBill(billID)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (repo acceptPayments) UpdateBill(billID int64, isActive bool) (*flipModel.EditBillingResponse, error) {
	transactionModel := repo.flipData.TransactionModel

	status := flipConstants.BillStatusActive
	if !isActive {
		status = flipConstants.BillStatusInActive
	}

	updateData := flipModel.EditBillingRequest{
		Title:                 repo.flipData.TransactionModel.ReferenceNo,
		Type:                  flipConstants.BillTypeSingle,
		Amount:                fmt.Sprintf("%.2f", repo.flipData.TotalItemFee),
		ExpiredDate:           transactionModel.ExpiredAt.Format(flipConstants.TimeFormatExpiredDate),
		RedirectUrl:           transactionModel.RedirectUrl,
		IsAddressRequired:     flipConstants.SelfieFlagTrue,
		IsPhoneNumberRequired: flipConstants.SelfieFlagTrue,
		Status:                status,
	}

	flipBase := flip.NewBaseFlip()
	bill, err := flipBase.EditBill(billID, updateData)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (repo acceptPayments) ConfirmCallback(token string) (bool, error) {
	isValid, err := flip.CheckCallback(token)
	return isValid, err
}
