package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// デフォルトの最大ファイルサイズ
	DefaultMaxImageSize = 5 * 1024 * 1024   // 5MB
	DefaultMaxVideoSize = 100 * 1024 * 1024 // 100MB

	// 一般的な最大ファイルサイズ
	MaxFileSize = 200 * 1024 * 1024 // 200MB
)

// UploadConfig はアップロード設定を保持する
type UploadConfig struct {
	MaxImageSize int64
	MaxVideoSize int64
	MaxFileSize  int64
}

// LoadUploadConfig は環境変数からアップロード設定を読み込む
func LoadUploadConfig() *UploadConfig {
	config := &UploadConfig{
		MaxImageSize: DefaultMaxImageSize,
		MaxVideoSize: DefaultMaxVideoSize,
		MaxFileSize:  MaxFileSize,
	}

	// 環境変数から最大画像サイズを取得
	if maxImageSize := os.Getenv("MAX_IMAGE_SIZE"); maxImageSize != "" {
		if size, err := strconv.ParseInt(maxImageSize, 10, 64); err == nil {
			config.MaxImageSize = size
		}
	}

	// 環境変数から最大動画サイズを取得
	if maxVideoSize := os.Getenv("MAX_VIDEO_SIZE"); maxVideoSize != "" {
		if size, err := strconv.ParseInt(maxVideoSize, 10, 64); err == nil {
			config.MaxVideoSize = size
		}
	}

	return config
}

// LimitFileSize はファイルサイズを制限するミドルウェア
func LimitFileSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Lengthヘッダーをチェック
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file_too_large",
				"message": fmt.Sprintf("ファイルサイズが上限を超えています（最大: %s）",
					formatFileSize(maxSize)),
			})
			c.Abort()
			return
		}

		// MaxBytesReaderでボディサイズを制限
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}

// LimitImageUpload は画像アップロード用のファイルサイズ制限ミドルウェア
func LimitImageUpload() gin.HandlerFunc {
	config := LoadUploadConfig()
	return LimitFileSize(config.MaxImageSize)
}

// LimitVideoUpload は動画アップロード用のファイルサイズ制限ミドルウェア
func LimitVideoUpload() gin.HandlerFunc {
	config := LoadUploadConfig()
	return LimitFileSize(config.MaxVideoSize)
}

// LimitMediaUpload はメディアアップロード用のファイルサイズ制限ミドルウェア
// メディアタイプに応じて動的にサイズ制限を適用
func LimitMediaUpload() gin.HandlerFunc {
	config := LoadUploadConfig()

	return func(c *gin.Context) {
		// multipart formからmedia_typeを取得
		mediaType := c.PostForm("media_type")

		var maxSize int64
		switch mediaType {
		case "image":
			maxSize = config.MaxImageSize
		case "video":
			maxSize = config.MaxVideoSize
		default:
			// メディアタイプが不明な場合は画像サイズを適用
			maxSize = config.MaxImageSize
		}

		// Content-Lengthヘッダーをチェック
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file_too_large",
				"message": fmt.Sprintf("ファイルサイズが上限を超えています（%s: 最大%s）",
					getMediaTypeName(mediaType), formatFileSize(maxSize)),
			})
			c.Abort()
			return
		}

		// MaxBytesReaderでボディサイズを制限
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}

// ValidateMultipartForm はmultipart/form-dataのバリデーションを行う
func ValidateMultipartForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Typeがmultipart/form-dataか確認
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_content_type",
				"message": "Content-Typeはmultipart/form-dataである必要があります",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireFileUpload はファイルアップロードが必須であることを検証する
func RequireFileUpload(fileFieldName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile(fileFieldName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_required",
				"message": "ファイルが選択されていません",
			})
			c.Abort()
			return
		}

		if file.Size == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "empty_file",
				"message": "ファイルが空です",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateFileSize はアップロードされたファイルのサイズを検証する
func ValidateFileSize(fileFieldName string, maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile(fileFieldName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "file_error",
				"message": "ファイルの読み取りに失敗しました",
			})
			c.Abort()
			return
		}

		if file.Size > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "file_too_large",
				"message": fmt.Sprintf("ファイルサイズが上限を超えています（最大: %s）",
					formatFileSize(maxSize)),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateContentLength はContent-Lengthヘッダーを検証する
func ValidateContentLength(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength < 0 {
			c.JSON(http.StatusLengthRequired, gin.H{
				"error":   "content_length_required",
				"message": "Content-Lengthヘッダーが必要です",
			})
			c.Abort()
			return
		}

		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request_too_large",
				"message": fmt.Sprintf("リクエストサイズが上限を超えています（最大: %s）",
					formatFileSize(maxSize)),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HandleUploadError はアップロードエラーをハンドリングする
func HandleUploadError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// エラーがあればレスポンス
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// MaxBytesReaderのエラーをチェック
			if strings.Contains(err.Error(), "http: request body too large") {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error":   "file_too_large",
					"message": "ファイルサイズが上限を超えています",
				})
				return
			}

			// その他のエラー
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "upload_error",
				"message": "ファイルのアップロードに失敗しました",
			})
		}
	}
}

// formatFileSize はバイト数を人間が読みやすい形式にフォーマットする
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// getMediaTypeName はメディアタイプの日本語名を返す
func getMediaTypeName(mediaType string) string {
	switch mediaType {
	case "image":
		return "画像"
	case "video":
		return "動画"
	default:
		return "ファイル"
	}
}

// SetMaxMultipartMemory はmultipart formの最大メモリサイズを設定する
func SetMaxMultipartMemory(maxMemory int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.ParseMultipartForm(maxMemory)
		c.Next()
	}
}

// GetUploadConfig は現在のアップロード設定を返す（デバッグ用）
func GetUploadConfig(c *gin.Context) {
	config := LoadUploadConfig()
	c.JSON(http.StatusOK, gin.H{
		"max_image_size": formatFileSize(config.MaxImageSize),
		"max_video_size": formatFileSize(config.MaxVideoSize),
		"max_file_size":  formatFileSize(config.MaxFileSize),
	})
}

// CombineUploadMiddleware は複数のアップロード用ミドルウェアを組み合わせる
func CombineUploadMiddleware(fileFieldName string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ValidateMultipartForm(),
		LimitMediaUpload(),
		RequireFileUpload(fileFieldName),
		HandleUploadError(),
	}
}

// ImageUploadMiddleware は画像アップロード用のミドルウェアチェーンを返す
func ImageUploadMiddleware(fileFieldName string) []gin.HandlerFunc {
	config := LoadUploadConfig()
	return []gin.HandlerFunc{
		ValidateMultipartForm(),
		LimitFileSize(config.MaxImageSize),
		RequireFileUpload(fileFieldName),
		HandleUploadError(),
	}
}

// VideoUploadMiddleware は動画アップロード用のミドルウェアチェーンを返す
func VideoUploadMiddleware(fileFieldName string) []gin.HandlerFunc {
	config := LoadUploadConfig()
	return []gin.HandlerFunc{
		ValidateMultipartForm(),
		LimitFileSize(config.MaxVideoSize),
		RequireFileUpload(fileFieldName),
		HandleUploadError(),
	}
}
