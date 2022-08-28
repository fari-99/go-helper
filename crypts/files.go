package crypts

import (
    "bufio"
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "errors"
    "fmt"
    "io"
    "os"
)

const BufferSize int = 4096
const IvSize int = 16
const HMACSize = sha256.Size

func (base *EncryptionBase) checkValidRunEncryptFile() error {
    if string(base.passphrase) == "" {
        return fmt.Errorf("please set passphare")
    }

    return nil
}

func (base *EncryptionBase) EncryptFile(filePathIn, filePathOut string) error {
    if err := base.checkValidRunEncryptFile(); err != nil {
        return err
    }

    inFile, err := os.Open(filePathIn)
    if err != nil {
        return err
    }
    defer inFile.Close()

    // create output files to store encrypted file
    outFile, err := os.Create(filePathOut)
    if err != nil {
        return err
    }
    defer outFile.Close()

    iv := make([]byte, IvSize)
    _, err = rand.Read(iv)
    if err != nil {
        return err
    }

    // get passphrase
    passphrase := []byte(base.createHash(base.passphrase))

    // generate new chiper
    chiperBlock, err := aes.NewCipher(passphrase)
    if err != nil {
        return err
    }

    // generate hmac from passphrase
    hmacHash := hmac.New(sha256.New, passphrase)

    // write iv to output file
    _, _ = outFile.Write(iv)

    // write iv to hmac
    hmacHash.Write(iv)

    // using CTR as base
    ctr := cipher.NewCTR(chiperBlock, iv)

    // create buffer
    buf := make([]byte, BufferSize)
    for {
        n, err := inFile.Read(buf)
        if err != nil && err != io.EOF {
            return err
        }

        // check if it's an end of file
        if err == io.EOF {
            break
        }

        // generate new output buffer
        outBuf := make([]byte, n)

        // encrypt buffer stream to output buffer
        ctr.XORKeyStream(outBuf, buf[:n])

        // write output buffer to hmac
        hmacHash.Write(outBuf)

        // write new buffer to output file
        _, _ = outFile.Write(outBuf)
    }

    // write hmac to output file
    _, _ = outFile.Write(hmacHash.Sum(nil))

    return nil
}

func (base *EncryptionBase) DecryptFile(filePathIn, filePathOut string) error {
    if err := base.checkValidRunEncryptFile(); err != nil {
        return err
    }
    
    inFile, err := os.Open(filePathIn)
    if err != nil {
        return err
    }
    defer inFile.Close()

    // create output files to store encrypted file
    outFile, err := os.Create(filePathOut)
    if err != nil {
        return err
    }
    defer outFile.Close()

    iv := make([]byte, IvSize)
    _, err = io.ReadFull(inFile, iv)
    if err != nil {
        return err
    }

    passphrase := []byte(base.createHash(base.passphrase))

    AES, err := aes.NewCipher(passphrase)
    if err != nil {
        return err
    }

    ctr := cipher.NewCTR(AES, iv)
    h := hmac.New(sha256.New, passphrase)
    h.Write(iv)
    mac := make([]byte, HMACSize)

    w := outFile

    buf := bufio.NewReaderSize(inFile, BufferSize)
    var limit int
    var b []byte
    for {
        b, err = buf.Peek(BufferSize)
        if err != nil && err != io.EOF {
            return err
        }

        limit = len(b) - HMACSize

        // We reached the end
        if err == io.EOF {

            left := buf.Buffered()
            if left < HMACSize {
                return errors.New("not enough left")
            }

            copy(mac, b[left-HMACSize:left])

            if left == HMACSize {
                break
            }
        }

        h.Write(b[:limit])

        // We always leave at least hmacSize bytes left in the buffer
        // That way, our next Peek() might be EOF, but we will still have enough
        outBuf := make([]byte, int64(limit))
        _, err = buf.Read(b[:limit])
        if err != nil {
            return err
        }

        ctr.XORKeyStream(outBuf, b[:limit])
        _, err = w.Write(outBuf)
        if err != nil {
            return err
        }

        if err == io.EOF {
            break
        }
    }

    if !hmac.Equal(mac, h.Sum(nil)) {
        return fmt.Errorf("invalid hmac")
    }

    return nil
}
