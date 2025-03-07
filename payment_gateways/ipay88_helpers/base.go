package ipay88_helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	ipay88Model "github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/models"
	"github.com/fari-99/go-helper/payment_gateways/models"
)

type BaseIpay88Helper struct {
	MerchantCode string
	MerchantKey  string

	TransactionModel           *models.Transactions
	TransactionUser            *models.TransactionUsers
	TransactionItems           []models.TransactionItems
	TransactionShippingAddress *models.TransactionAddress
	TransactionBillingAddress  *models.TransactionAddress
	TransactionCompanies       map[uint64]models.TransactionCompanies

	Invoice models.Invoices
}

func NewIpay88Helper() *BaseIpay88Helper {
	base := BaseIpay88Helper{
		MerchantCode: os.Getenv("IPAY88_MERCHANT_CODE"),
		MerchantKey:  os.Getenv("IPAY88_MERCHANT_KEY"),
	}

	return &base
}

func (base *BaseIpay88Helper) SetTransactionModel(data models.Transactions) *BaseIpay88Helper {
	base.TransactionModel = &data
	return base
}

func (base *BaseIpay88Helper) SetTransactionUser(data models.TransactionUsers) *BaseIpay88Helper {
	base.TransactionUser = &data
	return base
}

func (base *BaseIpay88Helper) SetTransactionItems(data []models.TransactionItems) *BaseIpay88Helper {
	base.TransactionItems = data
	return base
}

func (base *BaseIpay88Helper) SetBillingAddress(data models.TransactionAddress) *BaseIpay88Helper {
	base.TransactionBillingAddress = &data
	return base
}

func (base *BaseIpay88Helper) SetShippingAddress(data models.TransactionAddress) *BaseIpay88Helper {
	base.TransactionShippingAddress = &data
	return base
}

func (base *BaseIpay88Helper) SetTransactionCompanies(data []models.TransactionCompanies) *BaseIpay88Helper {
	transactionCompanies := make(map[uint64]models.TransactionCompanies)
	for _, transactionCompany := range data {
		if _, ok := transactionCompanies[transactionCompany.CompanyID]; !ok {
			transactionCompanies[transactionCompany.CompanyID] = transactionCompany
		}
	}

	base.TransactionCompanies = transactionCompanies
	return base
}

func (base *BaseIpay88Helper) generateSignature(refNo string, amount uint64, currency string) (string, error) {
	merchantKey := base.MerchantKey
	merchantCode := base.MerchantCode

	if merchantKey == "" || merchantCode == "" {
		return "", fmt.Errorf("merchant key or merchant code is empty")
	}

	hashData := fmt.Sprintf("||%s||%s||%s||%d||%s||", merchantKey, merchantCode, refNo, amount, currency)

	hs := sha256.New()
	hs.Write([]byte(hashData))
	return hex.EncodeToString(hs.Sum(nil)), nil
}

func (base *BaseIpay88Helper) generatePaymentRequestData() (*ipay88Model.PaymentRequestData, error) {
	itemTransactions, totalItemPrice, err := generateItemTransaction(base.TransactionItems, base.TransactionCompanies)
	if err != nil {
		return nil, err
	}

	shippingAddress, err := generateAddress(base.TransactionShippingAddress)
	if err != nil {
		return nil, err
	}

	billingAddress, err := generateAddress(base.TransactionBillingAddress)
	if err != nil {
		return nil, err
	}

	sellers, err := generateSellers(base.TransactionCompanies)
	if err != nil {
		return nil, err
	}

	additionalFee, totalAdditionalFee, err := getAdditionalFee()
	if err != nil {
		return nil, err
	}

	input := ipay88Model.PaymentRequestData{
		Items:           itemTransactions,
		ShippingAddress: *shippingAddress,
		BillingAddress:  *billingAddress,
		Sellers:         sellers,
		AdditionalFee:   additionalFee,

		TotalItemPrice:     totalItemPrice,
		TotalAdditionalFee: totalAdditionalFee,
	}

	return &input, nil
}

func getAdditionalFee() ([]ipay88Model.ItemTransactions, float64, error) {
	itemType := "ADDITIONAL_FEE"
	additionalFee := []ipay88Model.ItemTransactions{
		{
			ID:       uuid.New().String(),
			Name:     "Admin Fee",
			Quantity: "1",
			Amount:   fmt.Sprintf("%d", 5000),
			Type:     &itemType,
			// URL:        nil,
			// ImageURL:   nil,
			// Tenor:      nil,
			// CodePlan:   nil,
			// MerchantId: nil,
			// ParentType: nil,
			// ParentID:   nil,
		},
	}

	return additionalFee, 5000, nil
}

func generateItemTransaction(transactionItems []models.TransactionItems, transactionCompanies map[uint64]models.TransactionCompanies) ([]ipay88Model.ItemTransactions, float64, error) {
	if transactionItems == nil || len(transactionItems) == 0 {
		return nil, 0, fmt.Errorf("transaction items is empty")
	}

	var itemTransactions []ipay88Model.ItemTransactions
	var total float64
	for _, transactionItem := range transactionItems {
		itemTransaction := ipay88Model.ItemTransactions{
			ID:       transactionItem.TransactionItemUuid,
			Name:     transactionItem.ProductName,
			Quantity: fmt.Sprintf("%d", transactionItem.Qty), // seated only have 1 qty
			Amount:   cast.ToString(transactionItem.TotalPrice),
			Type:     &transactionItem.ProductCategoryName,
			URL:      &transactionItem.ItemUrl, // TODO: website item url
			ImageURL: nil,
			// Tenor:      nil,
			// CodePlan:   nil,
			// MerchantId: nil,
			ParentType: nil,
			ParentID:   nil,
		}

		itemTransactions = append(itemTransactions, itemTransaction)
		total += transactionItem.TotalPrice
	}

	return itemTransactions, total, nil
}

func generateAddress(transactionAddress *models.TransactionAddress) (*ipay88Model.TransactionAddress, error) {
	if transactionAddress == nil || transactionAddress.TransactionUuid == "" {
		return nil, fmt.Errorf("transaction address is empty")
	}

	address := ipay88Model.TransactionAddress{
		FirstName:   transactionAddress.FirstName,
		LastName:    transactionAddress.LastName,
		Address:     transactionAddress.Address,
		City:        transactionAddress.CityName,
		State:       transactionAddress.ProvinceName,
		PostalCode:  transactionAddress.Postcode,
		Phone:       transactionAddress.Phone,
		CountryCode: "ID", // default country code "ID"
	}

	return &address, nil
}

func generateSellers(companies map[uint64]models.TransactionCompanies) ([]ipay88Model.Sellers, error) {
	if companies == nil || len(companies) == 0 {
		return nil, fmt.Errorf("transaction address is empty")
	}

	var sellers []ipay88Model.Sellers
	for _, company := range companies {
		seller := ipay88Model.Sellers{
			ID:   cast.ToString(company.CompanyID),
			Name: company.Name,
			// LegalID:        "",
			// SellerIDNumber: nil,
			Email: company.Email,
			// URL:            "",
			Address: ipay88Model.TransactionAddress{
				// FirstName:   "",
				// LastName:    "",
				Address:     company.Address,
				City:        company.CityName,
				State:       company.ProvinceName,
				PostalCode:  company.Postcode,
				Phone:       company.MobilePhone,
				CountryCode: "ID",
			},
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}
