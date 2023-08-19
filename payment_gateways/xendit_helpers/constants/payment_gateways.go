package constants

import (
	"fmt"
)

const (
	ModuleInvoices = iota + 1
	ModulePayments
)

const (
	PaymentTypeCreditCards = iota + 1
	PaymentTypeEWallet
	PaymentTypePayLater
	PaymentTypeQRCodes
	PaymentTypeDirectDebit
	PaymentTypeVirtualAccount
	PaymentTypeRetailOutletOTC
)

func GetPaymentTypes() map[PaymentTypes]PaymentTypeDetail {
	data := map[PaymentTypes]PaymentTypeDetail{
		PaymentTypeCreditCards: {
			Name:           "Credit Cards",
			Label:          "Credit Cards",
			Code:           "CARD",
			PaymentMethods: GetCreditCards(),
		},
		PaymentTypeEWallet: {
			Name:           "E Wallet",
			Label:          "E-Wallet",
			Code:           "EWALLET",
			PaymentMethods: GetEWallets(),
		},
		PaymentTypePayLater: {
			Name:           "Pay Later",
			Label:          "Pay Later",
			Code:           "PAY_LATER",
			PaymentMethods: GetPayLater(),
		},
		PaymentTypeQRCodes: {
			Name:           "QR Codes",
			Label:          "QR Codes",
			Code:           "QR_CODE",
			PaymentMethods: GetQrCodes(),
		},
		PaymentTypeDirectDebit: {
			Name:           "Direct Debit",
			Label:          "Direct Debit",
			Code:           "DIRECT_DEBIT",
			PaymentMethods: GetDirectDebits(),
		},
		PaymentTypeVirtualAccount: {
			Name:           "Virtual Accounts",
			Label:          "Virtual Account",
			Code:           "VIRTUAL_ACCOUNT",
			PaymentMethods: GetVirtualAccounts(),
		},
		PaymentTypeRetailOutletOTC: {
			Name:           "Retail Outlets (OTC)",
			Label:          "Retail Outlets (OTC)",
			Code:           "OVER_THE_COUNTER",
			PaymentMethods: GetRetailOutletOTC(),
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

type PaymentMethods map[PaymentTypes]PaymentMethodList
type PaymentMethodList map[Country]PaymentMethodDetail
type PaymentMethodDetail struct {
	Name  string         `json:"name"`
	Label string         `json:"label"`
	Code  map[int]string `json:"code"`
}

const (
	creditCards = iota + 1
)

func GetCreditCards() PaymentMethods {
	creditCartMethod := PaymentMethodDetail{
		Name:  "Credit Cards",
		Label: "Credit Cards",
		Code: map[int]string{
			ModuleInvoices: "CREDIT_CARD",
			ModulePayments: "CARD",
		},
	}

	data := PaymentMethods{
		creditCards: PaymentMethodList{
			Indonesia: creditCartMethod,
		},
	}

	return data
}

func GetCreditCardDetail(dataType PaymentTypes) (*PaymentMethodDetail, error) {
	allData := GetCreditCards()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		paymentMethodDetail = value[Indonesia]

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("credit card type [%d] not found", dataType)
	}
}

// Indonesia
const (
	eWalletOvo = iota + 1
	eWalletDana
	eWalletLinkAja
	eWalletShopeePay
	eWalletAstraPay
	eWalletJeniusPay
	eWalletSakuKu
)

// Philippines
const (
	eWalletMaya = iota + 1
	eWalletGCash
	eWalletGrabPay
	eWalletShopeePayPhilippines
)

// Vietnam
const (
	eWalletShopeePayVietnam = iota + 1
	eWalletMomo
	eWalletZaloPay
	eWalletAppota
	eWalletVnptWallet
	eWalletViettelPay
)

func GetEWallets() PaymentMethods {
	ovo := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "OVO",
			ModulePayments: "OVO",
		},
	}

	dana := PaymentMethodDetail{
		Name:  "Dana",
		Label: "Dana",
		Code: map[int]string{
			ModuleInvoices: "DANA",
			ModulePayments: "DANA",
		},
	}

	linkAja := PaymentMethodDetail{
		Name:  "LinkAja",
		Label: "LinkAja",
		Code: map[int]string{
			ModuleInvoices: "LINKAJA",
			ModulePayments: "LINKAJA",
		},
	}

	shopeePay := PaymentMethodDetail{
		Name:  "ShopeePay",
		Label: "ShopeePay",
		Code: map[int]string{
			ModuleInvoices: "SHOPEEPAY",
			ModulePayments: "SHOPEEPAY",
		},
	}

	astraPay := PaymentMethodDetail{
		Name:  "AstraPay",
		Label: "AstraPay",
		Code: map[int]string{
			ModuleInvoices: "ASTRAPAY",
			ModulePayments: "ASTRAPAY",
		},
	}

	jeniusPay := PaymentMethodDetail{
		Name:  "JeniusPay",
		Label: "JeniusPay",
		Code: map[int]string{
			ModuleInvoices: "JENIUSPAY",
			ModulePayments: "JENIUSPAY",
		},
	}

	sakuKu := PaymentMethodDetail{
		Name:  "SakuKu",
		Label: "SakuKu",
		Code: map[int]string{
			ModuleInvoices: "SAKUKU",
			ModulePayments: "", // TODO: need information
		},
	}

	data := PaymentMethods{
		Indonesia: {
			eWalletOvo:       ovo,
			eWalletDana:      dana,
			eWalletLinkAja:   linkAja,
			eWalletShopeePay: shopeePay,
			eWalletAstraPay:  astraPay,
			eWalletJeniusPay: jeniusPay,
			eWalletSakuKu:    sakuKu,
		},

		Philippines: {
			// eWalletMaya:    "Maya",
			// eWalletGCash:   "GCash",
			// eWalletGrabPay: "GrabPay",
			// eWalletShopeePay : "",
		},

		Vietnam: {
			// eWalletShopeePay : "",
			// eWalletMomo:       "Momo",
			// eWalletZaloPay:    "ZaloPay",
			// eWalletAppota:     "Appota",
			// eWalletVnptWallet: "VnptWallet",
			// eWalletViettelPay: "ViettelPay",
		},
	}

	return data
}

func GetEWalletLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetEWallets()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("e-wallet type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("e-wallet type [%d] not found", dataType)
	}
}

// Indonesia
const (
	payLaterKredivo = iota + 1
	payLaterAkuLaku
	payLaterUangMe
	payLaterIndoDana
	payLaterAtome
)

// Philippines
const (
	payLaterBillease = iota + 1
	payLaterCashalo
	payLaterAtomePhilippines
)

func GetPayLater() PaymentMethods {
	kredivo := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "KREDIVO",
			ModulePayments: "", // NOT AVAILABLE
		},
	}

	akuLaku := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "AKULAKU",
			ModulePayments: "", // NOT AVAILABLE
		},
	}

	uangMe := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "UANGME",
			ModulePayments: "", // NOT AVAILABLE
		},
	}

	indoDana := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "INDODANA",
			ModulePayments: "", // NOT AVAILABLE
		},
	}

	atome := PaymentMethodDetail{
		Name:  "Ovo",
		Label: "Ovo",
		Code: map[int]string{
			ModuleInvoices: "ATOME",
			ModulePayments: "", // NOT AVAILABLE
		},
	}

	data := PaymentMethods{
		Indonesia: {
			payLaterKredivo:  kredivo,
			payLaterAkuLaku:  akuLaku,
			payLaterUangMe:   uangMe,
			payLaterIndoDana: indoDana,
			payLaterAtome:    atome,
		},

		Philippines: {
			// payLaterBillease: "Billease",
			// payLaterCashalo:  "Cashalo",
			// payLaterAtome
		},
	}

	return data
}

func GetPayLaterLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetPayLater()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("pay later type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("pay later type [%d] not found", dataType)
	}
}

const (
	qrCode = iota + 1
	qrCodeDana
	qrCodeLinkAja
)

func GetQrCodes() PaymentMethods {
	qrCodeDefault := PaymentMethodDetail{
		Name:  "QRIS",
		Label: "QRIS",
		Code: map[int]string{
			ModuleInvoices: "QRIS",
			ModulePayments: "", // TODO: need more information
		},
	}

	dana := PaymentMethodDetail{
		Name:  "Dana",
		Label: "QRIS - Dana",
		Code: map[int]string{
			ModuleInvoices: "QRIS",
			ModulePayments: "DANA",
		},
	}

	linkAja := PaymentMethodDetail{
		Name:  "LinkAja",
		Label: "QRIS - LinkAja",
		Code: map[int]string{
			ModuleInvoices: "QRIS",
			ModulePayments: "LINKAJA",
		},
	}

	data := PaymentMethods{
		Indonesia: {
			qrCode:        qrCodeDefault,
			qrCodeDana:    dana,
			qrCodeLinkAja: linkAja,
		},
	}

	return data
}

func GetQrCodeLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetQrCodes()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("qrcode type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("qr code type [%d] not found", dataType)
	}
}

const (
	// ID
	directDebitBri = iota + 1
	directDebitBcaOneKlik
	directDebitBcaKlikPay
	directDebitMandiri
)

// Philippines
const (
	directDebitBpi = iota + 1
	directDebitUnionBank
)

// Thailand
const (
	directDebitBbl = iota + 1
	directDebitKrungsri
	directDebitKtb
	directDebitScb
)

// Malaysia
const (
	directDebitFpx = iota + 1
)

func GetDirectDebits() PaymentMethods {
	bri := PaymentMethodDetail{
		Name:  "BRI Direct Debit",
		Label: "BRI Direct Debit",
		Code: map[int]string{
			ModuleInvoices: "DD_BRI",
			ModulePayments: "BRI",
		},
	}

	bcaOneKlik := PaymentMethodDetail{
		Name:  "BCA OneKlik",
		Label: "BCA OneKlik",
		Code: map[int]string{
			ModuleInvoices: "DD_BCA_ONEKLIK",
			ModulePayments: "", // TODO: need more information
		},
	}

	bcaKlikPay := PaymentMethodDetail{
		Name:  "BCA KlikPay",
		Label: "BCA KlikPay",
		Code: map[int]string{
			ModuleInvoices: "DD_BCA_KLIKPAY",
			ModulePayments: "", // TODO: need more information
		},
	}

	mandiri := PaymentMethodDetail{
		Name:  "Mandiri",
		Label: "Mandiri",
		Code: map[int]string{
			ModuleInvoices: "", // TODO: need more information
			ModulePayments: "MANDIRI",
		},
	}

	data := PaymentMethods{
		Indonesia: {
			directDebitBri:        bri,
			directDebitBcaOneKlik: bcaOneKlik,
			directDebitBcaKlikPay: bcaKlikPay,
			directDebitMandiri:    mandiri,
		},
		Philippines: {
			// directDebitBpi : "",
			// directDebitUnionBank : "",
		},
		Thailand: {
			// directDebitBbl : "",
			// directDebitKrungsri : "",
			// directDebitKtb : "",
			// directDebitScb : "",
		},
		Malaysia: {
			// directDebitFpx : "",
		},
	}

	return data
}

func GetDirectDebitLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetDirectDebits()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("direct debit type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("direct debit type [%d] not found", dataType)
	}
}

// Indonesia
const (
	virtualAccountBca = iota + 1
	virtualAccountBni
	virtualAccountBri
	virtualAccountBjb
	virtualAccountBsi
	virtualAccountBnc
	virtualAccountCimb
	virtualAccountDbs
	virtualAccountMandiri
	virtualAccountPermata
	virtualAccountSahabatSampoerna
	virtualAccountArtaJasa
)

// Philippines
const (
	virtualAccountPv = iota + 1
	virtualAccountVietCapital
	virtualAccountWoori
)

func GetVirtualAccounts() PaymentMethods {
	bca := PaymentMethodDetail{
		Name:  "Bank Central Asia",
		Label: "Bank Central Asia (BCA)",
		Code: map[int]string{
			ModuleInvoices: "BCA",
			ModulePayments: "BCA",
		},
	}

	bni := PaymentMethodDetail{
		Name:  "Bank National Indonesia",
		Label: "Bank National Indonesia (BNI)",
		Code: map[int]string{
			ModuleInvoices: "BNI",
			ModulePayments: "BNI",
		},
	}

	bri := PaymentMethodDetail{
		Name:  "Bank Rakyat Indonesia",
		Label: "Bank Rakyat Indonesia (BRI)",
		Code: map[int]string{
			ModuleInvoices: "BRI",
			ModulePayments: "BRI",
		},
	}

	bjb := PaymentMethodDetail{
		Name:  "Bank Pembangunan Daerah (Bank BJB)",
		Label: "Bank Pembangunan Daerah (Bank BJB)",
		Code: map[int]string{
			ModuleInvoices: "BJB",
			ModulePayments: "BJB",
		},
	}

	bsi := PaymentMethodDetail{
		Name:  "Bank Syariah Indonesia (BSI)",
		Label: "Bank Syariah Indonesia (BSI)",
		Code: map[int]string{
			ModuleInvoices: "BSI",
			ModulePayments: "BSI",
		},
	}

	bnc := PaymentMethodDetail{
		Name:  "Bank Neo Commerce (BNC)",
		Label: "Bank Neo Commerce (BNC)",
		Code: map[int]string{
			ModuleInvoices: "BNC",
			ModulePayments: "", // TODO: need more information
		},
	}

	cimb := PaymentMethodDetail{
		Name:  "Bank CIMB Niaga (CIMB)",
		Label: "Bank CIMB Niaga (CIMB)",
		Code: map[int]string{
			ModuleInvoices: "CIMB",
			ModulePayments: "CIMB",
		},
	}

	dbs := PaymentMethodDetail{
		Name:  "DBS Bank",
		Label: "DBS Bank",
		Code: map[int]string{
			ModuleInvoices: "DBS",
			ModulePayments: "", // TODO: need more information
		},
	}

	mandiri := PaymentMethodDetail{
		Name:  "Bank Mandiri",
		Label: "Bank Mandiri",
		Code: map[int]string{
			ModuleInvoices: "MANDIRI",
			ModulePayments: "MANDIRI",
		},
	}

	permata := PaymentMethodDetail{
		Name:  "Permata Bank (Bank Permata)",
		Label: "Permata Bank (Bank Permata)",
		Code: map[int]string{
			ModuleInvoices: "PERMATA",
			ModulePayments: "PERMATA",
		},
	}

	sahabatSampoerna := PaymentMethodDetail{
		Name:  "Bank Sahabar Sampoerna",
		Label: "Bank Sahabar Sampoerna",
		Code: map[int]string{
			ModuleInvoices: "SAHABAT_SAMPOERNA",
			ModulePayments: "SAHABAT_SAMPOERNA",
		},
	}

	artaJasa := PaymentMethodDetail{
		Name:  "PT. Artajasa Pembayaran Elektronis",
		Label: "PT. Artajasa Pembayaran Elektronis",
		Code: map[int]string{
			ModuleInvoices: "", // TODO: need more information
			ModulePayments: "ARTAJASA",
		},
	}

	data := PaymentMethods{
		Indonesia: {
			virtualAccountBca:              bca,
			virtualAccountBni:              bni,
			virtualAccountBri:              bri,
			virtualAccountBjb:              bjb,
			virtualAccountBsi:              bsi,
			virtualAccountBnc:              bnc,
			virtualAccountCimb:             cimb,
			virtualAccountDbs:              dbs,
			virtualAccountMandiri:          mandiri,
			virtualAccountPermata:          permata,
			virtualAccountSahabatSampoerna: sahabatSampoerna,
			virtualAccountArtaJasa:         artaJasa,
		},
		Philippines: {
			// virtualAccountPv : "",
			// virtualAccountVietCapital : "",
			// virtualAccountWoori : "",
		},
	}

	return data
}

func GetVirtualAccountLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetVirtualAccounts()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("virtual account type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("virtual account type [%d] not found", dataType)
	}
}

// Indonesia
const (
	retailOutletOtcAlfamart = iota + 1
	retailOutletOtcIndomaret
)

// Philippines
const (
	retailOutletOtc7Eleven = iota + 1
	retailOutletOtcCebuana
	retailOutletOtcEcPay
	retailOutletOtcPalawan
	retailOutletOtcMlhuillier
	retailOutletOtcDragonPay
)

func GetRetailOutletOTC() PaymentMethods {
	alfamart := PaymentMethodDetail{
		Name:  "Alfamart",
		Label: "Alfamart",
		Code: map[int]string{
			ModuleInvoices: "ALFAMART",
			ModulePayments: "ALFAMART",
		},
	}

	indomaret := PaymentMethodDetail{
		Name:  "Indomaret",
		Label: "Indomaret",
		Code: map[int]string{
			ModuleInvoices: "INDOMARET",
			ModulePayments: "INDOMARET",
		},
	}

	data := PaymentMethods{
		Indonesia: {
			retailOutletOtcAlfamart:  alfamart,
			retailOutletOtcIndomaret: indomaret,
		},
		Philippines: {
			// retailOutletOtc7Eleven : "",
			// retailOutletOtcCebuana : "",
			// retailOutletOtcEcPay : "",
			// retailOutletOtcPalawan : "",
			// retailOutletOtcMlhuillier : "",
			// retailOutletOtcDragonPay : "",
		},
	}

	return data
}

func GetRetailOutletOTCLabel(dataType PaymentTypes, country Country) (*PaymentMethodDetail, error) {
	allData := GetVirtualAccounts()
	if value, ok := allData[dataType]; ok {
		var paymentMethodDetail PaymentMethodDetail
		if country == 0 {
			paymentMethodDetail = value[Indonesia]
		} else {
			if paymentMethod, ok := value[country]; ok {
				paymentMethodDetail = paymentMethod
			} else {
				return nil, fmt.Errorf("retail outlet (otc) type [%d] for that country is not found", dataType)
			}
		}

		return &paymentMethodDetail, nil
	} else {
		return nil, fmt.Errorf("retail outlet (otc) type [%d] not found", dataType)
	}
}
