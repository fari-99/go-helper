package models

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

type BackendPostParams struct {
    MerchantCode       string `json:"MerchantCode"`
    PaymentID          string `json:"PaymentId"`
    RefNo              string `json:"RefNo"`
    Amount             string `json:"Amount"`
    Currency           string `json:"Currency"`
    Remark             string `json:"Remark"`
    TransID            string `json:"TransId"`
    AuthCode           string `json:"AuthCode"`
    TransactionStatus  string `json:"TransactionStatus"`
    ErrDesc            string `json:"ErrDesc"`
    Signature          string `json:"Signature"`
    IssuerBank         string `json:"IssuerBank"`
    PaymentDate        string `json:"PaymentDate"`
    Xfield1            string `json:"Xfield1"`
    DCCConversionRate  string `json:"DCCConversionRate"`
    OriginalAmount     string `json:"OriginalAmount"`
    OriginalCurrency   string `json:"OriginalCurrency"`
    SettlementAmount   string `json:"SettlementAmount"`
    SettlementCurrency string `json:"SettlementCurrency"`
    Binbank            string `json:"Binbank"`
}

func (model BackendPostParams) ValidateSignature(merchantKey string, merchantCode string) (Message, error) {
    signature := model.Signature

    if merchantCode != model.MerchantCode {
        return Message{
            Indonesia: "Kode Merchant parameter tidak sama dengan Kode Merchant yang saat ini dimiliki",
            English:   "Merchant code params is not the same as your merchant code",
        }, fmt.Errorf("merchant code params [%s] is not the same as your merchant code [%s]", model.MerchantCode, merchantCode)
    }

    checkSignature := "||"
    checkSignature += fmt.Sprintf("%s||", merchantKey)
    checkSignature += fmt.Sprintf("%s||", model.MerchantCode)
    checkSignature += fmt.Sprintf("%s||", model.PaymentID)
    checkSignature += fmt.Sprintf("%s||", model.RefNo)
    checkSignature += fmt.Sprintf("%s||", model.Amount)
    checkSignature += fmt.Sprintf("%s||", model.Currency)
    checkSignature += fmt.Sprintf("%s", model.TransactionStatus)
    checkSignature += "||"

    hs := sha256.New()
    hs.Write([]byte(checkSignature))
    calcSignature := hex.EncodeToString(hs.Sum(nil))
    if signature != calcSignature {
        return Message{
            Indonesia: "Signature parameter tidak sama dengan Signature yang di proses",
            English:   "Signature is not the same",
        }, fmt.Errorf("signature params [%s] is not the same [%s]", signature, calcSignature)
    }

    return Message{
        Indonesia: "Status Received",
        English:   "Pembayaran diterima",
    }, nil
}
