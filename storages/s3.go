package storages

import (
    "fmt"
    "image/jpeg"
    "io"
    "io/ioutil"
    "mime/multipart"
    "os"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/fari-99/aws-presignpost-s3-go"
)

func (base *StorageBase) s3Session() (sessionConfig *session.Session, err error) {
    awsAccessKey := aws.String(base.s3Enabled.AccessKey)
    awsSecretKey := aws.String(base.s3Enabled.SecretKey)
    token := ""

    credential := credentials.NewStaticCredentials(*awsAccessKey, *awsSecretKey, token)
    _, err = credential.Get()
    if err != nil {
        err = fmt.Errorf("bad AWS credentials, err := %s", err.Error())
        return
    }

    cfg := aws.NewConfig().
        WithRegion(base.s3Enabled.Region).
        WithCredentials(credential)

    sessionCfg, err := session.NewSession(cfg)
    if err != nil {
        err = fmt.Errorf("failed create session, err := %s", err.Error())
        return nil, err
    }

    return sessionCfg, nil
}

func (base *StorageBase) s3PresignUpload(presignConfig PresignUploadConfig) (s3Presign.Forms, error) {
    if base.s3Enabled == nil {
        return s3Presign.Forms{}, fmt.Errorf("aws s3 setting not exists")
    }

    awsConfig := s3Presign.AwsConfig{
        AwsAccessKey: base.s3Enabled.AccessKey,
        AwsRegion:    base.s3Enabled.Region,
        AwsSecretKey: base.s3Enabled.SecretKey,
        AwsBucket:    base.s3Enabled.BucketName,
    }

    timeExpired := presignConfig.ExpiredTime
    if timeExpired == nil {
        defaultExpired := time.Now().Add(15 * time.Minute)
        timeExpired = &defaultExpired
    }

    s3PolicyBase := s3Presign.NewS3Policy(awsConfig)
    s3PolicyBase.SetExpirationDate(*timeExpired)

    //set upload policy
    key := presignConfig.FilePath + presignConfig.Filename
    s3PolicyBase.SetKeyPolicy(s3Presign.ConditionMatchingExactMatch, key)
    s3PolicyBase.SetExpiresPolicy(*timeExpired)

    if presignConfig.MinSizeUpload >= 0 && presignConfig.MaxSizeUpload > 0 {
        s3PolicyBase.SetContentLengthPolicy(presignConfig.MinSizeUpload, presignConfig.MaxSizeUpload)
    }

    if presignConfig.ContentType != "" {
        s3PolicyBase.SetContentTypePolicy(s3Presign.ConditionMatchingExactMatch, presignConfig.ContentType)
    }

    // generated policy
    _, _, formsData := s3PolicyBase.GeneratePolicy()
    return formsData, nil
}

func (base *StorageBase) s3PresignDownload(storageType, storagePath, filename string) (presignUrl string, err error) {
    sessionCfg, err := base.s3Session()
    if err != nil {
        return "", err
    }

    filePath := base.s3Enabled.FilePath + "/" + storageType + storagePath + filename
    disposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
    expiredTime := 5 * time.Minute
    expired := time.Now().Add(expiredTime)
    objectInput := s3.GetObjectInput{
        Bucket:                     aws.String(base.s3Enabled.BucketName),
        Key:                        aws.String(filePath),
        ResponseContentDisposition: aws.String(disposition),
        ResponseExpires:            aws.Time(expired),
    }

    s3Client := s3.New(sessionCfg)

    req, _ := s3Client.GetObjectRequest(&objectInput)
    urlStr, err := req.Presign(expiredTime)
    return urlStr, err
}

func (base *StorageBase) s3Upload(contentTypeData FileData, scaled int, file multipart.File) error {
    sessionCfg, err := base.s3Session()
    if err != nil {
        return fmt.Errorf("failed create session, err := %s", err.Error())
    }

    uploader := s3manager.NewUploader(sessionCfg)

    // create temp file
    fileTemp, err := ioutil.TempFile(os.TempDir(), "prefix")
    if err != nil {
        return fmt.Errorf("bad AWS credentials, err := %s", err.Error())
    }

    if contentTypeData.IsImage {
        // encode all image.Image to jpeg
        // change all image mime to image/jpeg
        var opt jpeg.Options
        opt.Quality = scaled
        err = jpeg.Encode(fileTemp, contentTypeData.ImageFile, &opt)
        contentTypeData.ContentType = "image/jpeg"

        if err != nil {
            return fmt.Errorf("encode image failed -s3-, err := %s", err.Error())
        }
    } else {
        _, err = io.Copy(fileTemp, file)
        if err != nil {
            return fmt.Errorf("error copying data -s3-, err := %s", err.Error())
        }
    }

    _, err = fileTemp.Seek(0, 0)
    if err != nil {
        return fmt.Errorf("error seek file aws to start, err := %s", err.Error())
    }

    params := &s3manager.UploadInput{
        Bucket:      aws.String(base.s3Enabled.BucketName),
        Key:         aws.String(contentTypeData.StoragePath + contentTypeData.Filename),
        Body:        fileTemp,
        ContentType: aws.String(contentTypeData.ContentType),
    }

    _, err = uploader.Upload(params)

    if err != nil {
        return fmt.Errorf("upload to S3 failed, err := %s", err.Error())
    }

    return nil
}

func (base *StorageBase) s3GetFile(storageType, storagePath, filename string) (files *os.File, err error) {
    // Storage on S3
    sessionCfg, err := base.s3Session()
    if err != nil {
        return nil, fmt.Errorf("failed create session, err := %s", err.Error())
    }

    downloader := s3manager.NewDownloader(sessionCfg)
    fileTemp, err := ioutil.TempFile(os.TempDir(), "prefix")
    if err != nil {
        return nil, fmt.Errorf("failed to create temp file, err := %s", err.Error())
    }

    filePath := base.s3Enabled.FilePath + "/" + storageType + storagePath + filename
    _, err = downloader.Download(fileTemp, &s3.GetObjectInput{
        Bucket: aws.String(base.s3Enabled.BucketName),
        Key:    aws.String(filePath),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to download file, %v", err)
    }

    return fileTemp, nil
}
