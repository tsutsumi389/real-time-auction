package middleware

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLoadUploadConfig(t *testing.T) {
	// 環境変数を設定
	os.Setenv("MAX_IMAGE_SIZE", "10485760")  // 10MB
	os.Setenv("MAX_VIDEO_SIZE", "209715200") // 200MB
	defer os.Unsetenv("MAX_IMAGE_SIZE")
	defer os.Unsetenv("MAX_VIDEO_SIZE")

	config := LoadUploadConfig()

	assert.Equal(t, int64(10485760), config.MaxImageSize)
	assert.Equal(t, int64(209715200), config.MaxVideoSize)
}

func TestLoadUploadConfigDefaults(t *testing.T) {
	// 環境変数がない場合
	os.Unsetenv("MAX_IMAGE_SIZE")
	os.Unsetenv("MAX_VIDEO_SIZE")

	config := LoadUploadConfig()

	assert.Equal(t, int64(DefaultMaxImageSize), config.MaxImageSize)
	assert.Equal(t, int64(DefaultMaxVideoSize), config.MaxVideoSize)
}

func TestLimitFileSize(t *testing.T) {
	tests := []struct {
		name           string
		maxSize        int64
		contentLength  int64
		expectedStatus int
	}{
		{
			name:           "Within limit",
			maxSize:        1024,
			contentLength:  512,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceeds limit",
			maxSize:        1024,
			contentLength:  2048,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "Exactly at limit",
			maxSize:        1024,
			contentLength:  1024,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(LimitFileSize(tt.maxSize))
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("POST", "/upload", bytes.NewReader(make([]byte, tt.contentLength)))
			req.Header.Set("Content-Length", string(rune(tt.contentLength)))
			req.ContentLength = tt.contentLength

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLimitImageUpload(t *testing.T) {
	router := gin.New()
	router.Use(LimitImageUpload())
	router.POST("/upload", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// デフォルトサイズ内
	req := httptest.NewRequest("POST", "/upload", bytes.NewReader(make([]byte, 1024)))
	req.ContentLength = 1024
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLimitMediaUpload(t *testing.T) {
	tests := []struct {
		name           string
		mediaType      string
		contentLength  int64
		expectedStatus int
	}{
		{
			name:           "Image within limit",
			mediaType:      "image",
			contentLength:  1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Video within limit",
			mediaType:      "video",
			contentLength:  1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unknown type defaults to image limit",
			mediaType:      "unknown",
			contentLength:  1024,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(LimitMediaUpload())
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			writer.WriteField("media_type", tt.mediaType)
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.ContentLength = tt.contentLength

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestValidateMultipartForm(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Valid multipart/form-data",
			contentType:    "multipart/form-data; boundary=----WebKitFormBoundary",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid content type",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty content type",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ValidateMultipartForm())
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("POST", "/upload", nil)
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireFileUpload(t *testing.T) {
	tests := []struct {
		name           string
		includeFile    bool
		fileSize       int64
		expectedStatus int
	}{
		{
			name:           "File uploaded",
			includeFile:    true,
			fileSize:       1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No file uploaded",
			includeFile:    false,
			fileSize:       0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty file",
			includeFile:    true,
			fileSize:       0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RequireFileUpload("file"))
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.includeFile {
				part, _ := writer.CreateFormFile("file", "test.txt")
				if tt.fileSize > 0 {
					io.CopyN(part, bytes.NewReader(make([]byte, tt.fileSize)), tt.fileSize)
				}
			}
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name           string
		fileSize       int64
		maxSize        int64
		expectedStatus int
	}{
		{
			name:           "Within limit",
			fileSize:       512,
			maxSize:        1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceeds limit",
			fileSize:       2048,
			maxSize:        1024,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "At limit",
			fileSize:       1024,
			maxSize:        1024,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ValidateFileSize("file", tt.maxSize))
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.txt")
			io.CopyN(part, bytes.NewReader(make([]byte, tt.fileSize)), tt.fileSize)
			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestValidateContentLength(t *testing.T) {
	tests := []struct {
		name           string
		contentLength  int64
		maxSize        int64
		expectedStatus int
	}{
		{
			name:           "Within limit",
			contentLength:  512,
			maxSize:        1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Exceeds limit",
			contentLength:  2048,
			maxSize:        1024,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "Negative content length",
			contentLength:  -1,
			maxSize:        1024,
			expectedStatus: http.StatusLengthRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ValidateContentLength(tt.maxSize))
			router.POST("/upload", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("POST", "/upload", bytes.NewReader(make([]byte, 0)))
			req.ContentLength = tt.contentLength

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{512, "512 bytes"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{5 * 1024 * 1024, "5.00 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMediaTypeName(t *testing.T) {
	tests := []struct {
		mediaType string
		expected  string
	}{
		{"image", "画像"},
		{"video", "動画"},
		{"unknown", "ファイル"},
		{"", "ファイル"},
	}

	for _, tt := range tests {
		t.Run(tt.mediaType, func(t *testing.T) {
			result := getMediaTypeName(tt.mediaType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCombineUploadMiddleware(t *testing.T) {
	middlewares := CombineUploadMiddleware("file")

	// ミドルウェアが正しい数だけ返されることを確認
	assert.Equal(t, 4, len(middlewares))
}

func TestImageUploadMiddleware(t *testing.T) {
	middlewares := ImageUploadMiddleware("file")

	// ミドルウェアが正しい数だけ返されることを確認
	assert.Equal(t, 4, len(middlewares))
}

func TestVideoUploadMiddleware(t *testing.T) {
	middlewares := VideoUploadMiddleware("file")

	// ミドルウェアが正しい数だけ返されることを確認
	assert.Equal(t, 4, len(middlewares))
}

func TestSetMaxMultipartMemory(t *testing.T) {
	router := gin.New()
	router.Use(SetMaxMultipartMemory(8 * 1024 * 1024)) // 8MB
	router.POST("/upload", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUploadConfig(t *testing.T) {
	router := gin.New()
	router.GET("/config", GetUploadConfig)

	req := httptest.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "max_image_size")
	assert.Contains(t, w.Body.String(), "max_video_size")
	assert.Contains(t, w.Body.String(), "max_file_size")
}

// Integration test
func TestUploadMiddlewareIntegration(t *testing.T) {
	router := gin.New()

	// 画像アップロード用エンドポイント
	router.POST("/upload/image", append(ImageUploadMiddleware("file"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})...)

	// テストデータを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write(make([]byte, 1024)) // 1KB
	writer.Close()

	req := httptest.NewRequest("POST", "/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadMiddlewareIntegrationFailures(t *testing.T) {
	router := gin.New()

	router.POST("/upload/image", append(ImageUploadMiddleware("file"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})...)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
	}{
		{
			name: "Missing file",
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()

				req := httptest.NewRequest("POST", "/upload/image", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong content type",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("POST", "/upload/image", bytes.NewReader([]byte{}))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Benchmark tests
func BenchmarkLimitFileSize(b *testing.B) {
	router := gin.New()
	router.Use(LimitFileSize(5 * 1024 * 1024))
	router.POST("/upload", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/upload", bytes.NewReader(make([]byte, 1024)))
	req.ContentLength = 1024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkValidateMultipartForm(b *testing.B) {
	router := gin.New()
	router.Use(ValidateMultipartForm())
	router.POST("/upload", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
