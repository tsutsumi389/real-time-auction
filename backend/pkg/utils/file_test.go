package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:          "Normal filename",
			input:         "test.jpg",
			expected:      "test.jpg",
			expectedError: nil,
		},
		{
			name:          "Filename with spaces",
			input:         "my photo.jpg",
			expected:      "my photo.jpg",
			expectedError: nil,
		},
		{
			name:          "Filename with special characters",
			input:         "test@#$%.jpg",
			expected:      "test.jpg",
			expectedError: nil,
		},
		{
			name:          "Path traversal attempt",
			input:         "../../../etc/passwd",
			expected:      "",
			expectedError: ErrPathTraversal,
		},
		{
			name:          "Empty filename",
			input:         "",
			expected:      "",
			expectedError: ErrEmptyFilename,
		},
		{
			name:          "Windows reserved name",
			input:         "CON.txt",
			expected:      "",
			expectedError: ErrInvalidFilename,
		},
		{
			name:          "Multiple dots",
			input:         "file...jpg",
			expected:      "",
			expectedError: ErrPathTraversal,
		},
		{
			name:          "Leading and trailing dots",
			input:         ".test.jpg.",
			expected:      "test.jpg",
			expectedError: nil,
		},
		{
			name:          "Japanese filename",
			input:         "テスト画像.jpg",
			expected:      "テスト画像.jpg",
			expectedError: nil,
		},
		{
			name:          "With directory path",
			input:         "/path/to/file.jpg",
			expected:      "file.jpg",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeFilename(tt.input)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if err == nil && result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		checkExt string
	}{
		{
			name:     "JPEG file",
			input:    "test.jpg",
			checkExt: ".jpg",
		},
		{
			name:     "PNG file",
			input:    "image.png",
			checkExt: ".png",
		},
		{
			name:     "No extension",
			input:    "file",
			checkExt: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateUniqueFilename(tt.input)

			// UUIDフォーマットをチェック（36文字 + 拡張子）
			ext := filepath.Ext(result)
			if ext != tt.checkExt {
				t.Errorf("Expected extension '%s', got '%s'", tt.checkExt, ext)
			}

			// UUID部分の長さチェック（拡張子を除く）
			nameWithoutExt := strings.TrimSuffix(result, ext)
			if len(nameWithoutExt) != 36 {
				t.Errorf("Expected UUID length 36, got %d", len(nameWithoutExt))
			}

			// 2回呼び出して異なる値が返ることを確認
			result2 := GenerateUniqueFilename(tt.input)
			if result == result2 {
				t.Error("Generated filenames should be unique")
			}
		})
	}
}

func TestGeneratePrefixedFilename(t *testing.T) {
	prefix := "thumb"
	filename := "test.jpg"

	result := GeneratePrefixedFilename(prefix, filename)

	if !strings.HasPrefix(result, prefix+"_") {
		t.Errorf("Expected prefix '%s_', got '%s'", prefix, result)
	}

	if !strings.HasSuffix(result, ".jpg") {
		t.Errorf("Expected suffix '.jpg', got '%s'", result)
	}
}

func TestValidateFileExtension(t *testing.T) {
	tests := []struct {
		name              string
		filename          string
		allowedExtensions []string
		expectedError     error
	}{
		{
			name:              "Valid extension",
			filename:          "test.jpg",
			allowedExtensions: []string{".jpg", ".png"},
			expectedError:     nil,
		},
		{
			name:              "Valid extension without dot",
			filename:          "test.jpg",
			allowedExtensions: []string{"jpg", "png"},
			expectedError:     nil,
		},
		{
			name:              "Invalid extension",
			filename:          "test.pdf",
			allowedExtensions: []string{".jpg", ".png"},
			expectedError:     ErrInvalidFileExtension,
		},
		{
			name:              "Case insensitive",
			filename:          "test.JPG",
			allowedExtensions: []string{".jpg"},
			expectedError:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileExtension(tt.filename, tt.allowedExtensions)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.jpg", ".jpg"},
		{"file.tar.gz", ".gz"},
		{"noext", ""},
		{".hidden", ".hidden"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := GetFileExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetFileNameWithoutExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.jpg", "test"},
		{"file.tar.gz", "file.tar"},
		{"noext", "noext"},
		{".hidden", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := GetFileNameWithoutExtension(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSafeJoinPath(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "test-safejoin-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name          string
		base          string
		filename      string
		expectError   bool
		checkContains string
	}{
		{
			name:          "Normal join",
			base:          tmpDir,
			filename:      "test.jpg",
			expectError:   false,
			checkContains: "test.jpg",
		},
		{
			name:        "Path traversal attempt",
			base:        tmpDir,
			filename:    "../../../etc/passwd",
			expectError: true,
		},
		{
			name:          "Sanitized filename",
			base:          tmpDir,
			filename:      "test@#$.jpg",
			expectError:   false,
			checkContains: "test.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeJoinPath(tt.base, tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !strings.Contains(result, tt.checkContains) {
					t.Errorf("Expected result to contain '%s', got '%s'", tt.checkContains, result)
				}
			}
		})
	}
}

func TestCreateTempFile(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		suffix string
	}{
		{
			name:   "With prefix and suffix",
			prefix: "test",
			suffix: ".jpg",
		},
		{
			name:   "Empty prefix",
			prefix: "",
			suffix: ".png",
		},
		{
			name:   "No suffix",
			prefix: "test",
			suffix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := CreateTempFile(tt.prefix, tt.suffix)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// ファイルが存在することを確認
			if !FileExists(tmpFile.Name()) {
				t.Error("Temp file was not created")
			}

			// 拡張子を確認
			if tt.suffix != "" {
				ext := filepath.Ext(tmpFile.Name())
				expectedExt := tt.suffix
				if !strings.HasPrefix(expectedExt, ".") {
					expectedExt = "." + expectedExt
				}
				if ext != expectedExt {
					t.Errorf("Expected extension '%s', got '%s'", expectedExt, ext)
				}
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	// ソースファイルを作成
	srcFile, err := os.CreateTemp("", "src-*.txt")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	content := "Test content"
	if _, err := srcFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}
	srcFile.Close()

	// コピー先のパス
	dstPath := srcFile.Name() + ".copy"
	defer os.Remove(dstPath)

	// コピー実行
	if err := CopyFile(srcFile.Name(), dstPath); err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}

	// コピーしたファイルの内容を確認
	copiedContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copiedContent) != content {
		t.Errorf("Expected content '%s', got '%s'", content, string(copiedContent))
	}
}

func TestDeleteFile(t *testing.T) {
	// テストファイルを作成
	tmpFile, err := os.CreateTemp("", "test-delete-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	filePath := tmpFile.Name()

	// 削除実行
	if err := DeleteFile(filePath); err != nil {
		t.Errorf("Failed to delete file: %v", err)
	}

	// ファイルが削除されたことを確認
	if FileExists(filePath) {
		t.Error("File was not deleted")
	}

	// 存在しないファイルを削除してもエラーにならないことを確認
	if err := DeleteFile(filePath); err != nil {
		t.Errorf("Deleting non-existent file should not error: %v", err)
	}
}

func TestDeleteFiles(t *testing.T) {
	// 複数のテストファイルを作成
	var files []string
	for i := 0; i < 3; i++ {
		tmpFile, err := os.CreateTemp("", "test-delete-*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpFile.Close()
		files = append(files, tmpFile.Name())
	}

	// 削除実行
	if err := DeleteFiles(files); err != nil {
		t.Errorf("Failed to delete files: %v", err)
	}

	// すべてのファイルが削除されたことを確認
	for _, file := range files {
		if FileExists(file) {
			t.Errorf("File %s was not deleted", file)
		}
	}
}

func TestFileExists(t *testing.T) {
	// 存在するファイル
	tmpFile, err := os.CreateTemp("", "test-exists-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if !FileExists(tmpFile.Name()) {
		t.Error("FileExists should return true for existing file")
	}

	// 存在しないファイル
	if FileExists("/non/existent/file.txt") {
		t.Error("FileExists should return false for non-existent file")
	}

	// ディレクトリ
	tmpDir, _ := os.MkdirTemp("", "test-dir-")
	defer os.RemoveAll(tmpDir)

	if FileExists(tmpDir) {
		t.Error("FileExists should return false for directories")
	}
}

func TestGetFileSize(t *testing.T) {
	content := "Hello, World!"
	tmpFile, err := os.CreateTemp("", "test-size-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	size, err := GetFileSize(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to get file size: %v", err)
	}

	expectedSize := int64(len(content))
	if size != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, size)
	}
}

func TestIsPathSafe(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"normal/path/file.txt", true},
		{"../../../etc/passwd", false},
		{"path/../file.txt", false},
		{"/absolute/path/file.txt", true},
		{"./relative/path", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsPathSafe(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNormalizeExtension(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"jpg", ".jpg"},
		{".jpg", ".jpg"},
		{"JPG", ".jpg"},
		{".PNG", ".png"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeExtension(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestBuildStoragePath(t *testing.T) {
	result := BuildStoragePath("items", "123", "test.jpg")
	expected := filepath.Join("items", "123", "test.jpg")

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestExtractFilenameFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.jpg", "file.jpg"},
		{"file.jpg", "file.jpg"},
		{"/path/to/", "to"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ExtractFilenameFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsHiddenFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{".hidden", true},
		{"normal.txt", false},
		{".gitignore", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := IsHiddenFile(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateMediaPath(t *testing.T) {
	result := GenerateMediaPath(123, "original", "jpg")

	// パスのフォーマットを確認
	if !strings.Contains(result, "items/123/") {
		t.Errorf("Expected path to contain 'items/123/', got '%s'", result)
	}

	if !strings.HasPrefix(filepath.Base(result), "original_") {
		t.Errorf("Expected filename to start with 'original_', got '%s'", filepath.Base(result))
	}

	if !strings.HasSuffix(result, ".jpg") {
		t.Errorf("Expected path to end with '.jpg', got '%s'", result)
	}
}

func TestReplaceExtension(t *testing.T) {
	tests := []struct {
		filename     string
		newExtension string
		expected     string
	}{
		{"test.png", ".jpg", "test.jpg"},
		{"test.png", "jpg", "test.jpg"},
		{"file", ".txt", "file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := ReplaceExtension(tt.filename, tt.newExtension)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-ensure-dir")
	defer os.RemoveAll(tmpDir)

	// ディレクトリを作成
	if err := EnsureDirectoryExists(tmpDir); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// ディレクトリが存在することを確認
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// 既存のディレクトリに対して再度実行してもエラーにならないことを確認
	if err := EnsureDirectoryExists(tmpDir); err != nil {
		t.Errorf("EnsureDirectoryExists should not error for existing directory: %v", err)
	}
}

// Benchmark tests
func BenchmarkSanitizeFilename(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SanitizeFilename("test@#$%^&*().jpg")
	}
}

func BenchmarkGenerateUniqueFilename(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUniqueFilename("test.jpg")
	}
}

func BenchmarkIsPathSafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPathSafe("normal/path/to/file.txt")
	}
}
