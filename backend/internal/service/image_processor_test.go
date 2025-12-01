package service

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestImage テスト用画像を作成
func createTestImage(width, height int) (string, error) {
	// カラー画像を作成
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// グラデーション効果を追加
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// 一時ファイルに保存
	tmpFile, err := os.CreateTemp("", "test_image_*.png")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if err := png.Encode(tmpFile, img); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestNewImageProcessor(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	assert.Equal(t, 1920, processor.maxWidth)
	assert.Equal(t, 1080, processor.maxHeight)
	assert.Equal(t, 300, processor.thumbnailSize)
	assert.Equal(t, 80, processor.quality)
}

func TestProcessImage(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	// 大きな画像を作成（2400x1600）
	testImagePath, err := createTestImage(2400, 1600)
	require.NoError(t, err)
	defer os.Remove(testImagePath)

	// 画像処理実行
	result, err := processor.ProcessImage(testImagePath)
	require.NoError(t, err)
	require.NotNil(t, result)

	defer func() {
		os.Remove(result.OriginalPath)
		os.Remove(result.ThumbnailPath)
	}()

	// オリジナル画像が存在することを確認
	assert.FileExists(t, result.OriginalPath)
	assert.Contains(t, result.OriginalPath, "original_")
	assert.Contains(t, result.OriginalPath, ".jpg")

	// サムネイルが存在することを確認
	assert.FileExists(t, result.ThumbnailPath)
	assert.Contains(t, result.ThumbnailPath, "thumb_")
	assert.Contains(t, result.ThumbnailPath, ".jpg")

	// オリジナル画像のサイズチェック（リサイズされているか）
	originalImg, err := imaging.Open(result.OriginalPath)
	require.NoError(t, err)
	bounds := originalImg.Bounds()

	// 最大サイズ以下になっているか
	assert.LessOrEqual(t, bounds.Dx(), 1920)
	assert.LessOrEqual(t, bounds.Dy(), 1080)

	// アスペクト比が維持されているか（元画像は2400x1600 = 3:2）
	// リサイズ後も3:2になっているはず（1620x1080）
	expectedWidth := 1620 // 1080 * 3 / 2
	assert.Equal(t, expectedWidth, bounds.Dx())
	assert.Equal(t, 1080, bounds.Dy())
}

func TestResizeImage(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	tests := []struct {
		name           string
		inputWidth     int
		inputHeight    int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "大きな画像（横長）",
			inputWidth:     2400,
			inputHeight:    1600,
			expectedWidth:  1620,
			expectedHeight: 1080,
		},
		{
			name:           "大きな画像（縦長）",
			inputWidth:     1600,
			inputHeight:    2400,
			expectedWidth:  720,
			expectedHeight: 1080,
		},
		{
			name:           "小さな画像（リサイズ不要）",
			inputWidth:     800,
			inputHeight:    600,
			expectedWidth:  800,
			expectedHeight: 600,
		},
		{
			name:           "最大サイズぴったり",
			inputWidth:     1920,
			inputHeight:    1080,
			expectedWidth:  1920,
			expectedHeight: 1080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト画像作成
			img := image.NewRGBA(image.Rect(0, 0, tt.inputWidth, tt.inputHeight))

			// リサイズ実行
			resized := processor.resizeImage(img)
			bounds := resized.Bounds()

			assert.Equal(t, tt.expectedWidth, bounds.Dx())
			assert.Equal(t, tt.expectedHeight, bounds.Dy())
		})
	}
}

func TestGenerateThumbnail(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	// テスト画像作成（1200x800）
	testImagePath, err := createTestImage(1200, 800)
	require.NoError(t, err)
	defer os.Remove(testImagePath)

	img, err := imaging.Open(testImagePath)
	require.NoError(t, err)

	// サムネイル生成
	thumbnailPath, err := processor.GenerateThumbnail(img)
	require.NoError(t, err)
	defer os.Remove(thumbnailPath)

	// ファイル存在チェック
	assert.FileExists(t, thumbnailPath)
	assert.Contains(t, thumbnailPath, "thumb_")
	assert.Contains(t, thumbnailPath, ".jpg")

	// サムネイルサイズチェック（300x300の正方形）
	thumbnail, err := imaging.Open(thumbnailPath)
	require.NoError(t, err)
	bounds := thumbnail.Bounds()

	assert.Equal(t, 300, bounds.Dx())
	assert.Equal(t, 300, bounds.Dy())
}

func TestSaveAsJPEG(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	// テスト画像作成
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))

	// JPEG保存
	outputPath, err := processor.saveAsJPEG(img, "test")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// ファイル存在チェック
	assert.FileExists(t, outputPath)
	assert.Contains(t, outputPath, "test_")
	assert.Contains(t, outputPath, ".jpg")

	// ファイル拡張子チェック
	ext := filepath.Ext(outputPath)
	assert.Equal(t, ".jpg", ext)
}

func TestCleanupTempFiles(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	// テスト用一時ファイルを作成
	tmpFile1, err := os.CreateTemp("", "cleanup_test_1_*.txt")
	require.NoError(t, err)
	tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "cleanup_test_2_*.txt")
	require.NoError(t, err)
	tmpFile2.Close()

	paths := []string{tmpFile1.Name(), tmpFile2.Name(), ""}

	// ファイルが存在することを確認
	assert.FileExists(t, tmpFile1.Name())
	assert.FileExists(t, tmpFile2.Name())

	// クリーンアップ実行
	processor.CleanupTempFiles(paths)

	// ファイルが削除されたことを確認
	assert.NoFileExists(t, tmpFile1.Name())
	assert.NoFileExists(t, tmpFile2.Name())
}

func TestValidateImageFile(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	tests := []struct {
		name      string
		setupFunc func() (string, error)
		wantError bool
	}{
		{
			name: "正常な画像ファイル",
			setupFunc: func() (string, error) {
				return createTestImage(800, 600)
			},
			wantError: false,
		},
		{
			name: "存在しないファイル",
			setupFunc: func() (string, error) {
				return "/nonexistent/file.png", nil
			},
			wantError: true,
		},
		{
			name: "画像ではないファイル",
			setupFunc: func() (string, error) {
				tmpFile, err := os.CreateTemp("", "invalid_*.txt")
				if err != nil {
					return "", err
				}
				tmpFile.WriteString("This is not an image")
				tmpFile.Close()
				return tmpFile.Name(), nil
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}
			if filePath != "/nonexistent/file.png" {
				defer os.Remove(filePath)
			}

			err = processor.ValidateImageFile(filePath)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessImage_InvalidFile(t *testing.T) {
	processor := NewImageProcessor(1920, 1080, 300, 80)

	// 存在しないファイル
	result, err := processor.ProcessImage("/nonexistent/file.png")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to open image")
}
