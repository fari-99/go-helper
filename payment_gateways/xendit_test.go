package payment_gateways

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/fari-99/go-helper/payment_gateways/models"
	xenditConstant "github.com/fari-99/go-helper/payment_gateways/xendit_helpers/constants"
)

type XenditTestData struct {
	PaymentTypeID int
	Model         int
	Country       int
}

func GetTestXenditData(input XenditTestData) (models.Transactions, error) {
	paymentTypes, err := xenditConstant.GetPaymentTypeDetail(xenditConstant.PaymentTypes(input.PaymentTypeID))
	if err != nil {
		return models.Transactions{}, err
	}

	var paymentMethods []models.PaymentMethods
	for paymentTypeID, paymentType := range paymentTypes.PaymentMethods {
		for country, paymentMethod := range paymentType {
			if int(country) == input.Country {
				paymentMethods = append(paymentMethods, models.PaymentMethods{
					PaymentMethodTypeID: input.PaymentTypeID,
					PaymentMethodID:     int(paymentTypeID),
					Code:                paymentMethod.Code[input.Model],
				})
			}
		}
	}

	items := []models.TransactionItems{
		{
			TransactionUuid:     "uuid-123456879456",
			TransactionItemUuid: "item-uuid-123456789",
			Qty:                 1,        // Quantity between 1 and 10
			TotalPrice:          12345600, // Random price up to 100
			ExpiredAt:           GetRandomFutureTime(),
			ProductName:         "Product-123456",
			ProductCategoryName: "Category-123456",
			ItemUrl:             fmt.Sprintf("https://example.com/item/%v", rand.Intn(1000)),
		},
		{
			TransactionUuid:     "uuid-23456789",
			TransactionItemUuid: "item-uuid-23456789",
			Qty:                 10,        // Quantity between 1 and 10
			TotalPrice:          456789100, // Random price up to 100
			ExpiredAt:           GetRandomFutureTime(),
			ProductName:         "Product-456789",
			ProductCategoryName: "Category-456789",
			ItemUrl:             fmt.Sprintf("https://example.com/item/%v", rand.Intn(1000)),
		},
	}

	shippingAddress := models.TransactionAddress{
		TransactionUuid:   "uuid-12345",
		FirstName:         "John",
		LastName:          "Doe",
		Address:           "123 Main St",
		AdditionalAddress: "Apt 4B",
		Postcode:          "12345",
		EmailAddress:      "johndoe@example.com",
		Phone:             "+1-123-456-7890",
		CountryName:       "United States",
		CountryCode:       "US",
		ProvinceName:      "California",
		ProvinceCode:      "CA",
		CityName:          "Los Angeles",
		CityCode:          "LA",
	}

	billingAddress := shippingAddress

	users := models.TransactionUsers{
		TransactionUuid:   "uuid-12345",
		FirstName:         "John",
		LastName:          "Doe",
		Address:           "123 Main St",
		AdditionalAddress: "Apt 4B",
		Postcode:          "12345",
		EmailAddress:      "johndoe@example.com",
		Phone:             "+1-123-456-7890",
		CountryName:       "United States",
		CountryCode:       "US",
		ProvinceName:      "California",
		ProvinceCode:      "CA",
		CityName:          "Los Angeles",
		CityCode:          "LA",
	}

	company := []models.TransactionCompanies{
		{
			TransactionUuid: "uuid-67890",
			CompanyID:       1234567890,
			CompanyCode:     "C123",
			Name:            "Tech Corp.",
			Address:         "456 Tech St",
			Postcode:        "67890",
			Email:           "contact@techcorp.com",
			MobilePhone:     "+1-234-567-8901",
			CountryName:     "Canada",
			CountryCode:     "CA",
			ProvinceName:    "British Columbia",
			ProvinceCode:    "BC",
			CityName:        "Vancouver",
			CityCode:        "VAN",
		},
	}

	transactionModel := models.Transactions{
		TransactionUuid:  "uuid-9876543251",
		PaymentGatewayID: XenditID,
		// PaymentMethodType: int8(rand.Intn(5) + 1), // Let's assume 5 different types
		// PaymentMethodCode: fmt.Sprintf("code-%v", rand.Intn(100)),
		ReferenceNo:  "ref-123456789",
		Descriptions: "Description for transaction 123465789",
		RedirectUrl:  fmt.Sprintf("https://example.com/redirect/%v", rand.Intn(100)),
		ExpiredAt:    GetRandomFutureTime(),

		TransactionItems:           items,
		TransactionBillingAddress:  &billingAddress,
		TransactionShippingAddress: &shippingAddress,
		TransactionUsers:           &users,
		TransactionCompanies:       company,
		PaymentMethods:             paymentMethods,
	}

	return transactionModel, nil
}

func TestXenditCreateInvoice(t *testing.T) {
	os.Setenv("XENDIT_TEST", "true")
	os.Setenv("XENDIT_VERIFICATION_TOKEN", "")
	os.Setenv("XENDIT_SECRET_KEY", "")
	os.Setenv("XENDIT_REMINDER_UNIT", "hours")
	os.Setenv("XENDIT_REMINDER_TIME", "1")
	os.Setenv("DOMAIN_URL_CUSTOMER", "http://example.com/v1")

	testData := XenditTestData{
		PaymentTypeID: xenditConstant.PaymentTypeVirtualAccount,
		Model:         xenditConstant.ModuleInvoices,
		Country:       xenditConstant.Indonesia,
	}

	transactionModel, err := GetTestXenditData(testData)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}

	invoices, err := CreateInvoice(transactionModel)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}

	invoiceMarshal, _ := json.MarshalIndent(invoices, "", " ")
	log.Printf(string(invoiceMarshal))
	return
}
