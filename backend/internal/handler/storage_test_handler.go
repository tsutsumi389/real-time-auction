package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// StorageTestHandler はストレージサービスのテストハンドラー（開発用）
type StorageTestHandler struct {
	storageService service.StorageService
}

// NewStorageTestHandler は新しいStorageTestHandlerを作成
func NewStorageTestHandler(storageService service.StorageService) *StorageTestHandler {
	return &StorageTestHandler{
		storageService: storageService,
	}
}

// TestUpload はストレージへのアップロードをテスト
func (h *StorageTestHandler) TestUpload(c *gin.Context) {
	// テスト用の小さなテキストファイルを作成
	testContent := []byte("This is a test file for MinIO storage")
	reader := bytes.NewReader(testContent)

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "auction-media"
	}

	objectName := "test/test-file.txt"

	// アップロード実行
	url, err := h.storageService.Upload(
		c.Request.Context(),
		bucket,
		objectName,
		reader,
		int64(len(testContent)),
		"text/plain",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "upload_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     url,
		"bucket":  bucket,
		"object":  objectName,
	})
}

// TestHealthCheck はストレージサービスのヘルスチェックをテスト
func (h *StorageTestHandler) TestHealthCheck(c *gin.Context) {
	err := h.storageService.HealthCheck(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "health_check_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Storage service is healthy",
	})
}

// TestDelete はストレージからの削除をテスト
func (h *StorageTestHandler) TestDelete(c *gin.Context) {
	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "auction-media"
	}

	objectName := "test/test-file.txt"

	err := h.storageService.Delete(c.Request.Context(), bucket, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "delete_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File %s deleted successfully", objectName),
		"bucket":  bucket,
		"object":  objectName,
	})
}
