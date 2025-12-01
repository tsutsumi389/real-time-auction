package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

var (
	ErrInvalidFilename     = errors.New("invalid filename")
	ErrPathTraversal       = errors.New("path traversal detected")
	ErrEmptyFilename       = errors.New("filename cannot be empty")
	ErrFilenameTooLong     = errors.New("filename is too long")
	ErrInvalidFileExtension = errors.New("invalid file extension")
)

const (
	MaxFilenameLength = 255 // ファイルシステムの一般的な制限
	MaxPathLength     = 4096
)

// SanitizeFilename はファイル名から危険な文字を除去し、安全な形式に変換する
func SanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", ErrEmptyFilename
	}

	// パストラバーサル攻撃を防ぐ
	if strings.Contains(filename, "..") {
		return "", ErrPathTraversal
	}

	// ベース名のみを取得（ディレクトリパスを除去）
	filename = filepath.Base(filename)

	// 特殊なファイル名をチェック（Unix/Windows）
	if filename == "." || filename == ".." || filename == "/" || filename == "\\" {
		return "", ErrInvalidFilename
	}

	// Windows予約語をチェック
	windowsReserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperFilename := strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for _, reserved := range windowsReserved {
		if upperFilename == reserved {
			return "", ErrInvalidFilename
		}
	}

	// 制御文字を除去
	filename = removeControlCharacters(filename)

	// 危険な文字を除去（許可: 英数字、ハイフン、アンダースコア、ドット）
	filename = sanitizeCharacters(filename)

	// 連続するドットを単一のドットに変換
	filename = regexp.MustCompile(`\.+`).ReplaceAllString(filename, ".")

	// 先頭と末尾のドット・スペースを除去
	filename = strings.Trim(filename, ". ")

	// 長さチェック
	if len(filename) > MaxFilenameLength {
		return "", ErrFilenameTooLong
	}

	if filename == "" {
		return "", ErrEmptyFilename
	}

	return filename, nil
}

// GenerateUniqueFilename はUUIDを使用してユニークなファイル名を生成する
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

// GeneratePrefixedFilename はプレフィックス付きのユニークなファイル名を生成する
func GeneratePrefixedFilename(prefix, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s_%s%s", prefix, uuid.New().String(), ext)
}

// SanitizeAndGenerateFilename はファイル名をサニタイズし、UUID付きのユニークな名前を生成する
func SanitizeAndGenerateFilename(originalFilename string) (string, error) {
	// まずサニタイズして拡張子を取得
	sanitized, err := SanitizeFilename(originalFilename)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(sanitized)
	return GenerateUniqueFilename(ext), nil
}

// ValidateFileExtension はファイル拡張子が許可リストに含まれているか確認する
func ValidateFileExtension(filename string, allowedExtensions []string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowed := range allowedExtensions {
		if !strings.HasPrefix(allowed, ".") {
			allowed = "." + allowed
		}
		if ext == strings.ToLower(allowed) {
			return nil
		}
	}

	return ErrInvalidFileExtension
}

// GetFileExtension はファイルの拡張子を返す（ドット含む）
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// GetFileNameWithoutExtension は拡張子を除いたファイル名を返す
func GetFileNameWithoutExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// removeControlCharacters は制御文字を除去する
func removeControlCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// sanitizeCharacters は許可された文字以外を除去する
func sanitizeCharacters(s string) string {
	// 許可: 英数字、日本語、ハイフン、アンダースコア、ドット、スペース
	re := regexp.MustCompile(`[^a-zA-Z0-9\p{Han}\p{Hiragana}\p{Katakana}\-_. ]`)
	return re.ReplaceAllString(s, "")
}

// SafeJoinPath は安全にパスを結合する（パストラバーサル防止）
func SafeJoinPath(base, filename string) (string, error) {
	// ファイル名をサニタイズ
	sanitized, err := SanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	// パスを結合
	fullPath := filepath.Join(base, sanitized)

	// 絶対パスに変換
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	// ベースディレクトリ外へのアクセスを防ぐ
	if !strings.HasPrefix(absPath, absBase) {
		return "", ErrPathTraversal
	}

	return fullPath, nil
}

// CreateTempFile は一時ファイルを作成する
func CreateTempFile(prefix, suffix string) (*os.File, error) {
	// プレフィックスをサニタイズ
	prefix = sanitizeCharacters(prefix)
	if prefix == "" {
		prefix = "upload"
	}

	// サフィックス（拡張子）をサニタイズ
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}

	// 一時ファイルを作成
	pattern := fmt.Sprintf("%s-*%s", prefix, suffix)
	return os.CreateTemp("", pattern)
}

// CopyFile はファイルをコピーする
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Sync to ensure data is written to disk
	return destFile.Sync()
}

// DeleteFile はファイルを削除する（存在しない場合はエラーを返さない）
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// DeleteFiles は複数のファイルを削除する
func DeleteFiles(filePaths []string) error {
	var errs []error
	for _, path := range filePaths {
		if err := DeleteFile(path); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to delete %d files", len(errs))
	}

	return nil
}

// FileExists はファイルが存在するか確認する
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetFileSize はファイルサイズを返す（バイト単位）
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// IsPathSafe はパスが安全か（パストラバーサルを含まないか）確認する
func IsPathSafe(path string) bool {
	// 元のパスに..が含まれているかチェック
	if strings.Contains(path, "..") {
		return false
	}
	// クリーン後のパスもチェック
	cleanPath := filepath.Clean(path)
	return !strings.Contains(cleanPath, "..")
}

// NormalizeExtension は拡張子を正規化する（小文字化、ドット付加）
func NormalizeExtension(ext string) string {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

// BuildStoragePath はストレージ用のパスを構築する
// 例: BuildStoragePath("items", "123", "original_abc.jpg") → "items/123/original_abc.jpg"
func BuildStoragePath(parts ...string) string {
	return filepath.Join(parts...)
}

// ExtractFilenameFromPath はパスからファイル名を抽出する
func ExtractFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// ExtractDirectoryFromPath はパスからディレクトリ部分を抽出する
func ExtractDirectoryFromPath(path string) string {
	return filepath.Dir(path)
}

// IsHiddenFile はファイルが隠しファイルか判定する（Unix形式）
func IsHiddenFile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

// ValidateFilePath はファイルパスの妥当性を検証する
func ValidateFilePath(path string) error {
	if path == "" {
		return ErrEmptyFilename
	}

	if len(path) > MaxPathLength {
		return errors.New("path is too long")
	}

	if !IsPathSafe(path) {
		return ErrPathTraversal
	}

	return nil
}

// GenerateMediaPath はメディアファイルのストレージパスを生成する
// 例: GenerateMediaPath(123, "original", "jpg") → "items/123/original_uuid.jpg"
func GenerateMediaPath(itemID int64, prefix, extension string) string {
	uniqueID := uuid.New().String()
	ext := NormalizeExtension(extension)
	filename := fmt.Sprintf("%s_%s%s", prefix, uniqueID, ext)
	return filepath.Join("items", fmt.Sprintf("%d", itemID), filename)
}

// CleanupTempFiles は一時ファイルを削除する
func CleanupTempFiles(files []string) {
	for _, file := range files {
		_ = DeleteFile(file)
	}
}

// EnsureDirectoryExists はディレクトリが存在することを保証する（なければ作成）
func EnsureDirectoryExists(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetTempDir は一時ディレクトリのパスを返す
func GetTempDir() string {
	return os.TempDir()
}

// ConvertToUnixPath はパスをUnix形式に変換する（Windowsパス対策）
func ConvertToUnixPath(path string) string {
	return filepath.ToSlash(path)
}

// HasAllowedExtension はファイルが許可された拡張子を持つか確認する
func HasAllowedExtension(filename string, allowedExtensions []string) bool {
	return ValidateFileExtension(filename, allowedExtensions) == nil
}

// StripExtension はファイル名から拡張子を除去する
func StripExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// ReplaceExtension はファイル名の拡張子を置き換える
func ReplaceExtension(filename, newExtension string) string {
	withoutExt := StripExtension(filename)
	newExt := NormalizeExtension(newExtension)
	return withoutExt + newExt
}
