package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// MinIOStorage はMinIOベースのストレージサービス実装
type MinIOStorage struct {
	client    *minio.Client
	publicURL string
}

// MinIOConfig はMinIO接続設定
type MinIOConfig struct {
	Endpoint  string // MinIOエンドポイント（例: localhost:9000）
	AccessKey string // アクセスキー
	SecretKey string // シークレットキー
	UseSSL    bool   // SSL使用フラグ
	PublicURL string // 公開URL（例: http://localhost:9000）
}

// NewMinIOStorage は新しいMinIOStorageインスタンスを作成
func NewMinIOStorage(cfg MinIOConfig) (service.StorageService, error) {
	// MinIOクライアントを初期化
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOStorage{
		client:    client,
		publicURL: cfg.PublicURL,
	}, nil
}

// Upload uploads a file to MinIO
func (s *MinIOStorage) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	// バケットが存在するか確認
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// バケットが存在しない場合は作成
	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// オブジェクトをアップロード
	_, err = s.client.PutObject(ctx, bucket, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	// 公開URLを生成
	publicURL := s.GetPublicURL(bucket, objectName)
	return publicURL, nil
}

// Delete removes an object from MinIO
func (s *MinIOStorage) Delete(ctx context.Context, bucket, objectName string) error {
	err := s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// GetPublicURL returns the public URL for accessing an object
func (s *MinIOStorage) GetPublicURL(bucket, objectName string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicURL, bucket, objectName)
}

// HealthCheck verifies MinIO is accessible
func (s *MinIOStorage) HealthCheck(ctx context.Context) error {
	// List buckets to verify connection
	_, err := s.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO health check failed: %w", err)
	}
	return nil
}
