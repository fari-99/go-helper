package gohelper

import (
    "testing"

    "github.com/fari-99/go-helper/token_generator"
)

// for security, use more than 43 character or 256 bit
// you can use password generator website to generate this
const AccessToken = "qHtiap9l1#bX6^T0SE71@6tMBLrv%ntlbUiyiFweOpo"
const RefreshToken = "gG7oNYp8U@y7o3A0uPdr4cAR2G7OiPqFfdF@d3OLKdQ"
const SignMethod = "HS512"

var userDetails = token_generator.UserDetails{
    ID:        "1",
    Email:     "test@gmail.com",
    Username:  "test",
    UserRoles: []string{"admin", "pic", "finance"},
}

func TestToken(t *testing.T) {
    createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
    createTokenHelper, err := createTokenHelper.SetClaim(userDetails)
    if err != nil {
        t.Log("error set jwt claim")
        t.Log(err.Error())
        t.Fail()
        return
    }

    tokens, err := createTokenHelper.SignClaims()
    if err != nil {
        t.Log("error sign jwt claim")
        t.Log(err.Error())
        t.Fail()
        return
    }

    checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
    _, err = checkTokenHelper.ParseToken("access_token", tokens.AccessToken)
    if err != nil {
        t.Log("failed to parse access token")
        t.Log(err.Error())
        t.Fail()
        return
    }

    _, err = checkTokenHelper.ParseToken("refresh_token", tokens.RefreshToken)
    if err != nil {
        t.Log("failed to parse refresh token")
        t.Log(err.Error())
        t.Fail()
        return
    }

    //log.Printf("secret token := %s", tokens.AccessToken)
    //log.Printf("refresh token := %s", tokens.RefreshToken)
    //log.Printf("expired secret token := %s", time.Unix(tokens.AccessExpiredAt, 0).Format(formatDate))
    //log.Printf("expired refresh token := %s", time.Unix(tokens.RefreshExpiredAt, 0).Format(formatDate))
    //log.Printf("uuid := %s", tokens.Uuid)
}
