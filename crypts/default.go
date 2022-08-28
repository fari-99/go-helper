package crypts

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
)

func (base *EncryptionBase) Encrypt(data []byte) ([]byte, error) {
    blockKey, err := aes.NewCipher([]byte(base.createHash(base.passphrase)))
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create gcm, err := %s", err.Error())
    }

    var randomness io.Reader
    if base.useRandom {
        randomness = rand.Reader
    } else {
        randomness = bytes.NewReader(base.randomSetup)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(randomness, nonce); err != nil {
        return nil, fmt.Errorf("failed to create nonce, err := %s", err.Error())
    }

    cipherText := gcm.Seal(nonce, nonce, data, nil)

    result := cipherText
    if base.encodeUrlBase64 {
        result = []byte(base64.RawURLEncoding.EncodeToString(cipherText))
    }

    return result, nil
}

func (base *EncryptionBase) Decrypt(data []byte) ([]byte, error) {
    blockKey, err := aes.NewCipher([]byte(base.createHash(base.passphrase)))
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockKey)
    if err != nil {
        return nil, err
    }

    if base.encodeUrlBase64 {
        data, err = base64.RawURLEncoding.DecodeString(string(data))
        if err != nil {
            return nil, err
        }
    }

    nonceSize := gcm.NonceSize()
    nonce, cipherText := data[:nonceSize], data[nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}
