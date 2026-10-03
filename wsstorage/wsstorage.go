// Package wsstorage streams an io.Reader to S3, GCS or local disk without buffering the whole
// body in memory or in a temp file.
//
// Path and filename generation mirror the storages package, so records created here can be read
// back with storages.GetFiles.
package wsstorage

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	stdmime "mime"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	manager "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/api/option"
)

type Driver string

const (
	DriverLocal Driver = "local"
	DriverS3    Driver = "s3"
	DriverGCS   Driver = "gcs"
)

var (
	ErrEmptyFile     = errors.New("file is empty")
	ErrFileTooLarge  = errors.New("file exceeds the maximum allowed size")
	ErrInvalidType   = errors.New("invalid file type")
	errStoragePath   = errors.New("storage path is not configured")
	fileTypePattern  = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	extensionPattern = regexp.MustCompile(`^\.[A-Za-z0-9]{1,10}$`)
)

type S3Config struct {
	Path      string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
}

type GCSConfig struct {
	Path           string
	Bucket         string
	CredentialPath string
}

type Config struct {
	Driver    Driver
	LocalPath string
	S3        S3Config
	GCS       GCSConfig
}

// ConfigFromEnv reads STORAGE_DRIVER (local|s3|gcs, default local) and the same env vars go-helper uses.
func ConfigFromEnv() Config {
	driver := Driver(strings.ToLower(os.Getenv("STORAGE_DRIVER")))
	if driver != DriverS3 && driver != DriverGCS {
		driver = DriverLocal
	}

	return Config{
		Driver:    driver,
		LocalPath: os.Getenv("LOCAL_STORAGE_PATH"),
		S3: S3Config{
			Path:      os.Getenv("S3_STORAGE_PATH"),
			Bucket:    os.Getenv("S3_BUCKET"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			Region:    os.Getenv("S3_REGION"),
		},
		GCS: GCSConfig{
			Path:           os.Getenv("GCS_LOCATION"),
			Bucket:         os.Getenv("GCS_BUCKET_NAME"),
			CredentialPath: os.Getenv("GCS_CREDENTIAL_PATH"),
		},
	}
}

// Result has the same shape as go-helper's storages.StorageData, plus the stored size.
type Result struct {
	Type             string
	Path             string
	Filename         string
	Mime             string
	OriginalFilename string
	Size             int64
}

// ValidFileType reports whether fileType is safe to use as a storage path segment.
func ValidFileType(fileType string) bool {
	return fileTypePattern.MatchString(fileType)
}

// Upload streams body to the configured backend. At most maxBytes are accepted (maxBytes <= 0 means
// unlimited); a longer body aborts the upload with ErrFileTooLarge and nothing is stored. The MIME
// type is sniffed from the first bytes, never taken from the client. Content is stored unmodified.
func Upload(ctx context.Context, cfg Config, fileType, originalName string, body io.Reader, maxBytes int64) (*Result, error) {
	if !ValidFileType(fileType) {
		return nil, ErrInvalidType
	}

	limited := &limitReader{r: body, remaining: maxBytes, limited: maxBytes > 0}
	buffered := bufio.NewReaderSize(limited, 4096)

	head, err := buffered.Peek(3072)
	if len(head) == 0 {
		if err == nil || errors.Is(err, io.EOF) {
			return nil, ErrEmptyFile
		}
		return nil, err
	}

	mime := detectMime(head, originalName)
	name := generateName(originalName, safeExtension(originalName))
	datePath := time.Now().Local().Format("2006/01/02/")

	switch cfg.Driver {
	case DriverS3:
		err = uploadS3(ctx, cfg.S3, fileType, datePath, name, mime, buffered)
	case DriverGCS:
		err = uploadGCS(ctx, cfg.GCS, fileType, datePath, name, buffered)
	default:
		err = uploadLocal(cfg.LocalPath, fileType, datePath, name, buffered)
	}
	if err != nil {
		if limited.exceeded {
			return nil, ErrFileTooLarge
		}
		return nil, err
	}

	return &Result{
		Type:             fileType,
		Path:             datePath,
		Filename:         name,
		Mime:             mime,
		OriginalFilename: originalName,
		Size:             limited.read,
	}, nil
}

// detectMime sniffs the magic bytes; audio/video without any (octet-stream) is resolved by extension.
func detectMime(head []byte, originalName string) string {
	detected := mimetype.Detect(head).String()
	if detected != "application/octet-stream" {
		return detected
	}

	byExt, _, err := stdmime.ParseMediaType(stdmime.TypeByExtension(strings.ToLower(path.Ext(originalName))))
	if err == nil && (strings.HasPrefix(byExt, "audio/") || strings.HasPrefix(byExt, "video/")) {
		return byExt
	}
	return detected
}

// limitReader fails (instead of silently truncating) once more than the limit has been read.
type limitReader struct {
	r         io.Reader
	remaining int64
	limited   bool
	read      int64
	exceeded  bool
}

func (l *limitReader) Read(p []byte) (int, error) {
	if l.limited {
		if l.remaining <= 0 {
			// the limit is fully consumed; only EOF here means the body was exactly at the limit
			var one [1]byte
			n, err := l.r.Read(one[:])
			if n > 0 {
				l.exceeded = true
				return 0, ErrFileTooLarge
			}
			return 0, err
		}

		if int64(len(p)) > l.remaining {
			p = p[:l.remaining]
		}
	}

	n, err := l.r.Read(p)
	l.read += int64(n)
	l.remaining -= int64(n)
	return n, err
}

func objectKey(base, fileType, datePath, name string) string {
	return fmt.Sprintf("%s/%s/%s%s", base, fileType, datePath, name)
}

func uploadS3(ctx context.Context, cfg S3Config, fileType, datePath, name, mime string, body io.Reader) error {
	awsCfg := aws.Config{
		Region: cfg.Region,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	}

	uploader := manager.New(s3.NewFromConfig(awsCfg))
	_, err := uploader.UploadObject(ctx, &manager.UploadObjectInput{
		Bucket:      aws.String(cfg.Bucket),
		Key:         aws.String(objectKey(cfg.Path, fileType, datePath, name)),
		Body:        body,
		ContentType: aws.String(mime),
	})
	if err != nil {
		return fmt.Errorf("upload to S3 failed, err := %w", err)
	}

	return nil
}

func uploadGCS(ctx context.Context, cfg GCSConfig, fileType, datePath, name string, body io.Reader) error {
	var opts []option.ClientOption
	if cfg.CredentialPath != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialPath))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	defer client.Close()

	// cancelling the context before Close discards the partial object
	writeCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	writer := client.Bucket(cfg.Bucket).Object(objectKey(cfg.Path, fileType, datePath, name)).NewWriter(writeCtx)
	if _, err = io.Copy(writer, body); err != nil {
		cancel()
		_ = writer.Close()
		return fmt.Errorf("error copying data -gcs-, err := %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("error close writer -gcs-: %w", err)
	}

	return nil
}

func uploadLocal(basePath, fileType, datePath, name string, body io.Reader) (err error) {
	if basePath == "" {
		return errStoragePath
	}

	dir := fmt.Sprintf("%s/%s/%s", basePath, fileType, datePath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := dir + name
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("file not created, err := %w", err)
	}

	defer func() {
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			_ = os.Remove(filePath)
		}
	}()

	if _, err = io.Copy(out, body); err != nil {
		return fmt.Errorf("error upload file to local, err := %w", err)
	}

	return nil
}

func safeExtension(originalName string) string {
	ext := filepath.Ext(path.Base(strings.ReplaceAll(originalName, "\\", "/")))
	if !extensionPattern.MatchString(ext) {
		return ""
	}
	return strings.ToLower(ext)
}

// generateName mirrors go-helper: md5(name + timestamp) + extension. A nanosecond suffix is added so
// two uploads of the same name in the same second do not collide.
func generateName(originalName, ext string) string {
	seed := originalName + time.Now().Local().Format(time.UnixDate) + strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := md5.Sum([]byte(seed))
	return hex.EncodeToString(hash[:]) + ext
}
