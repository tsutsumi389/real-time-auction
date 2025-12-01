package service

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ImageProcessor 画像処理サービス
type ImageProcessor struct {
	maxWidth      int // 最大幅
	maxHeight     int // 最大高さ
	thumbnailSize int // サムネイルサイズ（正方形）
	quality       int // JPEG品質（1-100）
}

// NewImageProcessor 新しいImageProcessorを作成
func NewImageProcessor(maxWidth, maxHeight, thumbnailSize, quality int) *ImageProcessor {
	return &ImageProcessor{
		maxWidth:      maxWidth,
		maxHeight:     maxHeight,
		thumbnailSize: thumbnailSize,
		quality:       quality,
	}
}

// ProcessedImage 処理済み画像の情報
type ProcessedImage struct {
	OriginalPath  string // オリジナル画像のパス
	ThumbnailPath string // サムネイル画像のパス
}

// ProcessImage 画像をリサイズしてJPEG形式に変換
// srcPath: 元画像のパス
// 返り値: 処理済み画像情報、エラー
func (ip *ImageProcessor) ProcessImage(srcPath string) (*ProcessedImage, error) {
	// 元画像を開く
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	// リサイズ（アスペクト比を維持）
	resizedImg := ip.resizeImage(srcImg)

	// JPEG形式で保存
	originalPath, err := ip.saveAsJPEG(resizedImg, "original")
	if err != nil {
		return nil, fmt.Errorf("failed to save original image: %w", err)
	}

	// サムネイル生成
	thumbnailPath, err := ip.GenerateThumbnail(srcImg)
	if err != nil {
		// オリジナル画像は保存済みなので削除
		os.Remove(originalPath)
		return nil, fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	return &ProcessedImage{
		OriginalPath:  originalPath,
		ThumbnailPath: thumbnailPath,
	}, nil
}

// resizeImage 画像をリサイズ（アスペクト比維持）
func (ip *ImageProcessor) resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// すでに最大サイズ以下の場合はそのまま
	if width <= ip.maxWidth && height <= ip.maxHeight {
		return img
	}

	// アスペクト比を維持してリサイズ
	return imaging.Fit(img, ip.maxWidth, ip.maxHeight, imaging.Lanczos)
}

// GenerateThumbnail サムネイル画像を生成（正方形、中央クロップ）
func (ip *ImageProcessor) GenerateThumbnail(srcImg image.Image) (string, error) {
	// 正方形にクロップ（中央）
	thumbnail := imaging.Fill(srcImg, ip.thumbnailSize, ip.thumbnailSize, imaging.Center, imaging.Lanczos)

	// JPEG形式で保存
	thumbnailPath, err := ip.saveAsJPEG(thumbnail, "thumb")
	if err != nil {
		return "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return thumbnailPath, nil
}

// saveAsJPEG 画像をJPEG形式で保存
func (ip *ImageProcessor) saveAsJPEG(img image.Image, prefix string) (string, error) {
	// 一時ディレクトリにUUIDでファイル名生成
	filename := fmt.Sprintf("%s_%s.jpg", prefix, uuid.New().String())
	outputPath := filepath.Join(os.TempDir(), filename)

	// JPEGエンコード（imaging.Saveが品質オプションをサポート）
	err := imaging.Save(img, outputPath, imaging.JPEGQuality(ip.quality))
	if err != nil {
		return "", fmt.Errorf("failed to save JPEG: %w", err)
	}

	return outputPath, nil
}

// CleanupTempFiles 一時ファイルを削除
func (ip *ImageProcessor) CleanupTempFiles(paths []string) {
	for _, path := range paths {
		if path != "" {
			os.Remove(path)
		}
	}
}

// ValidateImageFile 画像ファイルのバリデーション（基本的なチェック）
func (ip *ImageProcessor) ValidateImageFile(filePath string) error {
	// ファイルを開いて画像として認識できるかチェック
	img, err := imaging.Open(filePath)
	if err != nil {
		return fmt.Errorf("invalid image file: %w", err)
	}

	// 画像サイズチェック
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return fmt.Errorf("invalid image dimensions: %dx%d", bounds.Dx(), bounds.Dy())
	}

	return nil
}
