package crypts

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "fmt"
    "os"
)

func (base *EncryptionBase) SetSignPrivateKey(key string) *EncryptionBase {
    base.signKey.PrivateKey = key
    return base
}

func (base *EncryptionBase) SetSignPublicKey(key string) *EncryptionBase {
    base.signKey.PublicKey = key
    return base
}

func (base *EncryptionBase) SignData(message string) (signature string, err error) {
    hashedMessage := base.createHash([]byte(message))

    privateKeyBase := base.signKey.PrivateKey
    if privateKeyBase == "" {
        privateKeyBase = os.Getenv("PRIVATE_KEY_SIGN_ENCRYPT")
    }

    privateKeyMarshal, err := base64.RawURLEncoding.DecodeString(privateKeyBase)
    if err != nil {
        return "", fmt.Errorf("error decode base64 rsa private key, err := %s", err.Error())
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyMarshal)
    if err != nil {
        return "", fmt.Errorf("error parse rsa private key, err := %s", err.Error())
    }

    signatureResult, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, []byte(hashedMessage), nil)
    if err != nil {
        return "", fmt.Errorf("failed to sign the message, err := %s", err.Error())
    }

    result := signatureResult
    if base.encodeUrlBase64 {
        result = []byte(base64.RawURLEncoding.EncodeToString(signatureResult))
    }

    return string(result), nil
}

func (base *EncryptionBase) VerifyData(message string, signature []byte) (isVerified bool, err error) {
    hashedMessage := base.createHash([]byte(message))

    publicKeyBase := base.signKey.PublicKey
    if publicKeyBase == "" {
        publicKeyBase = os.Getenv("PUBLIC_KEY_SIGN_ENCRYPT")
    }

    publicKeyMarshal, err := base64.RawURLEncoding.DecodeString(publicKeyBase)
    if err != nil {
        return false, fmt.Errorf("error decode base64 rsa private key, err := %s", err.Error())
    }

    publicKey, err := x509.ParsePKCS1PublicKey(publicKeyMarshal)
    if err != nil {
        return false, fmt.Errorf("error parse rsa private key, err := %s", err.Error())
    }

    if base.encodeUrlBase64 {
        signature, err = base64.RawURLEncoding.DecodeString(string(signature))
        if err != nil {
            return false, fmt.Errorf("error decode signatue, err := %s", err.Error())
        }
    }

    err = rsa.VerifyPSS(publicKey, crypto.SHA256, []byte(hashedMessage), signature, nil)
    if err != nil {
        return false, nil
    }

    return true, nil
}
