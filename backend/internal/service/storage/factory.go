package storage

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// NewStorageService は環境変数に基づいて適切なストレージサービスを作成
func NewStorageService() (service.StorageService, error) {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "minio" // デフォルトはMinIO
	}

	switch storageType {
	case "minio":
		return newMinIOFromEnv()
	case "s3":
		return newS3FromEnv()
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// newMinIOFromEnv は環境変数からMinIOストレージを作成
func newMinIOFromEnv() (service.StorageService, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT is required")
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY is required")
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("MINIO_SECRET_KEY is required")
	}

	publicURL := os.Getenv("MINIO_PUBLIC_URL")
	if publicURL == "" {
		publicURL = "http://localhost:9000" // デフォルト値
	}

	useSSLStr := os.Getenv("MINIO_USE_SSL")
	useSSL, _ := strconv.ParseBool(useSSLStr)

	cfg := MinIOConfig{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    useSSL,
		PublicURL: publicURL,
	}

	return NewMinIOStorage(cfg)
}

// newS3FromEnv は環境変数からS3ストレージを作成
func newS3FromEnv() (service.StorageService, error) {
	region := os.Getenv("S3_REGION")
	if region == "" {
		return nil, fmt.Errorf("S3_REGION is required")
	}

	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY_ID is required")
	}

	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, fmt.Errorf("S3_SECRET_ACCESS_KEY is required")
	}

	publicURL := os.Getenv("S3_PUBLIC_URL")
	if publicURL == "" {
		return nil, fmt.Errorf("S3_PUBLIC_URL is required")
	}

	cfg := S3Config{
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		PublicURL:       publicURL,
	}

	return NewS3Storage(cfg)
}
