package gohelper

import (
    "log"
    "testing"

    "github.com/fari-99/go-helper/crypts"
)

const passphrase = "9Jw0*9kGDLxhGWf5NAUh"
const testSentence = "The quick brown fox jumps over the lazy dog, 0123456789"

func TestDefault(t *testing.T) {
    encryptBase := crypts.NewEncryptionBase()
    encryptBase.SetPassphrase(passphrase)
    encryptResult, err := encryptBase.Encrypt([]byte(testSentence))
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    decryptBase := crypts.NewEncryptionBase()
    decryptBase.SetPassphrase(passphrase)
    decryptResult, err := decryptBase.Decrypt(encryptResult)
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    log.Printf("encrypt result \t= %s", string(encryptResult))
    log.Printf("decrypt result \t= %s", string(decryptResult))

    if string(decryptResult) != testSentence {
        t.Log(err.Error())
        t.Fail()
        return
    }
}

func TestRSA(t *testing.T) {
    generateKey := crypts.NewEncryptionBase()
    rsaKey, err := generateKey.GenerateRSAKey(2048)
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    encryptBase := crypts.NewEncryptionBase()
    encryptBase.SetPassphrase(passphrase)
    encryptBase.SetEncodeBase64(true)             // if you want the message got encoded before encrypted
    encryptBase.SetRsaPublicKey(rsaKey.PublicKey) // encrypt use public key
    encryptResult, err := encryptBase.EncryptRSA([]byte(testSentence))
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    decryptBase := crypts.NewEncryptionBase()
    decryptBase.SetPassphrase(passphrase)
    encryptBase.SetEncodeBase64(true)               // true if message got encoded before encrypted
    decryptBase.SetRsaPrivateKey(rsaKey.PrivateKey) // decrypt using private key
    decryptResult, err := decryptBase.DecryptRSA(encryptResult)
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    log.Printf("encrypt result RSA \t= %s", string(encryptResult))
    log.Printf("decrypt result RSA \t= %s", string(decryptResult))

    if string(decryptResult) != testSentence {
        t.Log(err.Error())
        t.Fail()
        return
    }
}

func TestSign(t *testing.T) {
    generateKey := crypts.NewEncryptionBase()
    rsaKey, err := generateKey.GenerateRSAKey(2048)
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    signBase := crypts.NewEncryptionBase()
    signBase.SetPassphrase(passphrase)
    signBase.SetEncodeBase64(true)                // if you want the message got encoded before encrypted
    signBase.SetSignPrivateKey(rsaKey.PrivateKey) // sign use private key
    signature, err := signBase.SignData(testSentence)
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    verifyBase := crypts.NewEncryptionBase()
    verifyBase.SetPassphrase(passphrase)
    verifyBase.SetEncodeBase64(true)              // if you want the message got encoded before encrypted
    verifyBase.SetSignPublicKey(rsaKey.PublicKey) // verify use public key
    isVerified, err := verifyBase.VerifyData(testSentence, []byte(signature))
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    log.Printf("sentence \t:= %s", testSentence)
    log.Printf("signature \t:= %s", signature)

    if !isVerified {
        t.Log("message signature is invalid")
        t.Fail()
        return
    }
}

func TestEncryptFile(t *testing.T) {
    encryptBase := crypts.NewEncryptionBase()
    encryptBase.SetPassphrase(passphrase)
    err := encryptBase.EncryptFile("./crypts/examples/infile.txt", "./crypts/examples/encrypt.txt")
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }

    decryptBase := crypts.NewEncryptionBase()
    decryptBase.SetPassphrase(passphrase)
    err = decryptBase.DecryptFile("./crypts/examples/encrypt.txt", "./crypts/examples/decrypt.txt")
    if err != nil {
        t.Log(err.Error())
        t.Fail()
        return
    }
}
