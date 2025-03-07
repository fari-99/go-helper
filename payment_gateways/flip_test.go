package payment_gateways

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/fari-99/go-helper/payment_gateways/models"
)

func GetTestFlipData() models.Transactions {
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
		PaymentGatewayID: FlipID,
		ReferenceNo:      "ref-123456789",
		Descriptions:     "Description for transaction 123465789",
		RedirectUrl:      fmt.Sprintf("https://example.com/redirect/%v", rand.Intn(100)),
		ExpiredAt:        GetRandomFutureTime(),

		TransactionItems:           items,
		TransactionBillingAddress:  &billingAddress,
		TransactionShippingAddress: &shippingAddress,
		TransactionUsers:           &users,
		TransactionCompanies:       company,
	}

	return transactionModel
}

func TestCreateBill(t *testing.T) {
	os.Setenv("FLIP_ENVIRONMENT", "dev")
	os.Setenv("FLIP_SECRET_TOKEN", "YOUR FLIP SECRET TOKEN")
	os.Setenv("FLIP_VALIDATION_TOKEN", "YOUR FLIP VALIDATION TOKEN")

	flipData := GetTestFlipData()
	invoices, err := CreateInvoice(flipData)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}

	invoiceMarshal, _ := json.MarshalIndent(invoices, "", " ")
	log.Printf(string(invoiceMarshal))
	return
}
