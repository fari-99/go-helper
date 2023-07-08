package virtual_accounts

import (
    "encoding/json"
    "fmt"
    "log"
    "testing"
    "time"

    "github.com/go-resty/resty/v2"
    "github.com/google/uuid"
    "github.com/xendit/xendit-go"
    "github.com/xendit/xendit-go/virtualaccount"
)

func TestXenditGetAvailableBanks(t *testing.T) {
    xendit.Opt.SecretKey = ""

    resp, errXendit := virtualaccount.GetAvailableBanks()
    if errXendit != nil {
        errMarshal, _ := json.Marshal(errXendit)
        t.Error(string(errMarshal))
        t.Fail()
        return
    }

    log.Printf("Success Get Available Banks")
    responseMarshal, _ := json.MarshalIndent(resp, " ", "")
    log.Printf(string(responseMarshal))

    /*
       [
           {
               "name": "Bank Central Asia",
               "code": "BCA",
               "country": "ID",
               "currency": "IDR",
               "is_activated": true
           },
           ...
       ]
    */
}

func TestXenditGenerateVirtualAccount(t *testing.T) {
    //_ = os.Setenv("XENDIT_MERCHANT_CODE", "9999")
    //virtualAccountNumber, err := GenerateVirtualAccount(BankBNI, "873618")
    //if err != nil {
    //    t.Log(err.Error())
    //    t.Fail()
    //    return
    //}

    // (is_closed: true) Closed amount means your customer can only pay amount specified by you.
    // Payment will be rejected if attempted payment amount deviates from the amount you specified.
    // Specify the amount using expected_amount parameter
    // (is_closed: false) Open amount means your customer can pay any amount to the Virtual Account
    isClosed := true

    // (is_single_use: true) can only be paid once.
    // Used VA number can be recreated for other customer/invoice/transaction
    // (is_single_use: false) allows your customer to pay to the same Virtual Account continuously
    isSingleUse := false

    expiredAt := time.Now().AddDate(0, 0, 1)
    identifierTransactions := uuid.New().String()

    xendit.Opt.SecretKey = ""
    data := virtualaccount.CreateFixedVAParams{
        //ForUserID:            "",
        ExternalID: identifierTransactions,
        BankCode:   BankBNC,
        Name:       "Michael Jackson",
        //VirtualAccountNumber: "", // better generated by xendit
        IsClosed:       &isClosed,
        IsSingleUse:    &isSingleUse,
        ExpirationDate: &expiredAt,
        //SuggestedAmount:      0, // some bank don't use suggested amount
        ExpectedAmount: 100000,
        //Description:          "", // some bank don't use descriptions
    }

    log.Printf("identifier transactions: %s", identifierTransactions)
    //log.Printf("virtual account number: %s", virtualAccountNumber)

    resp, errXendit := virtualaccount.CreateFixedVA(&data)
    if errXendit != nil {
        errMarshal, _ := json.Marshal(errXendit)
        t.Error(string(errMarshal))
        t.Fail()
        return
    }

    log.Printf("Success Create Virtual Account")
    responseMarshal, _ := json.MarshalIndent(resp, " ", "")
    log.Printf(string(responseMarshal))

    /*
       {
           "owner_id": "64a75cb6cd2807c03d11a6c6",
           "external_id": "20421c52-8b29-454f-b879-e96edcbf4d82",
           "bank_code": "BNC",
           "merchant_code": "90100010",
           "name": "Fadhlan Punya",
           "account_number": "901000109999766758",
           "is_closed": true,
           "id": "eb2d826c-30fc-41e1-bb02-80bcbdc5ef3d",
           "is_single_use": false,
           "status": "PENDING",
           "currency": "IDR",
           "expiration_date": "2023-07-09T03:31:26.486Z",
           "expected_amount": 100000
       }
    */

}

func TestXenditGetVirtualAccount(t *testing.T) {
    xendit.Opt.SecretKey = ""

    data := virtualaccount.GetFixedVAParams{
        ID: "eb2d826c-30fc-41e1-bb02-80bcbdc5ef3d",
    }

    resp, errXendit := virtualaccount.GetFixedVA(&data)
    if errXendit != nil {
        errMarshal, _ := json.Marshal(errXendit)
        t.Error(string(errMarshal))
        t.Fail()
        return
    }

    log.Printf("Success Get Virtual Account")
    responseMarshal, _ := json.MarshalIndent(resp, " ", "")
    log.Printf(string(responseMarshal))
}

func TestXenditUpdateVirtualAccount(t *testing.T) {
    xendit.Opt.SecretKey = ""

    expirationDate := time.Now().Add(5 * time.Minute)

    updateFixedVAData := virtualaccount.UpdateFixedVAParams{
        ID: "20421c52-8b29-454f-b879-e96edcbf4d82",
        //IsSingleUse:     nil,
        ExpirationDate:  &expirationDate,
        SuggestedAmount: 0,
        //SuggestedAmount: 0,
        ExpectedAmount: 20000,
        //Description:    "",
    }

    resp, errXendit := virtualaccount.UpdateFixedVA(&updateFixedVAData)
    if errXendit != nil {
        errMarshal, _ := json.Marshal(errXendit)
        t.Error(string(errMarshal))
        t.Fail()
        return
    }

    log.Printf("Success Update Virtual Account")
    responseMarshal, _ := json.MarshalIndent(resp, " ", "")
    log.Printf(string(responseMarshal))
}

func TestXenditSimulatePayment(t *testing.T) {
    secretKey := ""
    client := resty.New()

    externalID := "20421c52-8b29-454f-b879-e96edcbf4d82"
    urlSimulate := fmt.Sprintf("https://api.xendit.co/callback_virtual_accounts/external_id=%s/simulate_payment", externalID)
    resp, err := client.R().
        SetBasicAuth(secretKey, "").
        SetBody(map[string]interface{}{"amount": 20000}).
        Post(urlSimulate)
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }

    log.Printf("Success Simulate Payment Virtual Account")
    log.Printf(string(resp.Body()))

    /*
       {
           "status":"COMPLETED",
           "message":"Payment for the Fixed VA with external id 20421c52-8b29-454f-b879-e96edcbf4d82 is currently being processed.
            Please ensure that you have set a callback URL for VA payments via Dashboard Settings and contact us
            if you do not receive a VA payment callback within the next 5 mins."
       }
    */
}

func TestXenditGetPayments(t *testing.T) {
    xendit.Opt.SecretKey = ""

    getPayment := virtualaccount.GetPaymentParams{
        PaymentID: "ae2c6532-da11-484e-9d83-6dcae9ccc67d",
    }

    resp, errXendit := virtualaccount.GetPayment(&getPayment)
    if errXendit != nil {
        errMarshal, _ := json.Marshal(errXendit)
        t.Error(string(errMarshal))
        t.Fail()
        return
    }

    log.Printf("Success Get Virtual Account Payments")
    responseMarshal, _ := json.MarshalIndent(resp, " ", "")
    log.Printf(string(responseMarshal))
    /*
       {
            "id": "ae2c6532-da11-484e-9d83-6dcae9ccc67d",
            "payment_id": "ae2c6532-da11-484e-9d83-6dcae9ccc67d",
            "callback_virtual_account_id": "eb2d826c-30fc-41e1-bb02-80bcbdc5ef3d",
            "external_id": "20421c52-8b29-454f-b879-e96edcbf4d82",
            "account_number": "9999766758",
            "bank_code": "BNC",
            "amount": 20000,
            "transaction_timestamp": "2023-07-08T03:56:26Z",
            "merchant_code": "90100010",
            "currency": "IDR",
            "sender_name": ""
        }
    */
}
