package storages

import (
    "fmt"
    "image"
    "image/gif"
    "image/jpeg"
    "image/png"
    "mime/multipart"
    "net/http"
    "os"
    "path"
    "strings"

    gohelper "github.com/fari-99/go-helper"
)

type FileData struct {
    IsImage     bool
    Extension   string
    ImageFile   image.Image
    ContentType string
    StoragePath string
    Filename    string
}

type StorageBase struct {
    fileInput *multipart.FileHeader
    fileType  string
    localPath string

    // NonScaledTypes
    // if images not in this type, then file will be scaled
    // if images in this type, then file will not be scaled
    NonScaledTypes string

    s3Enabled  *S3Setup
    gcsEnabled *GCSSetup
}

type StorageData struct {
    Type             string `json:"type"`
    Path             string `json:"path"`
    Filename         string `json:"filename"`
    Mime             string `json:"mime"`
    OriginalFilename string `json:"original_filename"`
}

type S3Setup struct {
    FilePath   string
    BucketName string
    AccessKey  string
    SecretKey  string
    Region     string
}

type GCSSetup struct {
    FilePath       string
    ProjectID      string
    BucketName     string
    Region         string
    TimeOut        string
    CredentialPath string
}

func NewStorageBase(fileHeader *multipart.FileHeader, fileType string) *StorageBase {
    // default use local
    storageBase := &StorageBase{
        fileInput:  fileHeader,
        fileType:   fileType,
        localPath:  os.Getenv("LOCAL_STORAGE_PATH"),
        s3Enabled:  nil,
        gcsEnabled: nil,
    }

    return storageBase
}

func (base *StorageBase) setNonScaledType(nonScaledTypes []string) *StorageBase {
    base.NonScaledTypes = strings.Join(nonScaledTypes, ",")
    return base
}

func (base *StorageBase) SetAwsS3(s3Setup *S3Setup) *StorageBase {
    if s3Setup == nil {
        s3Setup = &S3Setup{
            FilePath:   os.Getenv("S3_STORAGE_PATH"),
            BucketName: os.Getenv("S3_BUCKET"),
            AccessKey:  os.Getenv("S3_ACCESS_KEY"),
            SecretKey:  os.Getenv("S3_SECRET_KEY"),
            Region:     os.Getenv("S3_REGION"),
        }
    }

    base.s3Enabled = s3Setup
    return base
}

func (base *StorageBase) SetGoogleGCS(gcsSetup *GCSSetup) *StorageBase {
    if gcsSetup == nil {
        gcsSetup = &GCSSetup{
            FilePath:       os.Getenv("GCS_LOCATION"),
            ProjectID:      os.Getenv("GCS_PROJECT_ID"),
            BucketName:     os.Getenv("GCS_BUCKET_NAME"),
            Region:         os.Getenv("GCS_REGION"),
            TimeOut:        os.Getenv("GCS_TIMEOUT"),
            CredentialPath: os.Getenv("GCS_CREDENTIAL_PATH"),
        }
    }

    base.gcsEnabled = gcsSetup
    return base
}

func (base *StorageBase) UploadFiles() (storageModel *StorageData, err error) {
    fileHeader := base.fileInput
    fileType := base.fileType

    file, err := fileHeader.Open()
    if err != nil {
        return
    }

    defer file.Close()

    var scaled = 80
    vals := strings.Split(base.NonScaledTypes, ",")

    if inArray, _, _ := gohelper.InArray(vals, fileType); inArray == true {
        scaled = 100
    }

    contentTypeData, err := base.getFileData(fileHeader)
    if err != nil {
        return
    }

    storagePath, datePath, err := base.generatePath(fileType)
    if err != nil {
        return
    }

    // Generate hash
    fileName := base.generateName(fileHeader.Filename, contentTypeData.Extension)

    contentTypeData.StoragePath = storagePath
    contentTypeData.Filename = fileName

    if base.s3Enabled != nil {
        err = base.s3Upload(contentTypeData, scaled, file)
    } else if base.gcsEnabled != nil {
        err = base.gcsUpload(contentTypeData, scaled, file)
    } else {
        err = base.localUpload(contentTypeData, scaled, file)
    }

    if err != nil {
        return nil, err
    }

    storageModel = &StorageData{
        Type:             fileType,
        Path:             datePath,
        Filename:         fileName,
        Mime:             contentTypeData.ContentType,
        OriginalFilename: fileHeader.Filename,
    }

    return storageModel, nil
}

func (base *StorageBase) GetFiles(storageType, storagePath, filename string) (files *os.File, err error) {
    if base.s3Enabled != nil {
        return base.s3GetFile(storageType, storagePath, filename)
    } else if base.gcsEnabled != nil {
        return base.gcsGetFile(storageType, storagePath, filename)
    } else {
        return base.localGetFile(storageType, storagePath, filename)
    }
}

func (base *StorageBase) getFileData(fileHeader *multipart.FileHeader) (contentTypeData FileData, err error) {
    file, err := fileHeader.Open()
    if err != nil {
        return
    }

    defer file.Close()

    buffer := make([]byte, 1024)
    _, err = file.Read(buffer)
    if err != nil {
        err = fmt.Errorf("file could not be read, err := %s", err.Error())
        return
    }

    _, _ = file.Seek(0, 0)
    contentType := http.DetectContentType(buffer)

    var img image.Image
    var isImage = true
    var ext string

    switch contentType {
    case ContentTypePNG:
        img, err = png.Decode(file)
        ext = ".jpg"
    case ContentTypeGIF:
        img, err = gif.Decode(file)
        ext = ".jpg"
    case ContentTypeJPEG:
        img, err = jpeg.Decode(file)
        ext = ".jpg"
    case ContentTypeJPG:
        img, err = jpeg.Decode(file)
        ext = ".jpg"
    default:
        isImage = false
        // Get file extension
        ext = path.Ext(fileHeader.Filename)
    }

    contentTypeData = FileData{
        IsImage:     isImage,
        Extension:   ext,
        ImageFile:   img,
        ContentType: contentType,
    }

    return
}
