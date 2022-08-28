package crypts

import (
    "crypto/sha256"
)

type EncryptionBase struct {
    passphrase      []byte
    encodeUrlBase64 bool

    useRandom   bool
    randomSetup []byte

    rsaPrivateKey string // encode to base64 raw url encoding
    rsaPublicKey  string // encode to base64 raw url encoding
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
