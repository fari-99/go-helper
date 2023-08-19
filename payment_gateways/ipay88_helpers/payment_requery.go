package ipay88_helpers

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/constants"
	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/models"
)

func (base *BaseIpay88Helper) PaymentRequery() ([]byte, error) {
	url, err := constants.GetIpay88Url(constants.Ipay88PaymentRequeryUrl)
	if err != nil {
		return nil, err
	}

	invoices := base.Invoice
	paymentRequery := models.PaymentRequery{
		MerchantCode: base.MerchantCode,
		RefNo:        invoices.TransactionUuid,
		Amount:       cast.ToString(invoices.TotalPrice),
	}

	var query map[string]string
	queryMarshal, _ := json.Marshal(paymentRequery)
	_ = json.Unmarshal(queryMarshal, &query)

	client := resty.New()
	resp, err := client.R().
		SetQueryParams(query).
		Get(*url)

	log.Printf(string(resp.Body()))
	return resp.Body(), err
}
