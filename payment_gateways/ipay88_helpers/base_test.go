package ipay88_helpers

import (
	"os"
	"testing"

	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/constants"
)

func TestGenerateSignature(t *testing.T) {
	_ = os.Setenv("IPAY88_MERCHANT_KEY", "Apple")
	_ = os.Setenv("IPAY88_MERCHANT_CODE", "ID00001")

	baseHelper := NewIpay88Helper()

	currencyData, _ := constants.GetCurrencyLabel(constants.CurrencyIDR)
	signature, err := baseHelper.generateSignature("A00000001", 3000, *currencyData)
	if err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}

	if string(signature) != "3ee767e49d29c46b6187db2fe511287ccd986e874a057fa2e1bb222442b68f63" {
		t.Fail()
		t.Log("signature is not the same")
		return
	}

	t.Log("success generate signature")
	return
}
