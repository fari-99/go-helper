package storages

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "os"
    "time"
)

//generate path
func (base *StorageBase) generatePath(sType string) (string, string, error) {
    ctime := time.Now().Local()
    datePath := ctime.Format("2006/01/02/")

    var storagePath string
    if base.s3Enabled != nil {
        storagePath = base.s3Enabled.FilePath
    } else if base.gcsEnabled != nil {
        storagePath = base.gcsEnabled.FilePath
    } else {
        storagePath = base.localPath
    }

    if storagePath == "" {
        panic(fmt.Errorf("storage path not set, please set first"))
    }

    filePath := fmt.Sprintf("%s/%s/%s", storagePath, sType, datePath)

    if base.s3Enabled != nil && base.gcsEnabled != nil {
        err := os.MkdirAll(filePath, 0711)
        if err != nil {
            return filePath, datePath, err
        }
    }

    return filePath, datePath, nil
}

// generate filename
func (base *StorageBase) generateName(tmpName string, ext string) string {
    ctime := time.Now().Local()
    filename := tmpName + ctime.Format(time.UnixDate)
    hash := md5.Sum([]byte(filename))

    // encode hash to string
    newName := hex.EncodeToString(hash[:])

    return newName + ext
}
