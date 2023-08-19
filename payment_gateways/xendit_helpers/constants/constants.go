package constants

import (
	"fmt"
	"os"

	"github.com/go-playground/locales/currency"
)

const (
	Indonesia = iota + 1
	Philippines
	Vietnam
	Thailand
	Malaysia
)

type Country int

func GetCurrencyCode() map[currency.Type]string {
	currencyData := map[currency.Type]string{
		currency.IDR: "IDR",
		currency.PHP: "PHP",
		currency.THB: "THB",
		currency.MYR: "MYR",
		currency.VND: "VND",
	}

	return currencyData
}

const (
	UrlSuccessRedirectUrl = iota
	UrlFailedRedirectUrl
)

func GetAllUrlRedirect() map[int]string {
	url := map[int]string{
		UrlSuccessRedirectUrl: os.Getenv("DOMAIN_URL_CUSTOMER") + "/payments/success",
		UrlFailedRedirectUrl:  os.Getenv("DOMAIN_URL_CUSTOMER") + "/payments/failed",
	}

	return url
}

func GetUrlRedirect(urlRedirectType int) (*string, error) {
	allUrl := GetAllUrlRedirect()
	if value, ok := allUrl[urlRedirectType]; ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("url type [%d] not found", urlRedirectType)
	}
}

const (
	InvoicePaid    = "PAID"
	InvoiceExpired = "EXPIRED"
)
