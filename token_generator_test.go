package gohelper

import (
	"net/http"
	"testing"
	"time"

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
	TwoFAModels: token_generator.TwoFAModels{
		TOTP:         true,
		RecoveryCode: false,
		Email:        true,
	},
}

func TestToken(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetIssuer("current_app_name")
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
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
	accessClaims, err := checkTokenHelper.ParseToken("access_token", tokens.AccessToken)
	if err != nil {
		t.Log("failed to parse access token")
		t.Log(err.Error())
		t.Fail()
		return
	}

	refreshClaims, err := checkTokenHelper.ParseToken("refresh_token", tokens.RefreshToken)
	if err != nil {
		t.Log("failed to parse refresh token")
		t.Log(err.Error())
		t.Fail()
		return
	}

	if accessClaims.Uuid != tokens.Uuid || refreshClaims.Uuid != tokens.Uuid {
		t.Errorf("expected both claims to carry uuid %q, got access=%q refresh=%q",
			tokens.Uuid, accessClaims.Uuid, refreshClaims.Uuid)
	}

	if accessClaims.UserDetails == nil || accessClaims.UserDetails.Email != userDetails.Email {
		t.Errorf("expected decrypted user details with email %q", userDetails.Email)
	}

	if accessClaims.TwoFAModels != userDetails.TwoFAModels {
		t.Errorf("expected TwoFAModels %+v, got %+v", userDetails.TwoFAModels, accessClaims.TwoFAModels)
	}

	//log.Printf("secret token := %s", tokens.AccessToken)
	//log.Printf("refresh token := %s", tokens.RefreshToken)
	//log.Printf("expired secret token := %s", time.Unix(tokens.AccessExpiredAt, 0).Format(formatDate))
	//log.Printf("expired refresh token := %s", time.Unix(tokens.RefreshExpiredAt, 0).Format(formatDate))
	//log.Printf("uuid := %s", tokens.Uuid)
}

func TestToken_SetRoles(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
	createTokenHelper.SetRoles([]string{"editor", "viewer"}, "editor")

	createTokenHelper, err := createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	claims, err := checkTokenHelper.ParseToken("access_token", tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}

	if claims.HasuraClaim.DefaultRole != "editor" {
		t.Errorf("expected default role editor, got %q", claims.HasuraClaim.DefaultRole)
	}

	wantRoles := map[string]bool{"admin": true, "editor": true, "viewer": true}
	if len(claims.HasuraClaim.AllowedRoles) != len(wantRoles) {
		t.Fatalf("expected %d allowed roles, got %v", len(wantRoles), claims.HasuraClaim.AllowedRoles)
	}
	for _, role := range claims.HasuraClaim.AllowedRoles {
		if !wantRoles[role] {
			t.Errorf("unexpected role in allowed roles: %q", role)
		}
	}
}

func TestToken_ClientIPAndUserAgent(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("User-Agent", "unit-test-agent")

	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
	createTokenHelper.SetCtx(req)
	createTokenHelper.SetClientIP("203.0.113.10")

	createTokenHelper, err = createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	claims, err := checkTokenHelper.ParseToken("access_token", tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}

	appData := claims.TokenData.AppData
	if appData == nil {
		t.Fatal("expected app data on parsed claims")
	}
	if appData.AppName != "app_name_ask_for_auth" {
		t.Errorf("expected AppName app_name_ask_for_auth, got %q", appData.AppName)
	}
	if appData.UserAgent != "unit-test-agent" {
		t.Errorf("expected UserAgent unit-test-agent, got %q", appData.UserAgent)
	}
	if len(appData.IPList) != 1 || appData.IPList[0] != "203.0.113.10" {
		t.Errorf("expected IPList [203.0.113.10], got %v", appData.IPList)
	}
}

func TestToken_ParseToken_InvalidToken(t *testing.T) {
	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	if _, err := checkTokenHelper.ParseToken("access_token", "not-a-valid-jwt"); err == nil {
		t.Error("expected error parsing malformed token")
	}
}

func TestToken_ParseToken_WrongSecret(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
	createTokenHelper, err := createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	wrongSecretHelper := token_generator.NewJwt("a-completely-different-access-secret", RefreshToken, SignMethod)
	if _, err := wrongSecretHelper.ParseToken("access_token", tokens.AccessToken); err == nil {
		t.Error("expected error parsing token signed with a different secret")
	}
}

func TestToken_ParseToken_CrossTokenTypeRejected(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
	createTokenHelper, err := createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	if _, err := checkTokenHelper.ParseToken("refresh_token", tokens.AccessToken); err == nil {
		t.Error("expected error parsing access token as a refresh token")
	}
	if _, err := checkTokenHelper.ParseToken("access_token", tokens.RefreshToken); err == nil {
		t.Error("expected error parsing refresh token as an access token")
	}
}

func TestToken_ParseToken_TamperedTokenRejected(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")
	createTokenHelper, err := createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	tampered := tokens.AccessToken[:len(tokens.AccessToken)-2] + "xx"

	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	if _, err := checkTokenHelper.ParseToken("access_token", tampered); err == nil {
		t.Error("expected error parsing tampered token")
	}
}

func TestToken_SetExpired_ValidationError(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	if _, err := createTokenHelper.SetExpired(2, time.Hour, 1, time.Hour); err == nil {
		t.Error("expected error when refresh expiry is shorter than access expiry")
	}
}

func TestToken_SetExpired_AffectsSignedTokenExpiry(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	createTokenHelper.SetOrigin("app_name_ask_for_auth")

	createTokenHelper, err := createTokenHelper.SetExpired(1, time.Hour, 6, time.Hour)
	if err != nil {
		t.Fatalf("SetExpired returned error: %v", err)
	}

	createTokenHelper, err = createTokenHelper.SetClaim(userDetails)
	if err != nil {
		t.Fatalf("error set jwt claim: %v", err)
	}

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	if !tokens.RefreshExpiredAt.After(tokens.AccessExpiredAt) {
		t.Errorf("expected RefreshExpiredAt (%s) to be after AccessExpiredAt (%s)",
			tokens.RefreshExpiredAt, tokens.AccessExpiredAt)
	}

	gap := tokens.RefreshExpiredAt.Sub(tokens.AccessExpiredAt)
	if gap < 4*time.Hour || gap > 6*time.Hour {
		t.Errorf("expected roughly a 5h gap between refresh and access expiry, got %s", gap)
	}
}

func TestToken_SetClaimApp(t *testing.T) {
	createTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	appData := token_generator.AppData{AppName: "service-a", UserAgent: "svc"}
	createTokenHelper = createTokenHelper.SetClaimApp(appData)

	tokens, err := createTokenHelper.SignClaims()
	if err != nil {
		t.Fatalf("error sign jwt claim: %v", err)
	}

	checkTokenHelper := token_generator.NewJwt(AccessToken, RefreshToken, SignMethod)
	claims, err := checkTokenHelper.ParseToken("access_token", tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}

	if !claims.TokenData.Authorized {
		t.Error("expected TokenData.Authorized to be true")
	}
	if claims.TokenData.AppData == nil || claims.TokenData.AppData.AppName != "service-a" {
		t.Errorf("expected AppName service-a, got %+v", claims.TokenData.AppData)
	}
}
