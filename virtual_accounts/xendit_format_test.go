package virtual_accounts

import (
    "log"
    "os"
    "testing"
)

func TestGenerateVirtualAccount(t *testing.T) {
    os.Setenv("XENDIT_MERCHANT_CODE", "1010")
    virtualAccount, err := GenerateVirtualAccount(BankBCA, "999999")
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    log.Printf(virtualAccount)
}
