package token_generator

import (
    _ "crypto"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt"
    "github.com/google/uuid"
)

const accessToken = "access_token"
const refreshToken = "refresh_token"

type BaseJwt struct {
    accessSecret  interface{}
    refreshSecret interface{}
    signingMethod jwt.SigningMethod
    mapClaims     *JwtMapClaims

    expiredAccess  int // in days, default 1
    expiredRefresh int // in days, default 7

    requestCtx *http.Request
}

type JwtMapClaims struct {
    Uuid        string       `json:"uuid"`
    TokenData   TokenData    `json:"token_data"`
    UserDetails *UserDetails `json:"user_details"`
    HasuraClaim HasuraClaim  `json:"hasura_claim"`
    jwt.StandardClaims
}

type SignedToken struct {
    Uuid             string `json:"uuid"`
    AccessToken      string `json:"access_token"`
    RefreshToken     string `json:"refresh_token"`
    AccessExpiredAt  int64  `json:"access_expired_at"`
    RefreshExpiredAt int64  `json:"refresh_expired_at"`
}

func NewJwt(accessSecret, refreshSecret, signedMethod string) *BaseJwt {
    base := BaseJwt{
        accessSecret:   []byte(accessSecret),
        refreshSecret:  []byte(refreshSecret),
        signingMethod:  jwt.GetSigningMethod(signedMethod),
        expiredAccess:  1,
        expiredRefresh: 7,
    }

    return &base
}

func (base *BaseJwt) SetCtx(ctx *http.Request) *BaseJwt {
    base.requestCtx = ctx
    return base
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

func (base *BaseJwt) SetClaim(userDetails UserDetails) (*BaseJwt, error) {
    timeDate := time.Now()

    encryptUserDetails, err := EncryptUserDetails(userDetails)
    if err != nil {
        return nil, err
    }

    claim := JwtMapClaims{
        TokenData: TokenData{
            Origin:      os.Getenv("GO_API_NAME"),
            UserDetails: encryptUserDetails,
            AppData:     base.getAppData(),
        },
        HasuraClaim: HasuraClaim{
            AllowedRoles: []string{"customer", "merchant"},
            DefaultRole:  "customer",
        },
        StandardClaims: jwt.StandardClaims{
            IssuedAt:  timeDate.Unix(),
            ExpiresAt: timeDate.AddDate(0, 0, 1).Unix(), // default 1 day
            Issuer:    os.Getenv("APP_NAME"),
        },
    }

    base.mapClaims = &claim
    return base, nil
}

func (base *BaseJwt) getAppData() *AppData {
    var appData AppData
    if base.requestCtx != nil {
        requestCtx := base.requestCtx
        appData.UserAgent = requestCtx.UserAgent()
        appData.IPList = append(appData.IPList, requestCtx.RemoteAddr)
    }

    appData.AppName = os.Getenv("GO_API_NAME")
    return &appData
}

func (base *BaseJwt) SetClaimApp(appData AppData) *BaseJwt {
    claim := JwtMapClaims{
        TokenData: TokenData{
            Authorized: true,
            AppData:    &appData,
        },
        StandardClaims: jwt.StandardClaims{},
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
func (base *BaseJwt) signClaims(typeClaims string, accessUuid string) (signedToken string, expiredAt int64, err error) {
    expiredDate := base.getExpiredDate(typeClaims)

    mapClaims := base.mapClaims
    mapClaims.Uuid = accessUuid
    mapClaims.StandardClaims.ExpiresAt = expiredDate

    token := jwt.NewWithClaims(base.signingMethod, mapClaims)
    signedToken, err = token.SignedString(base.getSecret(typeClaims))

    return signedToken, expiredDate, err
}

func (base *BaseJwt) getExpiredDate(typeClaims string) int64 {
    timeDate := time.Now()

    var expiredTime time.Time
    switch typeClaims {
    case accessToken:
        expiredTime = timeDate.AddDate(0, 0, base.expiredAccess) // expired for 1 day
    case refreshToken:
        expiredTime = timeDate.AddDate(0, 0, base.expiredRefresh) // expired after 7 day
    }

    return expiredTime.Unix()
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
