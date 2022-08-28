package crypts

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func (base *EncryptionBase) SetRsaPrivateKey(key string) *EncryptionBase {
    base.rsaKey.PrivateKey = key
    return base
}

func (base *EncryptionBase) SetRsaPublicKey(key string) *EncryptionBase {
    base.rsaKey.PublicKey = key
    return base
}

func (base *EncryptionBase) EncryptRSA(secretMessage []byte) ([]byte, error) {
    passphrase := base.passphrase
    encodeBase64 := base.encodeUrlBase64

    publicKeyBase := base.rsaKey.PublicKey
    if publicKeyBase == "" {
        publicKeyBase = os.Getenv("PUBLIC_KEY_RSA_ENCRYPT")
    }

    publicKeyMarshal, err := base64.RawURLEncoding.DecodeString(publicKeyBase)
    if err != nil {
        return nil, fmt.Errorf("error decode base64 rsa public key, err := %s", err.Error())
    }

    publicKey, err := x509.ParsePKCS1PublicKey(publicKeyMarshal)
    if err != nil {
        return nil, fmt.Errorf("error parse rsa public key, err := %s", err.Error())
    }

    var randomness io.Reader
    if base.useRandom {
        randomness = rand.Reader
    } else {
        randomness = bytes.NewReader(base.randomSetup)
    }

    cipherText, err := rsa.EncryptOAEP(sha256.New(), randomness, publicKey, secretMessage, passphrase)

    result := cipherText
    if encodeBase64 {
        result = []byte(base64.RawURLEncoding.EncodeToString(cipherText))
    }

    return result, err
}

func (base *EncryptionBase) DecryptRSA(secretMessage []byte) ([]byte, error) {
    passphrase := base.passphrase
    encodeBase64 := base.encodeUrlBase64

    privateKeyBase := base.rsaKey.PrivateKey
    if privateKeyBase == "" {
        privateKeyBase = os.Getenv("PRIVATE_KEY_RSA_ENCRYPT")
    }

    privateKeyMarshal, err := base64.RawURLEncoding.DecodeString(privateKeyBase)
    if err != nil {
        return nil, fmt.Errorf("error decode base64 rsa private key, err := %s", err.Error())
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyMarshal)
    if err != nil {
        return nil, fmt.Errorf("error parse rsa private key, err := %s", err.Error())
    }

    var randomness io.Reader
    if base.useRandom {
        randomness = rand.Reader
    } else {
        randomness = bytes.NewReader(base.randomSetup)
    }

    var message []byte
    if encodeBase64 {
        message, err = base64.RawURLEncoding.DecodeString(string(secretMessage))
        if err != nil {
            return nil, fmt.Errorf("message is not base64 encoded, err := %s", err.Error())
        }
    }

    result, err := rsa.DecryptOAEP(sha256.New(), randomness, privateKey, message, passphrase)
    return result, err
}
