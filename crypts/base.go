package crypts

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "fmt"
)

type EncryptionBase struct {
    passphrase      []byte
    encodeUrlBase64 bool

    useRandom   bool
    randomSetup []byte

    rsaKey  Keys
    signKey Keys
}

type Keys struct {
    PrivateKey string // encode to base64 raw url encoding
    PublicKey  string // encode to base64 raw url encoding
}

func NewEncryptionBase() *EncryptionBase {
    base := &EncryptionBase{
        encodeUrlBase64: true,
        useRandom:       true,
    }

    return base
}

func (base *EncryptionBase) SetPassphrase(passphrase string) *EncryptionBase {
    if passphrase != "" {
        base.passphrase = []byte(passphrase)
    }

    return base
}

func (base *EncryptionBase) GenerateRSAKey(bitSize int) (rsaKey *Keys, err error) {
    if bitSize <= 0 {
        bitSize = 2048
    }

    randomness := rand.Reader
    key, err := rsa.GenerateKey(randomness, bitSize)
    if err != nil {
        return nil, fmt.Errorf("error generate RSA key, err := %s", err.Error())
    }

    keyPrivate := x509.MarshalPKCS1PrivateKey(key)
    keyPublic := x509.MarshalPKCS1PublicKey(&key.PublicKey)

    // encode to base 64
    privateKey := base64.RawURLEncoding.EncodeToString(keyPrivate)
    publicKey := base64.RawURLEncoding.EncodeToString(keyPublic)

    base.rsaKey = Keys{
        PrivateKey: privateKey,
        PublicKey:  publicKey,
    }

    return &base.rsaKey, nil
}

func (base *EncryptionBase) SetUseRandomness(useRandomness bool, keyRandom string) *EncryptionBase {
    base.useRandom = useRandomness

    if !useRandomness {
        if keyRandom != "" {
            base.randomSetup = []byte(keyRandom)
        } else {
            // randomness used for security data, if you are not using it. it still need some key random,
            // so you must set it.
            panic("need to set random key when use random is false.")
        }
    }

    return base
}

func (base *EncryptionBase) SetEncodeBase64(useEncode64 bool) *EncryptionBase {
    base.encodeUrlBase64 = useEncode64
    return base
}

func (base *EncryptionBase) createHash(passphrase []byte) string {
    hash := sha256.New()
    hash.Write(passphrase)
    return string(hash.Sum(nil))
}
