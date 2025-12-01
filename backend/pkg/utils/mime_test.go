package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

// createTestImageFile はテスト用の画像ファイルを作成する
func createTestImageFile(t *testing.T, format string) string {
	t.Helper()

	// 1x1ピクセルの画像を作成
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test-image-*."+format)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// フォーマットに応じてエンコード
	switch format {
	case "jpg", "jpeg":
		if err := jpeg.Encode(tmpFile, img, nil); err != nil {
			t.Fatalf("Failed to encode JPEG: %v", err)
		}
	case "png":
		if err := png.Encode(tmpFile, img); err != nil {
			t.Fatalf("Failed to encode PNG: %v", err)
		}
	default:
		t.Fatalf("Unsupported format: %s", format)
	}

	return tmpFile.Name()
}

// createTestTextFile はテスト用のテキストファイルを作成する
func createTestTextFile(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-text-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString("This is a text file"); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestValidateMimeType(t *testing.T) {
	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		allowedTypes  []string
		expectedError error
	}{
		{
			name:          "Valid JPEG image",
			setupFile:     func(t *testing.T) string { return createTestImageFile(t, "jpg") },
			allowedTypes:  AllowedImageTypes,
			expectedError: nil,
		},
		{
			name:          "Valid PNG image",
			setupFile:     func(t *testing.T) string { return createTestImageFile(t, "png") },
			allowedTypes:  AllowedImageTypes,
			expectedError: nil,
		},
		{
			name:          "Invalid text file as image",
			setupFile:     createTestTextFile,
			allowedTypes:  AllowedImageTypes,
			expectedError: ErrInvalidMimeType,
		},
		{
			name:          "Non-existent file",
			setupFile:     func(t *testing.T) string { return "/non/existent/file.jpg" },
			allowedTypes:  AllowedImageTypes,
			expectedError: ErrFileRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile(t)
			if tt.expectedError != ErrFileRead {
				defer os.Remove(filePath)
			}

			err := ValidateMimeType(filePath, tt.allowedTypes)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestDetectMimeType(t *testing.T) {
	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		expectedType string
	}{
		{
			name:         "Detect JPEG",
			setupFile:    func(t *testing.T) string { return createTestImageFile(t, "jpg") },
			expectedType: "image/jpeg",
		},
		{
			name:         "Detect PNG",
			setupFile:    func(t *testing.T) string { return createTestImageFile(t, "png") },
			expectedType: "image/png",
		},
		{
			name:         "Detect text file",
			setupFile:    createTestTextFile,
			expectedType: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile(t)
			defer os.Remove(filePath)

			mimeType, err := DetectMimeType(filePath)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if mimeType != tt.expectedType {
				t.Errorf("Expected MIME type %s, got %s", tt.expectedType, mimeType)
			}
		})
	}
}

func TestValidateMimeTypeFromBytes(t *testing.T) {
	// JPEG magic bytes
	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	// PNG magic bytes
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	// Text
	textBytes := []byte("This is text")

	tests := []struct {
		name          string
		data          []byte
		allowedTypes  []string
		expectedError error
	}{
		{
			name:          "Valid JPEG bytes",
			data:          jpegBytes,
			allowedTypes:  AllowedImageTypes,
			expectedError: nil,
		},
		{
			name:          "Valid PNG bytes",
			data:          pngBytes,
			allowedTypes:  AllowedImageTypes,
			expectedError: nil,
		},
		{
			name:          "Invalid text bytes",
			data:          textBytes,
			allowedTypes:  AllowedImageTypes,
			expectedError: ErrInvalidMimeType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMimeTypeFromBytes(tt.data, tt.allowedTypes)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  func(t *testing.T) string
		wantIsImage bool
	}{
		{
			name:       "JPEG is image",
			setupFile:  func(t *testing.T) string { return createTestImageFile(t, "jpg") },
			wantIsImage: true,
		},
		{
			name:       "PNG is image",
			setupFile:  func(t *testing.T) string { return createTestImageFile(t, "png") },
			wantIsImage: true,
		},
		{
			name:       "Text is not image",
			setupFile:  createTestTextFile,
			wantIsImage: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile(t)
			defer os.Remove(filePath)

			isImage := IsImage(filePath)

			if isImage != tt.wantIsImage {
				t.Errorf("Expected IsImage=%v, got %v", tt.wantIsImage, isImage)
			}
		})
	}
}

func TestGetMediaType(t *testing.T) {
	tests := []struct {
		name              string
		setupFile         func(t *testing.T) string
		expectedMediaType string
		expectError       bool
	}{
		{
			name:              "Image file returns 'image'",
			setupFile:         func(t *testing.T) string { return createTestImageFile(t, "jpg") },
			expectedMediaType: "image",
			expectError:       false,
		},
		{
			name:              "Text file returns error",
			setupFile:         createTestTextFile,
			expectedMediaType: "",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile(t)
			defer os.Remove(filePath)

			mediaType, err := GetMediaType(filePath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if mediaType != tt.expectedMediaType {
					t.Errorf("Expected media type %s, got %s", tt.expectedMediaType, mediaType)
				}
			}
		})
	}
}

func TestIsMimeTypeAllowed(t *testing.T) {
	tests := []struct {
		name         string
		mimeType     string
		allowedTypes []string
		expected     bool
	}{
		{
			name:         "Allowed type",
			mimeType:     "image/jpeg",
			allowedTypes: AllowedImageTypes,
			expected:     true,
		},
		{
			name:         "Not allowed type",
			mimeType:     "application/pdf",
			allowedTypes: AllowedImageTypes,
			expected:     false,
		},
		{
			name:         "Type with charset parameter",
			mimeType:     "image/jpeg; charset=utf-8",
			allowedTypes: AllowedImageTypes,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMimeTypeAllowed(tt.mimeType, tt.allowedTypes)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetFileExtensionFromMime(t *testing.T) {
	tests := []struct {
		mimeType string
		expected string
	}{
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"image/webp", ".webp"},
		{"image/gif", ".gif"},
		{"video/mp4", ".mp4"},
		{"video/quicktime", ".mov"},
		{"application/pdf", ".pdf"},
		{"unknown/type", ""},
		{"image/jpeg; charset=utf-8", ".jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			ext := GetFileExtensionFromMime(tt.mimeType)

			if ext != tt.expected {
				t.Errorf("Expected extension %s, got %s", tt.expected, ext)
			}
		})
	}
}

func TestGetMimeTypeCategory(t *testing.T) {
	tests := []struct {
		mimeType string
		expected string
	}{
		{"image/jpeg", "image"},
		{"video/mp4", "video"},
		{"text/plain", "text"},
		{"application/json", "application"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.mimeType, func(t *testing.T) {
			category := GetMimeTypeCategory(tt.mimeType)

			if category != tt.expected {
				t.Errorf("Expected category %s, got %s", tt.expected, category)
			}
		})
	}
}

func TestValidateMultipartFile(t *testing.T) {
	// JPEG magic bytes
	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	// Text
	textBytes := []byte("This is text")

	tests := []struct {
		name          string
		data          []byte
		allowedTypes  []string
		expectedError error
	}{
		{
			name:          "Valid JPEG data",
			data:          jpegBytes,
			allowedTypes:  AllowedImageTypes,
			expectedError: nil,
		},
		{
			name:          "Invalid text data",
			data:          textBytes,
			allowedTypes:  AllowedImageTypes,
			expectedError: ErrInvalidMimeType,
		},
		{
			name:          "Empty data",
			data:          []byte{},
			allowedTypes:  AllowedImageTypes,
			expectedError: ErrFileRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMultipartFile(tt.data, tt.allowedTypes)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestBufferedReader(t *testing.T) {
	data := []byte("Hello, World!")
	reader := bytes.NewReader(data)

	br, err := NewBufferedReader(reader)
	if err != nil {
		t.Fatalf("Failed to create BufferedReader: %v", err)
	}

	// Bytesで全データを取得
	if !bytes.Equal(br.Bytes(), data) {
		t.Errorf("Expected %v, got %v", data, br.Bytes())
	}

	// Readで読み取り
	buf := make([]byte, 5)
	n, err := br.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
	if n != 5 {
		t.Errorf("Expected to read 5 bytes, got %d", n)
	}
	if string(buf) != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", string(buf))
	}

	// Resetして再度読み取り
	br.Reset()
	n, err = br.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read after reset: %v", err)
	}
	if string(buf) != "Hello" {
		t.Errorf("Expected 'Hello' after reset, got '%s'", string(buf))
	}
}

func TestCreateMimeValidator(t *testing.T) {
	validator := CreateMimeValidator(AllowedImageTypes)

	jpegFile := createTestImageFile(t, "jpg")
	defer os.Remove(jpegFile)

	textFile := createTestTextFile(t)
	defer os.Remove(textFile)

	// 画像ファイルは通過
	if err := validator(jpegFile); err != nil {
		t.Errorf("JPEG file should be valid: %v", err)
	}

	// テキストファイルはエラー
	if err := validator(textFile); err == nil {
		t.Error("Text file should be invalid")
	}
}

func TestDetectMimeTypeFromReader(t *testing.T) {
	// JPEG magic bytes
	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	reader := bytes.NewReader(jpegBytes)

	mimeType, err := DetectMimeTypeFromReader(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mimeType != "image/jpeg" {
		t.Errorf("Expected 'image/jpeg', got '%s'", mimeType)
	}
}

func TestReadFileHeader(t *testing.T) {
	data := []byte("This is a test file with more than 512 bytes of data")
	reader := bytes.NewReader(data)

	header, err := ReadFileHeader(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(header) != len(data) {
		t.Errorf("Expected header length %d, got %d", len(data), len(header))
	}

	if !bytes.Equal(header, data) {
		t.Error("Header content does not match original data")
	}
}

func TestValidateImageFile(t *testing.T) {
	jpegFile := createTestImageFile(t, "jpg")
	defer os.Remove(jpegFile)

	if err := ValidateImageFile(jpegFile); err != nil {
		t.Errorf("JPEG file should be valid: %v", err)
	}

	textFile := createTestTextFile(t)
	defer os.Remove(textFile)

	if err := ValidateImageFile(textFile); err == nil {
		t.Error("Text file should be invalid")
	}
}

// Benchmark tests
func BenchmarkDetectMimeType(b *testing.B) {
	// テスト用のJPEGファイルを作成
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	tmpFile, _ := os.CreateTemp("", "bench-*.jpg")
	jpeg.Encode(tmpFile, img, nil)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectMimeType(tmpFile.Name())
	}
}

func BenchmarkValidateMimeType(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	tmpFile, _ := os.CreateTemp("", "bench-*.jpg")
	jpeg.Encode(tmpFile, img, nil)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateMimeType(tmpFile.Name(), AllowedImageTypes)
	}
}

func BenchmarkValidateMimeTypeFromBytes(b *testing.B) {
	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateMimeTypeFromBytes(jpegBytes, AllowedImageTypes)
	}
}
