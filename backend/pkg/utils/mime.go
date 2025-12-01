package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ErrInvalidMimeType = errors.New("invalid MIME type")
	ErrFileRead        = errors.New("failed to read file")
)

// AllowedImageTypes はアップロード可能な画像のMIMEタイプ
var AllowedImageTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	"image/gif",
}

// AllowedVideoTypes はアップロード可能な動画のMIMEタイプ
var AllowedVideoTypes = []string{
	"video/mp4",
	"video/quicktime",
	"video/x-msvideo",
}

// ValidateMimeType はファイルのMIMEタイプを検証する（マジックバイトを使用）
func ValidateMimeType(filePath string, allowedTypes []string) error {
	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return ErrFileRead
	}
	defer file.Close()

	// 最初の512バイトを読み取る（MIMEタイプ検出に必要）
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return ErrFileRead
	}

	// マジックバイトからMIMEタイプを検出
	detectedType := http.DetectContentType(buffer[:n])

	// セミコロン以降のパラメータを除去（例: "image/jpeg; charset=utf-8" → "image/jpeg"）
	detectedType = strings.Split(detectedType, ";")[0]
	detectedType = strings.TrimSpace(detectedType)

	// 許可されたタイプに含まれているか確認
	for _, allowedType := range allowedTypes {
		if detectedType == allowedType {
			return nil
		}
	}

	return ErrInvalidMimeType
}

// ValidateMimeTypeFromBytes はバイト列からMIMEタイプを検証する
func ValidateMimeTypeFromBytes(data []byte, allowedTypes []string) error {
	// マジックバイトからMIMEタイプを検出
	detectedType := http.DetectContentType(data)

	// セミコロン以降のパラメータを除去
	detectedType = strings.Split(detectedType, ";")[0]
	detectedType = strings.TrimSpace(detectedType)

	// 許可されたタイプに含まれているか確認
	for _, allowedType := range allowedTypes {
		if detectedType == allowedType {
			return nil
		}
	}

	return ErrInvalidMimeType
}

// DetectMimeType はファイルのMIMEタイプを検出する
func DetectMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", ErrFileRead
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", ErrFileRead
	}

	detectedType := http.DetectContentType(buffer[:n])
	detectedType = strings.Split(detectedType, ";")[0]
	detectedType = strings.TrimSpace(detectedType)

	return detectedType, nil
}

// DetectMimeTypeFromReader はio.Readerからマジックバイトを読み取ってMIMEタイプを検出する
func DetectMimeTypeFromReader(reader io.Reader) (string, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", ErrFileRead
	}

	detectedType := http.DetectContentType(buffer[:n])
	detectedType = strings.Split(detectedType, ";")[0]
	detectedType = strings.TrimSpace(detectedType)

	return detectedType, nil
}

// IsImage はファイルが画像かどうかを判定する
func IsImage(filePath string) bool {
	mimeType, err := DetectMimeType(filePath)
	if err != nil {
		return false
	}

	for _, allowedType := range AllowedImageTypes {
		if mimeType == allowedType {
			return true
		}
	}

	return false
}

// IsVideo はファイルが動画かどうかを判定する
func IsVideo(filePath string) bool {
	mimeType, err := DetectMimeType(filePath)
	if err != nil {
		return false
	}

	for _, allowedType := range AllowedVideoTypes {
		if mimeType == allowedType {
			return true
		}
	}

	return false
}

// GetMediaType はファイルのメディアタイプ（image/video）を返す
func GetMediaType(filePath string) (string, error) {
	if IsImage(filePath) {
		return "image", nil
	}
	if IsVideo(filePath) {
		return "video", nil
	}
	return "", ErrInvalidMimeType
}

// ValidateImageFile は画像ファイルのMIMEタイプを検証する
func ValidateImageFile(filePath string) error {
	return ValidateMimeType(filePath, AllowedImageTypes)
}

// ValidateVideoFile は動画ファイルのMIMEタイプを検証する
func ValidateVideoFile(filePath string) error {
	return ValidateMimeType(filePath, AllowedVideoTypes)
}

// ReadFileHeader はファイルの先頭部分を読み取る（MIMEタイプ検出用）
func ReadFileHeader(reader io.Reader) ([]byte, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, ErrFileRead
	}
	return buffer[:n], nil
}

// IsMimeTypeAllowed は指定されたMIMEタイプが許可リストに含まれているか確認する
func IsMimeTypeAllowed(mimeType string, allowedTypes []string) bool {
	mimeType = strings.Split(mimeType, ";")[0]
	mimeType = strings.TrimSpace(mimeType)

	for _, allowedType := range allowedTypes {
		if mimeType == allowedType {
			return true
		}
	}
	return false
}

// CreateMimeValidator は特定のMIMEタイプリストに対するバリデーター関数を返す
func CreateMimeValidator(allowedTypes []string) func(string) error {
	return func(filePath string) error {
		return ValidateMimeType(filePath, allowedTypes)
	}
}

// GetMimeTypeCategory はMIMEタイプのカテゴリ（image, video, application等）を返す
func GetMimeTypeCategory(mimeType string) string {
	parts := strings.Split(mimeType, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// ValidateMultipartFile はmultipart.FileHeaderからMIMEタイプを検証する
func ValidateMultipartFile(fileData []byte, allowedTypes []string) error {
	if len(fileData) == 0 {
		return ErrFileRead
	}

	// 最大512バイトを使用してMIMEタイプを検出
	maxLen := 512
	if len(fileData) < maxLen {
		maxLen = len(fileData)
	}

	detectedType := http.DetectContentType(fileData[:maxLen])
	detectedType = strings.Split(detectedType, ";")[0]
	detectedType = strings.TrimSpace(detectedType)

	for _, allowedType := range allowedTypes {
		if detectedType == allowedType {
			return nil
		}
	}

	return ErrInvalidMimeType
}

// GetFileExtensionFromMime はMIMEタイプから推奨ファイル拡張子を返す
func GetFileExtensionFromMime(mimeType string) string {
	mimeType = strings.Split(mimeType, ";")[0]
	mimeType = strings.TrimSpace(mimeType)

	extensions := map[string]string{
		"image/jpeg":         ".jpg",
		"image/png":          ".png",
		"image/webp":         ".webp",
		"image/gif":          ".gif",
		"video/mp4":          ".mp4",
		"video/quicktime":    ".mov",
		"video/x-msvideo":    ".avi",
		"application/pdf":    ".pdf",
		"application/zip":    ".zip",
		"application/json":   ".json",
		"text/plain":         ".txt",
		"text/html":          ".html",
		"text/css":           ".css",
		"text/javascript":    ".js",
		"application/xml":    ".xml",
	}

	if ext, ok := extensions[mimeType]; ok {
		return ext
	}

	// デフォルトは空文字
	return ""
}

// BufferedReader はReaderをBufferに変換し、再利用可能にする
type BufferedReader struct {
	buffer   *bytes.Buffer
	original []byte
}

// NewBufferedReader は新しいBufferedReaderを作成する
func NewBufferedReader(reader io.Reader) (*BufferedReader, error) {
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return nil, err
	}
	original := buffer.Bytes()
	return &BufferedReader{
		buffer:   bytes.NewBuffer(original),
		original: original,
	}, nil
}

// Read はBufferから読み取る
func (br *BufferedReader) Read(p []byte) (n int, err error) {
	return br.buffer.Read(p)
}

// Bytes はBuffer全体を返す
func (br *BufferedReader) Bytes() []byte {
	return br.original
}

// Reset はBufferを先頭に戻す
func (br *BufferedReader) Reset() {
	br.buffer = bytes.NewBuffer(br.original)
}
