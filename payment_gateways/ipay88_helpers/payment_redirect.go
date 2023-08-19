package ipay88_helpers

import (
	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/constants"
	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/models"
)

func (base *BaseIpay88Helper) PaymentRedirectUrl(checkoutID, signature string) (*models.PaymentRedirectUrlData, error) {
	url, err := constants.GetIpay88Url(constants.Ipay88PaymentRedirectUrl)
	if err != nil {
		return nil, err
	}

	data := models.PaymentRedirectUrlData{
		Url:        *url,
		CheckoutID: checkoutID,
		Signature:  signature,
	}

	return &data, nil
}
