package virtual_accounts

import (
    "fmt"
    "math/rand"
    "os"
    "time"

    "github.com/spf13/cast"
)

const BankBCA = "BCA"         // Bank Central Asia (BCA)
const BankCIMB = "CIMB"       // Bank CIMB (CIMB)
const BankDBS = "DBS"         // Bank Development Bank of Singapore (DBS)
const BankBJB = "BJB"         // Bank Jabar Banten (BJB)
const BankMandiri = "Mandiri" // Bank Mandiri
const BankBNI = "BNI"         // Bank Negara Indonesia (BNI)
const BankBNC = "BNC"         // Bank Neo Commerce (BNC)
const BankPermata = "Permata" // Bank Permata
const BankBRI = "BRI"         // Bank Rakyat Indonesia (BRI)
const BankBSS = "BSS"         // Bank Sahabat Sampoerna (BSS)
const BankBSI = "BSI"         // Bank Syariah Indonesia (BSI)

func GetXenditVirtualAccountBanks() map[string]bool {
    return map[string]bool{
        BankBCA:     true,
        BankCIMB:    true,
        BankDBS:     true,
        BankBJB:     true,
        BankMandiri: true,
        BankBNI:     true,
        BankBNC:     true,
        BankPermata: true,
        BankBRI:     true,
        BankBSS:     true,
        BankBSI:     true,
    }
}

var bcaMerchantCode = []int{7007, 38165, 38166}
var cimbMerchantCode = []int{93490}
var dbsMerchantCode = []int{9488}
var bjbMerchantCode = []int{12345}
var mandiriMerchantCode = []int{88608, 88908}
var bniMerchantCode = []int{8808, 8930, 7151, 7152}
var bncMerchantCode = []int{90100011}
var permataMerchantCode = []int{8214, 7293}
var briMerchantCode = []int{26215, 92001, 13281, 13282, 13404, 13405}
var bssMerchantCode = []int{40102}
var bsiMerchantCode = []int{9347, 9655}

func getMerchantCode(bankCode string) string {
    getXenditBank := GetXenditVirtualAccountBanks()
    if isActive, ok := getXenditBank[bankCode]; !ok {
        panic(fmt.Sprintf("not such bank code [%s]", bankCode))
    } else if !isActive {
        panic(fmt.Sprintf("bank [%s] is not active", bankCode))
    }

    var bankMerchantCode []int
    switch bankCode {
    case BankBCA:
        bankMerchantCode = bcaMerchantCode
    case BankCIMB:
        bankMerchantCode = cimbMerchantCode
    case BankDBS:
        bankMerchantCode = dbsMerchantCode
    case BankBJB:
        bankMerchantCode = bjbMerchantCode
    case BankMandiri:
        bankMerchantCode = mandiriMerchantCode
    case BankBNI:
        bankMerchantCode = bniMerchantCode
    case BankBNC:
        bankMerchantCode = bncMerchantCode
    case BankPermata:
        bankMerchantCode = permataMerchantCode
    case BankBRI:
        bankMerchantCode = briMerchantCode
    case BankBSS:
        bankMerchantCode = bssMerchantCode
    case BankBSI:
        bankMerchantCode = bsiMerchantCode
    default:
        panic(fmt.Sprintf("not such bank code [%s]", bankCode))
    }

    rand.Seed(time.Now().UnixNano())
    randomKey := rand.Intn(len(bankMerchantCode))
    return fmt.Sprintf("%d", bankMerchantCode[randomKey])
}

const XenditVirtualAccountLength = 18  // max length
const XenditMaxVirtualAccount = 999999 // Xendit offers 999,999 Virtual Account numbers to new customers to start with.

func GenerateVirtualAccount(bankCode string, identifier string) (virtualNumber string, err error) {
    xenditMerchantCode := os.Getenv("XENDIT_MERCHANT_CODE")
    merchantCode := getMerchantCode(bankCode)
    identifier = generateIdentifier(identifier)

    virtualAccountNumber := fmt.Sprintf("%s%s%s", merchantCode, xenditMerchantCode, identifier)
    if len(virtualAccountNumber) > XenditVirtualAccountLength {
        return "", fmt.Errorf("generated virtual account more than %d", XenditVirtualAccountLength)
    }

    return virtualAccountNumber, nil
}

func generateIdentifier(identifier string) string {
    if identifier == "" {
        panic("please give us identifier number")
    }

    maxVirtualAccount := cast.ToInt(os.Getenv("XENDIT_MAX_VIRTUAL_ACCOUNT"))
    if maxVirtualAccount == 0 {
        maxVirtualAccount = XenditMaxVirtualAccount
    }

    lengthIdentifier := len(identifier)
    lengthMax := len(fmt.Sprintf("%d", maxVirtualAccount))

    if lengthMax < lengthIdentifier {
        identifier = identifier[lengthIdentifier-lengthMax : lengthIdentifier]
    } else {
        for i := 0; i < lengthMax-lengthIdentifier; i++ {
            identifier = "0" + identifier
        }
    }

    return identifier
}
