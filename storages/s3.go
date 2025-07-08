package storages

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fari-99/aws-presignpost-s3-go"
)

func (base *StorageBase) s3Config() (aws.Config, error) {
	awsAccessKey := aws.String(base.s3Enabled.AccessKey)
	awsSecretKey := aws.String(base.s3Enabled.SecretKey)
	token := ""

	credential := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			*awsAccessKey, *awsSecretKey, token,
		),
	)

	cfg := aws.Config{
		Credentials: credential,
	}

	// Optionally test credentials
	_, err := cfg.Credentials.Retrieve(base.ctx)
	if err != nil {
		errMsg := "Failed to get AWS credentials. Error: " + err.Error()
		log.Println(errMsg)
		return aws.Config{}, err
	}

	return cfg, nil
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

	// set upload policy
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
	s3Config, err := base.s3Config()
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

	s3Client := s3.NewFromConfig(s3Config)
	presigner := s3.NewPresignClient(s3Client)

	preSignedReq, err := presigner.PresignGetObject(
		base.ctx,
		&objectInput,
		s3.WithPresignExpires(expiredTime),
	)

	if err != nil {
		errMsg := "Failed to get preSign url. Error: " + err.Error()
		log.Println(errMsg)
		return "", err
	}

	return preSignedReq.URL, err
}

func (base *StorageBase) s3Upload(contentTypeData FileData, scaled int, file multipart.File) error {
	s3Config, err := base.s3Config()
	if err != nil {
		return fmt.Errorf("failed create session, err := %s", err.Error())
	}

	s3Client := s3.NewFromConfig(s3Config)
	uploader := manager.NewUploader(s3Client)

	// create temp file
	fileTemp, err := os.CreateTemp(os.TempDir(), "prefix")
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

	params := s3.PutObjectInput{
		Bucket:      aws.String(base.s3Enabled.BucketName),
		Key:         aws.String(contentTypeData.StoragePath + contentTypeData.Filename),
		Body:        fileTemp,
		ContentType: aws.String(contentTypeData.ContentType),
	}

	_, err = uploader.Upload(base.ctx, &params)

	if err != nil {
		return fmt.Errorf("upload to S3 failed, err := %s", err.Error())
	}

	return nil
}

func (base *StorageBase) s3GetFile(storageType, storagePath, filename string) (files *os.File, err error) {
	// Storage on S3
	s3Config, err := base.s3Config()
	if err != nil {
		return nil, fmt.Errorf("failed create session, err := %s", err.Error())
	}

	s3Client := s3.NewFromConfig(s3Config)
	downloader := manager.NewDownloader(s3Client)

	fileTemp, err := os.CreateTemp(os.TempDir(), "prefix")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file, err := %s", err.Error())
	}

	filePath := base.s3Enabled.FilePath + "/" + storageType + storagePath + filename
	_, err = downloader.Download(base.ctx, fileTemp, &s3.GetObjectInput{
		Bucket: aws.String(base.s3Enabled.BucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file, %v", err)
	}

	return fileTemp, nil
}
