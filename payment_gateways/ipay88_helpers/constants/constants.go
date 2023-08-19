package constants

import (
	"fmt"
	"os"

	"github.com/go-playground/locales/currency"
	"github.com/spf13/cast"
)

const ApiVersions = "2.0"
const BackendPostResponseError = "0"
const BackendPostResponseSuccess = "1"

const (
	Ipay88PaymentSuccess = "1"
	Ipay88PaymentFail    = "2"
	Ipay88PaymentPending = "3"
)

const (
	Ipay88PaymentRequestUrl = iota
	Ipay88PaymentRedirectUrl
	Ipay88PaymentRequeryUrl
)

type BaseUrl map[int]map[bool]string // true: sandbox, false: prod

func GetAllIpay88Url() BaseUrl {
	baseUrl := BaseUrl{
		Ipay88PaymentRequestUrl: map[bool]string{
			true:  "https://sandbox.ipay88.co.id/ePayment/WebService/PaymentAPI/Checkout",
			false: "https://payment.ipay88.co.id/ePayment/WebService/PaymentAPI/Checkout",
		},
		Ipay88PaymentRedirectUrl: map[bool]string{
			true:  "https://sandbox.ipay88.co.id/PG",
			false: "https://payment.ipay88.co.id/PG",
		},
		Ipay88PaymentRequeryUrl: map[bool]string{
			true:  "https://sandbox.ipay88.co.id/epayment/enquiry.asp",
			false: "https://payment.ipay88.co.id/epayment/enquiry.asp",
		},
	}

	return baseUrl
}

func GetIpay88Url(urlType int) (*string, error) {
	isDev := cast.ToBool(os.Getenv("IPAY88_TEST"))

	baseUrl := GetAllIpay88Url()
	if value, ok := baseUrl[urlType]; ok {
		url := value[isDev]
		return &url, nil
	} else {
		return nil, fmt.Errorf("url type [%d] is not found", urlType)
	}
}

const (
	PaymentResponseUrl = iota
	PaymentBackendUrl
)

func GetAllApiUrl() map[int]string {
	allUrl := map[int]string{
		PaymentResponseUrl: os.Getenv("DOMAIN_URL_CUSTOMER") + "/payments/success",   // success payments page
		PaymentBackendUrl:  os.Getenv("DOMAIN_URL_API") + "/payments/ipay88/backend", // callback to api from ipay88
	}

	return allUrl
}

func GetApiUrl(apiUrlType int) (*string, error) {
	baseUrl := GetAllApiUrl()
	if value, ok := baseUrl[apiUrlType]; ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("url type [%d] is not found", apiUrlType)
	}
}

const (
	CurrencyIDR = iota + 1
	CurrencyUSD
	CurrencyTHB
)

type Currency map[int]string

func GetAllCurrency() Currency {
	return Currency{
		CurrencyIDR: "IDR",
		CurrencyUSD: "USD",
		CurrencyTHB: "THB",
	}
}

func GetCurrencyLabel(currency int) (*string, error) {
	allCurrency := GetAllCurrency()
	if value, ok := allCurrency[currency]; ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("currency type [%d] is not found", currency)
	}
}

const (
	RequestTypeRedirect = iota + 1
	RequestTypeSeamless
)

type RequestType map[int]string

func GetAllRequestType() RequestType {
	return RequestType{
		RequestTypeRedirect: "REDIRECT",
		RequestTypeSeamless: "SEAMLESS",
	}
}

func GetRequestTypeLabel(requestType int) (*string, error) {
	allRequestType := GetAllRequestType()
	if value, ok := allRequestType[requestType]; ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("request type [%d] is not found", requestType)
	}
}

const (
	EncodingISO_8859_1 = iota + 1 // English
	EncodingUTF_8                 // Unicode
	EncodingGB2312                // Chinese Simplified
	EncodingGD18030               // Chinese Simplified
	EncodingBIG5                  // Chinese Traditional
)

type EncodingLanguage map[int]string

func GetAllEncoding() EncodingLanguage {
	return EncodingLanguage{
		EncodingISO_8859_1: "ISO-8859-1",
		EncodingUTF_8:      "UTF-8",
		EncodingGB2312:     "GB2312",
		EncodingGD18030:    "GD18030",
		EncodingBIG5:       "BIG5",
	}
}

func GetEncodingLabel(encoding int) (*string, error) {
	allEncoding := GetAllEncoding()
	if value, ok := allEncoding[encoding]; ok {
		return &value, nil
	} else {
		return nil, fmt.Errorf("encoding type [%d] is not found", encoding)
	}
}

func GetCurrencyCode() map[currency.Type]string {
	currencyData := map[currency.Type]string{
		currency.IDR: "IDR",
	}

	return currencyData
}
