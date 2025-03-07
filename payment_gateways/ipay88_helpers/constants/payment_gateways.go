package constants

import (
    "fmt"
)

const (
    PaymentTypeCreditCards = iota + 1
    PaymentTypeOnlineBanking
    PaymentTypeVirtualAccount
    PaymentTypeEWallet
    PaymentTypeQRCode
    PaymentTypeOverTheCounter
    PaymentTypeOnlineCredit
    PaymentTypeOthers
)

func GetPaymentTypes() map[PaymentTypes]PaymentTypeDetail {
    data := map[PaymentTypes]PaymentTypeDetail{
        PaymentTypeCreditCards: {
            Name:           "Credit Cards",
            Label:          "Credit Cards",
            Code:           "CREDIT_CARDS",
            PaymentMethods: GetCreditCards(),
        },
        PaymentTypeEWallet: {
            Name:           "E Wallet",
            Label:          "E-Wallet",
            Code:           "EWALLET",
            PaymentMethods: GetEWallets(),
        },
        PaymentTypeOnlineBanking: {
            Name:           "Online Banking",
            Label:          "Online Banking",
            Code:           "ONLINE_BANKING",
            PaymentMethods: GetOnlineBanking(),
        },
        PaymentTypeVirtualAccount: {
            Name:           "QR Codes",
            Label:          "QR Codes",
            Code:           "QR_CODE",
            PaymentMethods: GetVirtualAccounts(),
        },
        PaymentTypeQRCode: {
            Name:           "Virtual Accounts",
            Label:          "Virtual Account",
            Code:           "VIRTUAL_ACCOUNT",
            PaymentMethods: getQRCodes(),
        },
        PaymentTypeOverTheCounter: {
            Name:           "Direct Debit",
            Label:          "Direct Debit",
            Code:           "DIRECT_DEBIT",
            PaymentMethods: getOverTheCounters(),
        },
        PaymentTypeOnlineCredit: {
            Name:           "Retail Outlets (OTC)",
            Label:          "Retail Outlets (OTC)",
            Code:           "OVER_THE_COUNTER",
            PaymentMethods: getOnlineCredits(),
        },
        PaymentTypeOthers: {
            Name:           "Others",
            Label:          "Others",
            Code:           "OTHERS",
            PaymentMethods: getOthers(),
        },
    }

    return data
}

func GetPaymentTypeDetail(dataType PaymentTypes) (*PaymentTypeDetail, error) {
    allData := GetPaymentTypes()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("payment type [%d] not found", dataType)
    }
}

type PaymentTypes int
type PaymentTypeDetail struct {
    Name           string         `json:"name"`
    Label          string         `json:"label"`
    Code           string         `json:"code"`
    PaymentMethods PaymentMethods `json:"payment_methods"`
}

type PaymentMethods map[PaymentTypes]PaymentMethodDetail
type PaymentMethodDetail struct {
    Name  string `json:"name"`
    Label string `json:"label"`
    Code  string `json:"code"`
}

const (
    creditCardBCA = iota + 1
    creditCardBRI
    creditCardCIMB
    creditCardCIMBAuthorization
    creditCardCIMBIPG
    creditCardDanamon
    creditCardMandiri
    creditCardMaybank
    creditCardUnionPay
    creditCardUOB
    creditCardGPN
)

func GetCreditCards() PaymentMethods {
    bca := PaymentMethodDetail{
        Name:  "BCA",
        Label: "BCA",
        Code:  "52",
    }

    bri := PaymentMethodDetail{
        Name:  "BRI",
        Label: "BRI",
        Code:  "35",
    }

    cimb := PaymentMethodDetail{
        Name:  "CIMB",
        Label: "CIMB",
        Code:  "42",
    }

    cimbAuth := PaymentMethodDetail{
        Name:  "CIMB Authorization",
        Label: "CIMB Authorization",
        Code:  "56",
    }

    cimbIpg := PaymentMethodDetail{
        Name:  "CIMB IPG",
        Label: "CIMB IPG",
        Code:  "34",
    }

    danamon := PaymentMethodDetail{
        Name:  "Danamon",
        Label: "Danamon",
        Code:  "45",
    }

    mandiri := PaymentMethodDetail{
        Name:  "Mandiri",
        Label: "Mandiri",
        Code:  "53",
    }

    maybank := PaymentMethodDetail{
        Name:  "Maybank",
        Label: "Maybank",
        Code:  "43",
    }

    unionPay := PaymentMethodDetail{
        Name:  "UnionPay",
        Label: "UnionPay",
        Code:  "54",
    }

    uob := PaymentMethodDetail{
        Name:  "UOB",
        Label: "UOB",
        Code:  "46",
    }

    gpn := PaymentMethodDetail{
        Name:  "GPN",
        Label: "GPN",
        Code:  "49",
    }

    data := PaymentMethods{
        creditCardBCA:               bca,
        creditCardBRI:               bri,
        creditCardCIMB:              cimb,
        creditCardCIMBAuthorization: cimbAuth,
        creditCardCIMBIPG:           cimbIpg,
        creditCardDanamon:           danamon,
        creditCardMandiri:           mandiri,
        creditCardMaybank:           maybank,
        creditCardUnionPay:          unionPay,
        creditCardUOB:               uob,
        creditCardGPN:               gpn,
    }

    return data
}

func GetCreditCardDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := GetCreditCards()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("credit card type [%d] not found", dataType)
    }
}

const (
    onlineBankingBCAKlikPay = iota + 1
    onlineBankingCIMBClicks
    onlineBankingMuamalatIB
    onlineBankingDanamonOnlineBanking
)

func GetOnlineBanking() PaymentMethods {
    bcaKlikPay := PaymentMethodDetail{
        Name:  "BCA KlikPay",
        Label: "BCA KlikPay",
        Code:  "8",
    }

    cimbClick := PaymentMethodDetail{
        Name:  "CIMB Clicks",
        Label: "CIMB Clicks",
        Code:  "11",
    }

    muamalatIB := PaymentMethodDetail{
        Name:  "Muamalat IB",
        Label: "Muamalat IB",
        Code:  "14",
    }

    danamonOnlineBanking := PaymentMethodDetail{
        Name:  "Danamon Online Banking",
        Label: "Danamon Online Banking",
        Code:  "23",
    }

    data := PaymentMethods{
        onlineBankingBCAKlikPay:           bcaKlikPay,
        onlineBankingCIMBClicks:           cimbClick,
        onlineBankingMuamalatIB:           muamalatIB,
        onlineBankingDanamonOnlineBanking: danamonOnlineBanking,
    }

    return data
}

func GetOnlineBankingDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := GetOnlineBanking()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("online banking type [%d] not found", dataType)
    }
}

const (
    virtualAccountMaybank = iota + 1
    virtualAccountPermata
    virtualAccountMandiri
    virtualAccountBCA
    virtualAccountBRI
    virtualAccountBNI
    virtualAccountCIMB
)

func GetVirtualAccounts() PaymentMethods {
    maybank := PaymentMethodDetail{
        Name:  "Maybank VA",
        Label: "Maybank VA",
        Code:  "9",
    }

    permata := PaymentMethodDetail{
        Name:  "Permata VA",
        Label: "Permata VA",
        Code:  "31",
    }

    mandiri := PaymentMethodDetail{
        Name:  "Mandiri VA",
        Label: "Mandiri VA",
        Code:  "17",
    }

    bca := PaymentMethodDetail{
        Name:  "BCA VA",
        Label: "BCA VA",
        Code:  "25",
    }

    bri := PaymentMethodDetail{
        Name:  "BRI VA",
        Label: "BRI VA",
        Code:  "61",
    }

    bni := PaymentMethodDetail{
        Name:  "BNI VA",
        Label: "BNI VA",
        Code:  "26",
    }

    cimb := PaymentMethodDetail{
        Name:  "CIMB VA",
        Label: "CIMB VA",
        Code:  "86",
    }

    data := PaymentMethods{
        virtualAccountMaybank: maybank,
        virtualAccountPermata: permata,
        virtualAccountMandiri: mandiri,
        virtualAccountBCA:     bca,
        virtualAccountBRI:     bri,
        virtualAccountBNI:     bni,
        virtualAccountCIMB:    cimb,
    }

    return data
}

func GetVirtualAccountDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := GetVirtualAccounts()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("virtual account type [%d] not found", dataType)
    }
}

const (
    eWalletOvo = iota + 1
    eWalletDana
    eWalletLinkAja
    eWalletShopeePay
)

func GetEWallets() PaymentMethods {
    ovo := PaymentMethodDetail{
        Name:  "OVO",
        Label: "OVO",
        Code:  "63",
    }

    dana := PaymentMethodDetail{
        Name:  "DANA",
        Label: "DANA",
        Code:  "77",
    }

    linkAja := PaymentMethodDetail{
        Name:  "LinkAja",
        Label: "LinkAja",
        Code:  "13",
    }

    shopeePay := PaymentMethodDetail{
        Name:  "ShopeePay JumpApp",
        Label: "ShopeePay JumpApp",
        Code:  "76",
    }

    data := PaymentMethods{
        eWalletOvo:       ovo,
        eWalletDana:      dana,
        eWalletLinkAja:   linkAja,
        eWalletShopeePay: shopeePay,
    }

    return data
}

func GetEWalletDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := GetEWallets()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("e-wallet type [%d] not found", dataType)
    }
}

const (
    qrCodeShopeePay = iota + 0
    qrCodeNobuBank
)

func getQRCodes() PaymentMethods {
    shopee := PaymentMethodDetail{
        Name:  "ShopeePay QR",
        Label: "ShopeePay QR",
        Code:  "75",
    }

    nobu := PaymentMethodDetail{
        Name:  "Nobu Bank QRIS",
        Label: "Nobu Bank QRIS",
        Code:  "78",
    }

    data := PaymentMethods{
        qrCodeShopeePay: shopee,
        qrCodeNobuBank:  nobu,
    }

    return data
}

func getQRCodeDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := getQRCodes()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("qr-code type [%d] not found", dataType)
    }
}

const (
    overTheCounterAlfamart = iota + 1
    overTheCounterIndomaret
)

func getOverTheCounters() PaymentMethods {
    alfamart := PaymentMethodDetail{
        Name:  "Alfamart",
        Label: "Alfamart",
        Code:  "60",
    }

    indomaret := PaymentMethodDetail{
        Name:  "Indomaret",
        Label: "Indomaret",
        Code:  "65",
    }

    data := PaymentMethods{
        overTheCounterAlfamart:  alfamart,
        overTheCounterIndomaret: indomaret,
    }

    return data
}

func getOverTheCounterDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := getOverTheCounters()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("over-the-counter type [%d] not found", dataType)
    }
}

const (
    onlineCreditAkuLaku = iota + 1
    onlineCreditIndoDana
    onlineCreditKredivo
)

func getOnlineCredits() PaymentMethods {
    akulaku := PaymentMethodDetail{
        Name:  "Akulaku",
        Label: "Akulaku",
        Code:  "71",
    }

    indodana := PaymentMethodDetail{
        Name:  "Indodana",
        Label: "Indodana",
        Code:  "70",
    }

    kredivo := PaymentMethodDetail{
        Name:  "Kredivo",
        Label: "Kredivo",
        Code:  "55",
    }

    data := PaymentMethods{
        onlineCreditAkuLaku:  akulaku,
        onlineCreditIndoDana: indodana,
        onlineCreditKredivo:  kredivo,
    }

    return data
}

func getOnlineCreditDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := getOnlineCredits()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("over-the-counter type [%d] not found", dataType)
    }
}

const (
    paypalPayment = iota + 1
)

func getOthers() PaymentMethods {
    paypal := PaymentMethodDetail{
        Name:  "Paypal",
        Label: "Paypal",
        Code:  "6",
    }

    data := PaymentMethods{
        paypalPayment: paypal,
    }

    return data
}

func getOthersDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
    allData := getOthers()
    if value, ok := allData[dataType]; ok {
        return &value, nil
    } else {
        return nil, fmt.Errorf("other payment type [%d] not found", dataType)
    }
}
