package ipay88_helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"

	"github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/constants"
	ipay88Model "github.com/fari-99/go-helper/payment_gateways/ipay88_helpers/models"
	"github.com/fari-99/go-helper/payment_gateways/models"
)

func (base *BaseIpay88Helper) CreatPaymentRequest() (*models.Invoices, error) {
	currencyLabel, _ := constants.GetCurrencyLabel(constants.CurrencyIDR)
	requestType, _ := constants.GetRequestTypeLabel(constants.RequestTypeRedirect)
	encodingType, _ := constants.GetEncodingLabel(constants.EncodingUTF_8)

	transactionDetails := base.TransactionModel
	transactionUser := base.TransactionUser
	transactionUuid := transactionDetails.TransactionUuid

	paymentRequestData, err := base.generatePaymentRequestData()
	if err != nil {
		return nil, err
	}

	totalPrice := paymentRequestData.TotalItemPrice + paymentRequestData.TotalAdditionalFee
	totalPriceUint := cast.ToUint64(totalPrice)
	signature, err := base.generateSignature(transactionUuid, totalPriceUint, *currencyLabel)
	if err != nil {
		return nil, err
	}

	// Response page URL is the page at the merchant website that will receive payment status from iPay88 OPSG.
	responseUrl, _ := constants.GetApiUrl(constants.PaymentResponseUrl)

	// Backend response page
	callbackUrl, _ := constants.GetApiUrl(constants.PaymentBackendUrl)

	paymentRequestInput := ipay88Model.PaymentRequests{
		APIVersion:   constants.ApiVersions,
		MerchantCode: os.Getenv("IPAY88_MERCHANT_CODE"),
		PaymentID:    os.Getenv("IPAY88_MERCHANT_KEY"),
		Currency:     *currencyLabel,
		RefNo:        transactionUuid,
		Amount:       cast.ToString(totalPriceUint),
		ProdDesc:     "",
		UserName:     fmt.Sprintf("%s %s", transactionUser.FirstName, transactionUser.LastName),
		UserEmail:    transactionUser.EmailAddress,
		UserContact:  transactionUser.Phone,
		RequestType:  requestType,
		Remark:       nil,
		Lang:         encodingType,
		ResponseURL:  *responseUrl,
		BackendURL:   *callbackUrl,
		Signature:    signature,
		// FullTransactionAmount: nil,
		// MiscFee:               nil,
		ItemTransactions: paymentRequestData.Items,
		ShippingAddress:  paymentRequestData.ShippingAddress,
		BillingAddress:   paymentRequestData.BillingAddress,
		Sellers:          paymentRequestData.Sellers,
		SettingField:     nil,
	}

	url, err := constants.GetIpay88Url(constants.Ipay88PaymentRequestUrl)
	if err != nil {
		return nil, err
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(paymentRequestInput).
		Post(*url)
	if err != nil {
		return nil, err
	}

	var responseData ipay88Model.PaymentRequestResponse
	_ = json.Unmarshal(resp.Body(), &responseData)

	redirectUrl, _ := constants.GetIpay88Url(constants.Ipay88PaymentRedirectUrl)
	redirectParams := map[string]string{
		"CheckoutID": responseData.CheckoutID,
		"Signature":  responseData.Signature,
	}

	invoiceRespMarshal, _ := json.Marshal(responseData)
	redirectParamMarshal, _ := json.Marshal(redirectParams)
	invoiceModel := models.Invoices{
		PaymentGatewayID:  transactionDetails.PaymentGatewayID,
		PaymentMethodType: transactionDetails.PaymentMethodType,
		PaymentMethodCode: transactionDetails.PaymentMethodCode,
		TransactionUuid:   transactionUuid,
		// Status:            config.GetConstStatus("STATUS_ACTIVE", nil), // TODO: adding status

		TotalPrice:     totalPrice,
		Identifier:     responseData.RefNo,
		ExpiredAt:      responseData.TransactionExpiryDate,
		RedirectUrl:    *redirectUrl,
		RedirectParams: string(redirectParamMarshal),
		ResponseJson:   string(invoiceRespMarshal),
	}

	return &invoiceModel, nil
}
