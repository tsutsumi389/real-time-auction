package service

import (
	"context"
	"io"
)

// StorageService はオブジェクトストレージの抽象化インターフェース
// MinIO、AWS S3、GCS等の実装を切り替え可能にする
type StorageService interface {
	// Upload uploads a file to the specified bucket
	// ctx: context for cancellation
	// bucket: bucket name
	// objectName: object key/name in the bucket
	// reader: file content reader
	// fileSize: size of the file in bytes
	// contentType: MIME type of the file
	// Returns the public URL of the uploaded object and any error
	Upload(ctx context.Context, bucket, objectName string, reader io.Reader, fileSize int64, contentType string) (string, error)

	// Delete removes an object from the specified bucket
	// ctx: context for cancellation
	// bucket: bucket name
	// objectName: object key/name to delete
	Delete(ctx context.Context, bucket, objectName string) error

	// GetPublicURL returns the public URL for accessing an object
	// bucket: bucket name
	// objectName: object key/name
	GetPublicURL(bucket, objectName string) string

	// HealthCheck verifies the storage service is accessible
	// ctx: context for cancellation
	HealthCheck(ctx context.Context) error
}
