package xendit_helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"

	"github.com/fari-99/go-helper/payment_gateways/models"
	"github.com/fari-99/go-helper/payment_gateways/xendit_helpers/constants"
)

type Invoices interface {
	GetAllInvoices() ([]xendit.Invoice, *xendit.Error)
	GetInvoiceByID(xenditInvoiceID string) (*xendit.Invoice, *xendit.Error)
	CreateInvoice() (*models.Invoices, error)
	CancelInvoice(xenditInvoiceID string) (*xendit.Invoice, *xendit.Error)
}

type invoices struct {
	base *BaseXenditHelpers
}

func NewInvoices(base *BaseXenditHelpers) Invoices {
	return invoices{base: base}
}

func (repo invoices) GetAllInvoices() ([]xendit.Invoice, *xendit.Error) {
	params := invoice.GetAllParams{
		ForUserID:          "",
		Statuses:           nil,
		Limit:              0,
		CreatedAfter:       time.Time{},
		CreatedBefore:      time.Time{},
		PaidAfter:          time.Time{},
		PaidBefore:         time.Time{},
		ExpiredAfter:       time.Time{},
		ExpiredBefore:      time.Time{},
		LastInvoiceID:      "",
		ClientTypes:        nil,
		PaymentChannels:    nil,
		OnDemandLink:       "",
		RecurringPaymentID: "",
	}

	xendit.Opt.SecretKey = repo.base.XenditSecretKey
	invoiceList, errXendit := invoice.GetAll(&params)
	return invoiceList, errXendit
}

func (repo invoices) GetInvoiceByID(xenditInvoiceID string) (*xendit.Invoice, *xendit.Error) {
	params := invoice.GetParams{
		ID: xenditInvoiceID,
	}

	xendit.Opt.SecretKey = repo.base.XenditSecretKey
	invoiceData, errXendit := invoice.Get(&params)
	return invoiceData, errXendit
}

func (repo invoices) CreateInvoice() (*models.Invoices, error) {
	transactionUuid := repo.base.TransactionUuid
	transactionDetails := repo.base.TransactionModel

	xenditInvoiceData, err := repo.base.generateXenditData(constants.ModuleInvoices)
	if err != nil {
		return nil, err
	}

	totalAmount := xenditInvoiceData.totalItemFee + xenditInvoiceData.totalAdditionalFee
	urlSuccessRedirect, _ := constants.GetUrlRedirect(constants.UrlSuccessRedirectUrl)
	urlFailedRedirect, _ := constants.GetUrlRedirect(constants.UrlFailedRedirectUrl)

	shouldSendEmail := true
	invoiceParams := invoice.CreateParams{
		ExternalID:                     transactionUuid,
		Amount:                         totalAmount,
		Description:                    xenditInvoiceData.descriptions,
		PayerEmail:                     xenditInvoiceData.user.Email,
		ShouldSendEmail:                &shouldSendEmail,
		Customer:                       *xenditInvoiceData.user,
		CustomerNotificationPreference: *xenditInvoiceData.notifications,
		InvoiceDuration:                int(transactionDetails.ExpiredAt.Sub(time.Now()).Seconds()),
		SuccessRedirectURL:             *urlSuccessRedirect,
		FailureRedirectURL:             *urlFailedRedirect,
		PaymentMethods:                 xenditInvoiceData.paymentMethods,
		Currency:                       repo.base.Currency,
		ReminderTimeUnit:               repo.base.ReminderTimeUnit,
		ReminderTime:                   repo.base.ReminderTime,
		Locale:                         "id", // default ID
		Items:                          xenditInvoiceData.invoiceItems,
		Fees:                           xenditInvoiceData.additionalFee,
		// MidLabel:                       "test-mid", // if using credit cards
	}

	xendit.Opt.SecretKey = repo.base.XenditSecretKey
	invoiceResp, errXendit := invoice.Create(&invoiceParams)
	if errXendit != nil {
		return nil, fmt.Errorf(errXendit.Error())
	}

	invoiceRespMarshal, _ := json.Marshal(invoiceResp)
	invoiceModel := models.Invoices{
		PaymentGatewayID:  transactionDetails.PaymentGatewayID,
		PaymentMethodType: transactionDetails.PaymentMethodType,
		PaymentMethodCode: transactionDetails.PaymentMethodCode,
		TransactionUuid:   transactionUuid,
		TotalPrice:        invoiceResp.Amount,
		Identifier:        invoiceResp.ID,
		RedirectUrl:       invoiceResp.InvoiceURL,
		ResponseJson:      string(invoiceRespMarshal),
		ExpiredAt:         invoiceResp.ExpiryDate.String(),
		// Status:            config.GetConstStatus("STATUS_ACTIVE", nil), // TODO: adding status
	}

	return &invoiceModel, nil
}

func (repo invoices) CancelInvoice(xenditInvoiceID string) (*xendit.Invoice, *xendit.Error) {
	params := invoice.ExpireParams{
		ID: xenditInvoiceID,
	}

	xendit.Opt.SecretKey = repo.base.XenditSecretKey
	invoiceData, errXendit := invoice.Expire(&params)
	return invoiceData, errXendit
}
