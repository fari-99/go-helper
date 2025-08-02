package token_generator

import (
	_ "crypto"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const accessToken = "access_token"
const refreshToken = "refresh_token"
const adminRole = "admin"

type BaseJwt struct {
	accessSecret  interface{}
	refreshSecret interface{}
	signingMethod jwt.SigningMethod
	mapClaims     *JwtMapClaims

	expiredAccess  int // in days, default 1
	expiredRefresh int // in days, default 7

	origin string
	issuer string

	requestCtx *http.Request

	allowedRoles []string
	defaultRoles string
}

type JwtMapClaims struct {
	Uuid        string       `json:"uuid"`
	TokenData   TokenData    `json:"token_data"`
	UserDetails *UserDetails `json:"user_details"`
	HasuraClaim HasuraClaim  `json:"hasura_claim"`
	jwt.RegisteredClaims
}

type SignedToken struct {
	Uuid             string    `json:"uuid"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

func NewJwt(accessSecret, refreshSecret, signedMethod string) *BaseJwt {
	base := BaseJwt{
		accessSecret:   []byte(accessSecret),
		refreshSecret:  []byte(refreshSecret),
		signingMethod:  jwt.GetSigningMethod(signedMethod),
		expiredAccess:  1,
		expiredRefresh: 7,
		defaultRoles:   adminRole,
		allowedRoles:   []string{adminRole},
	}

	return &base
}

func (base *BaseJwt) SetCtx(ctx *http.Request) *BaseJwt {
	base.requestCtx = ctx
	return base
}

func (base *BaseJwt) SetOrigin(origin string) *BaseJwt {
	base.origin = origin
	return base
}

func (base *BaseJwt) SetIssuer(issuer string) *BaseJwt {
	base.issuer = issuer
	return base
}

func (base *BaseJwt) SetRoles(allowedRoles []string, defaultRoles string) *BaseJwt {
	for _, allowedRole := range allowedRoles {
		base.allowedRoles = append(base.allowedRoles, allowedRole)
	}

	if defaultRoles != "" {
		base.defaultRoles = defaultRoles
	}

	return base
}

func (base *BaseJwt) SetClaim(userDetails UserDetails) (*BaseJwt, error) {
	timeDate := time.Now()

	encryptUserDetails, err := EncryptUserDetails(userDetails)
	if err != nil {
		return nil, err
	}

	claim := JwtMapClaims{
		TokenData: TokenData{
			Origin:      base.origin,
			UserDetails: encryptUserDetails,
			AppData:     base.getAppData(),
		},
		HasuraClaim: HasuraClaim{
			AllowedRoles: base.allowedRoles,
			DefaultRole:  base.defaultRoles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeDate),
			ExpiresAt: jwt.NewNumericDate(timeDate.AddDate(0, 0, 1)), // default 1 day
			Issuer:    base.issuer,
		},
	}

	base.mapClaims = &claim
	return base, nil
}

func (base *BaseJwt) SetClaimApp(appData AppData) *BaseJwt {
	claim := JwtMapClaims{
		TokenData: TokenData{
			Authorized: true,
			AppData:    &appData,
		},
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	base.mapClaims = &claim
	return base
}

func (base *BaseJwt) SignClaims() (signedToken *SignedToken, err error) {
	accessUuid := uuid.New().String()
	accessSignedToken, accessExpired, err := base.signClaims(accessToken, accessUuid)
	if err != nil {
		return nil, err
	}

	refreshSignedToken, refreshExpired, err := base.signClaims(refreshToken, accessUuid)
	if err != nil {
		return nil, err
	}

	token := &SignedToken{
		Uuid:             accessUuid,
		AccessToken:      accessSignedToken,
		AccessExpiredAt:  accessExpired,
		RefreshToken:     refreshSignedToken,
		RefreshExpiredAt: refreshExpired,
	}

	return token, nil
}

func (base *BaseJwt) ParseToken(typeClaims, jwtToken string) (*JwtMapClaims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &JwtMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		secret := base.getSecret(typeClaims)
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtMapClaims); ok && token.Valid {
		if claims.TokenData.UserDetails != "" {
			userDetails, err := DecryptUserDetails(claims.TokenData.UserDetails)
			if err != nil {
				return nil, err
			}

			claims.UserDetails = &userDetails
		}

		return claims, nil
	}

	return nil, err
}

func (base *BaseJwt) setExpired(expiredAccess, expiredRefresh int) (*BaseJwt, error) {
	// no expired on JWT token is a bad practice
	if expiredAccess <= 0 { // if set 0, then default is 1
		expiredAccess = base.expiredAccess
	}

	if expiredRefresh == 0 { // if set 0, then default is 7
		expiredRefresh = base.expiredRefresh
	}

	if expiredAccess < expiredRefresh {
		return nil, fmt.Errorf("refresh expired should be more or equal than access expired")
	}

	base.expiredAccess = expiredRefresh
	base.expiredRefresh = expiredRefresh
	return base, nil
}

func (base *BaseJwt) getAppData() *AppData {
	var appData AppData
	if base.requestCtx != nil {
		requestCtx := base.requestCtx
		appData.UserAgent = requestCtx.UserAgent()
		appData.IPList = append(appData.IPList, requestCtx.RemoteAddr)
	}

	appData.AppName = base.origin
	return &appData
}

func (base *BaseJwt) getSecret(typeClaims string) interface{} {
	var secret interface{}
	if typeClaims == accessToken {
		secret = base.accessSecret
	} else if typeClaims == refreshToken {
		secret = base.refreshSecret
	}

	return secret
}

/**
All JWT uuid, must be signed with ACCESS UUID from JWT UUID
so that we can made only one refresh token for many device,
but still have one access token
*/
func (base *BaseJwt) signClaims(typeClaims string, accessUuid string) (signedToken string, expiredAt time.Time, err error) {
	expiredDate := base.getExpiredDate(typeClaims)

	mapClaims := base.mapClaims
	mapClaims.Uuid = accessUuid
	mapClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiredDate)

	token := jwt.NewWithClaims(base.signingMethod, mapClaims)
	signedToken, err = token.SignedString(base.getSecret(typeClaims))

	return signedToken, expiredDate, err
}

func (base *BaseJwt) getExpiredDate(typeClaims string) time.Time {
	timeDate := time.Now()

	var expiredTime time.Time
	switch typeClaims {
	case accessToken:
		expiredTime = timeDate.AddDate(0, 0, base.expiredAccess) // expired for 1 day
	case refreshToken:
		expiredTime = timeDate.AddDate(0, 0, base.expiredRefresh) // expired after 7 day
	}

	return expiredTime
}
