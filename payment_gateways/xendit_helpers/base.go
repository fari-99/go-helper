package xendit_helpers

import (
	"fmt"
	"os"

	"github.com/spf13/cast"
	"github.com/xendit/xendit-go"

	"github.com/fari-99/go-helper/payment_gateways/models"
	xenditConstant "github.com/fari-99/go-helper/payment_gateways/xendit_helpers/constants"
)

type BaseXenditHelpers struct {
	XenditSecretKey         string
	XenditVerificationToken string
	Currency                string
	ReminderTimeUnit        string
	ReminderTime            int
	TransactionUuid         string

	TransactionModel        *models.Transactions
	TransactionAddressModel *models.TransactionAddress
	TransactionItemModels   []models.TransactionItems
	TransactionUser         *models.TransactionUsers
	PaymentMethods          []models.PaymentMethods
}

type XenditInvoiceData struct {
	descriptions string

	notifications  *xendit.InvoiceCustomerNotificationPreference
	additionalFee  []xendit.InvoiceFee
	user           *xendit.InvoiceCustomer
	address        []xendit.InvoiceCustomerAddress
	invoiceItems   []xendit.InvoiceItem
	paymentMethods []string

	totalAdditionalFee float64
	totalItemFee       float64
}

func NewXenditHelpers(transactionUuid string) *BaseXenditHelpers {
	currency := os.Getenv("CURRENCY_DEFAULT")
	if currency == "" {
		currency = "IDR"
	}

	if os.Getenv("XENDIT_REMINDER_UNIT") == "" || os.Getenv("XENDIT_REMINDER_TIME") == "" {
		panic("reminder unit and or reminder time is empty")
	}

	if os.Getenv("XENDIT_SECRET_KEY") == "" {
		panic("xendit secret key is empty")
	}

	if os.Getenv("XENDIT_VERIFICATION_TOKEN") == "" {
		panic("xendit verification token is empty")
	}

	base := &BaseXenditHelpers{
		TransactionUuid:         transactionUuid,
		Currency:                currency,
		ReminderTimeUnit:        os.Getenv("XENDIT_REMINDER_UNIT"),
		ReminderTime:            cast.ToInt(os.Getenv("XENDIT_REMINDER_TIME")),
		XenditSecretKey:         os.Getenv("XENDIT_SECRET_KEY"),
		XenditVerificationToken: os.Getenv("XENDIT_VERIFICATION_TOKEN"),
	}

	return base
}

func (base *BaseXenditHelpers) SetTransactionDetails(transactionModel models.Transactions) *BaseXenditHelpers {
	base.TransactionModel = &transactionModel
	return base
}

func (base *BaseXenditHelpers) SetTransactionAddress(transactionAddressModel models.TransactionAddress) *BaseXenditHelpers {
	base.TransactionAddressModel = &transactionAddressModel
	return base
}

func (base *BaseXenditHelpers) SetTransactionItems(transactionItemModels []models.TransactionItems) *BaseXenditHelpers {
	base.TransactionItemModels = transactionItemModels
	return base
}

func (base *BaseXenditHelpers) SetTransactionUser(transactionUser models.TransactionUsers) *BaseXenditHelpers {
	base.TransactionUser = &transactionUser
	return base
}

func (base *BaseXenditHelpers) SetPaymentMethods(paymentMethods []models.PaymentMethods) *BaseXenditHelpers {
	base.PaymentMethods = paymentMethods
	return base
}

func (base *BaseXenditHelpers) CheckCallbackToken(callbackToken string) error {
	if callbackToken != base.XenditVerificationToken {
		return fmt.Errorf("verification token is invalid")
	}

	return nil
}

func (base *BaseXenditHelpers) generateXenditData(module int) (*XenditInvoiceData, error) {
	paymentMethods, err := generateXenditPaymentMethods(base.PaymentMethods, module)
	if err != nil {
		return nil, err
	}

	address, err := generateXenditAddress(base.TransactionAddressModel)
	if err != nil {
		return nil, err
	}

	user, err := generateXenditUser(base.TransactionUser)
	if err != nil {
		return nil, err
	}

	user.Address = address

	notifications, err := getXenditNotifications()
	if err != nil {
		return nil, err
	}

	invoiceItems, totalItem, err := generateXenditItems(base.TransactionItemModels)
	if err != nil {
		return nil, err
	}

	additionalFee, totalAdditionalFee, err := getAdditionalFee()
	if err != nil {
		return nil, err
	}

	xenditData := XenditInvoiceData{
		descriptions: base.TransactionModel.Descriptions,

		notifications:      &notifications,
		additionalFee:      additionalFee,
		user:               user,
		address:            address,
		invoiceItems:       invoiceItems,
		totalAdditionalFee: totalAdditionalFee,
		totalItemFee:       totalItem,
		paymentMethods:     paymentMethods,
	}

	return &xenditData, nil
}

func getXenditNotifications() (xendit.InvoiceCustomerNotificationPreference, error) {
	notifications := xendit.InvoiceCustomerNotificationPreference{
		InvoiceCreated:  []string{"whatsapp", "sms", "email"},
		InvoiceReminder: []string{"whatsapp", "sms", "email"},
		InvoicePaid:     []string{"whatsapp", "sms", "email"},
		InvoiceExpired:  []string{"whatsapp", "sms", "email"},
	}

	return notifications, nil
}

// TODO: generate additional fee
func getAdditionalFee() ([]xendit.InvoiceFee, float64, error) {
	additionalFee := []xendit.InvoiceFee{
		{
			Type:  "ADMIN",
			Value: 5000,
		},
	}

	return additionalFee, 5000, nil
}

func generateXenditPaymentMethods(paymentMethodModels []models.PaymentMethods, module int) ([]string, error) {
	if paymentMethodModels == nil || len(paymentMethodModels) == 0 { // default payment method is Virtual Accounts
		paymentTypeDetails, _ := xenditConstant.GetPaymentTypeDetail(xenditConstant.PaymentTypeVirtualAccount)

		var paymentMethods []string
		for _, paymentMethod := range paymentTypeDetails.PaymentMethods {
			if value, ok := paymentMethod[xenditConstant.Indonesia].Code[module]; ok {
				paymentMethods = append(paymentMethods, value)
			}
		}

		return paymentMethods, nil
	}

	var paymentMethods []string
	for _, paymentMethodModel := range paymentMethodModels {
		paymentMethods = append(paymentMethods, paymentMethodModel.Code)
	}

	return paymentMethods, nil
}

func generateXenditUser(transactionUser *models.TransactionUsers) (*xendit.InvoiceCustomer, error) {
	if transactionUser == nil {
		return nil, fmt.Errorf("transaction user is empty")
	}

	user := xendit.InvoiceCustomer{
		GivenNames:   fmt.Sprintf("%s %s", transactionUser.FirstName, transactionUser.LastName),
		Email:        transactionUser.EmailAddress,
		MobileNumber: transactionUser.Phone,
	}

	return &user, nil
}

func generateXenditAddress(transactionAddress *models.TransactionAddress) ([]xendit.InvoiceCustomerAddress, error) {
	if transactionAddress == nil {
		return nil, fmt.Errorf("transaction address is empty")
	}

	address := []xendit.InvoiceCustomerAddress{
		{
			City:        transactionAddress.CityName,
			Country:     transactionAddress.CountryName,
			PostalCode:  transactionAddress.Postcode,
			State:       transactionAddress.ProvinceName,
			StreetLine1: transactionAddress.Address,
			StreetLine2: transactionAddress.AdditionalAddress,
		},
	}

	return address, nil
}

func generateXenditItems(transactionItems []models.TransactionItems) ([]xendit.InvoiceItem, float64, error) {
	var total float64
	var invoiceItems []xendit.InvoiceItem
	for _, transactionItem := range transactionItems {
		itemName := transactionItem.ProductName

		invoiceItem := xendit.InvoiceItem{
			Name:     itemName,
			Price:    transactionItem.TotalPrice,
			Quantity: int(transactionItem.Qty),
			Category: transactionItem.ProductCategoryName,
			Url:      transactionItem.ItemUrl,
		}

		invoiceItems = append(invoiceItems, invoiceItem)
		total += transactionItem.TotalPrice
	}

	return invoiceItems, total, nil
}
