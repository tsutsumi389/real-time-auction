package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"github.com/tsutsumi389/real-time-auction/internal/service"
	"github.com/tsutsumi389/real-time-auction/pkg/utils"
)

// MediaHandler handles media upload and management HTTP requests
type MediaHandler struct {
	mediaRepo       *repository.ItemMediaRepository
	itemRepo        *repository.ItemRepository
	storageService  service.StorageService
	imageProcessor  *service.ImageProcessor
	storageBucket   string
	maxImageSize    int64
	maxVideoSize    int64
	maxImagesPerItem int
	maxVideosPerItem int
}

// MediaHandlerConfig holds configuration for MediaHandler
type MediaHandlerConfig struct {
	MediaRepo        *repository.ItemMediaRepository
	ItemRepo         *repository.ItemRepository
	StorageService   service.StorageService
	ImageProcessor   *service.ImageProcessor
	StorageBucket    string
	MaxImageSize     int64
	MaxVideoSize     int64
	MaxImagesPerItem int
	MaxVideosPerItem int
}

// NewMediaHandler creates a new MediaHandler instance
func NewMediaHandler(config *MediaHandlerConfig) *MediaHandler {
	return &MediaHandler{
		mediaRepo:        config.MediaRepo,
		itemRepo:         config.ItemRepo,
		storageService:   config.StorageService,
		imageProcessor:   config.ImageProcessor,
		storageBucket:    config.StorageBucket,
		maxImageSize:     config.MaxImageSize,
		maxVideoSize:     config.MaxVideoSize,
		maxImagesPerItem: config.MaxImagesPerItem,
		maxVideosPerItem: config.MaxVideosPerItem,
	}
}

// UploadMedia handles POST /api/admin/items/:id/media
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Validate UUID format
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Verify item exists
	item, err := h.itemRepo.FindItemDetailByID(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Item not found",
		})
		return
	}

	// Get media type from form
	mediaTypeStr := c.PostForm("media_type")
	if mediaTypeStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "media_type is required (image or video)",
		})
		return
	}

	mediaType := domain.MediaType(mediaTypeStr)
	if mediaType != domain.MediaTypeImage && mediaType != domain.MediaTypeVideo {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid media_type: must be 'image' or 'video'",
		})
		return
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "No file uploaded",
		})
		return
	}

	// Validate file size
	maxSize := h.maxImageSize
	if mediaType == domain.MediaTypeVideo {
		maxSize = h.maxVideoSize
	}
	if fileHeader.Size > maxSize {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("File size exceeds limit (%d MB)", maxSize/1024/1024),
		})
		return
	}

	// Check media count limits
	if mediaType == domain.MediaTypeImage {
		count, err := h.mediaRepo.CountByItemIDAndType(itemUUID, domain.MediaTypeImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
			return
		}
		if count >= int64(h.maxImagesPerItem) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: fmt.Sprintf("Maximum %d images allowed per item", h.maxImagesPerItem),
			})
			return
		}
	} else {
		count, err := h.mediaRepo.CountByItemIDAndType(itemUUID, domain.MediaTypeVideo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
			return
		}
		if count >= int64(h.maxVideosPerItem) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: fmt.Sprintf("Maximum %d videos allowed per item", h.maxVideosPerItem),
			})
			return
		}
	}

	// Validate MIME type
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to open uploaded file",
		})
		return
	}
	defer file.Close()

	mimeType, err := utils.DetectMimeTypeFromReader(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Failed to detect file type",
		})
		return
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process file",
		})
		return
	}

	// Validate MIME type
	allowedTypes := utils.AllowedImageTypes
	if mediaType == domain.MediaTypeVideo {
		allowedTypes = utils.AllowedVideoTypes
	}
	if !utils.IsMimeTypeAllowed(mimeType, allowedTypes) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("Invalid file type: %s", mimeType),
		})
		return
	}

	// Save uploaded file to temporary location
	tempDir := os.TempDir()
	tempFilename := fmt.Sprintf("upload_%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
	tempPath := filepath.Join(tempDir, tempFilename)

	if err := c.SaveUploadedFile(fileHeader, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to save uploaded file",
		})
		return
	}

	// Handle image processing or direct upload
	var mediaURL, thumbnailURL string
	var cleanupPaths []string
	cleanupPaths = append(cleanupPaths, tempPath)

	if mediaType == domain.MediaTypeImage {
		// Process image (resize, convert to JPEG, generate thumbnail)
		result, err := h.uploadImageWithProcessing(c, itemUUID, tempPath, mimeType)
		if err != nil {
			// Cleanup temp files
			os.Remove(tempPath)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		mediaURL = result.OriginalURL
		thumbnailURL = result.ThumbnailURL
		cleanupPaths = append(cleanupPaths, result.CleanupPaths...)
	} else {
		// Upload video directly without processing
		objectName := fmt.Sprintf("items/%s/video_%s%s", itemUUID.String(), uuid.New().String(), filepath.Ext(fileHeader.Filename))

		file, err := os.Open(tempPath)
		if err != nil {
			os.Remove(tempPath)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to open file for upload",
			})
			return
		}
		defer file.Close()

		mediaURL, err = h.storageService.Upload(c.Request.Context(), h.storageBucket, objectName, file, fileHeader.Size, mimeType)
		if err != nil {
			os.Remove(tempPath)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to upload video to storage",
			})
			return
		}
	}

	// Get next display order
	displayOrder, err := h.mediaRepo.GetNextDisplayOrder(itemUUID)
	if err != nil {
		// Cleanup uploaded files from storage
		h.cleanupStorageFiles(c, mediaURL, thumbnailURL)
		for _, path := range cleanupPaths {
			os.Remove(path)
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to determine display order",
		})
		return
	}

	// Save to database
	media := &domain.ItemMedia{
		ItemID:       itemUUID,
		MediaType:    mediaType,
		URL:          mediaURL,
		DisplayOrder: displayOrder,
	}

	if thumbnailURL != "" {
		media.ThumbnailURL = &thumbnailURL
	}

	if err := h.mediaRepo.Create(media); err != nil {
		// Cleanup uploaded files from storage
		h.cleanupStorageFiles(c, mediaURL, thumbnailURL)
		for _, path := range cleanupPaths {
			os.Remove(path)
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to save media to database",
		})
		return
	}

	// Cleanup temporary files
	for _, path := range cleanupPaths {
		os.Remove(path)
	}

	// Return success response
	response := &domain.UploadMediaResponse{
		ID:           media.ID,
		ItemID:       media.ItemID,
		MediaType:    media.MediaType,
		URL:          media.URL,
		ThumbnailURL: media.ThumbnailURL,
		DisplayOrder: media.DisplayOrder,
		CreatedAt:    media.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// uploadImageResult holds the result of image upload with processing
type uploadImageResult struct {
	OriginalURL   string
	ThumbnailURL  string
	CleanupPaths  []string
}

// uploadImageWithProcessing processes and uploads an image
func (h *MediaHandler) uploadImageWithProcessing(c *gin.Context, itemID uuid.UUID, srcPath string, mimeType string) (*uploadImageResult, error) {
	// Process image (resize + thumbnail)
	processed, err := h.imageProcessor.ProcessImage(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	cleanupPaths := []string{processed.OriginalPath, processed.ThumbnailPath}

	// Upload original (processed) image
	originalObjectName := fmt.Sprintf("items/%s/original_%s.jpg", itemID.String(), uuid.New().String())
	originalFile, err := os.Open(processed.OriginalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open processed image: %w", err)
	}
	defer originalFile.Close()

	originalStat, err := originalFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	originalURL, err := h.storageService.Upload(c.Request.Context(), h.storageBucket, originalObjectName, originalFile, originalStat.Size(), "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("failed to upload original image: %w", err)
	}

	// Upload thumbnail
	thumbnailObjectName := fmt.Sprintf("items/%s/thumb_%s.jpg", itemID.String(), uuid.New().String())
	thumbnailFile, err := os.Open(processed.ThumbnailPath)
	if err != nil {
		// Cleanup original from storage
		h.storageService.Delete(c.Request.Context(), h.storageBucket, originalObjectName)
		return nil, fmt.Errorf("failed to open thumbnail: %w", err)
	}
	defer thumbnailFile.Close()

	thumbnailStat, err := thumbnailFile.Stat()
	if err != nil {
		h.storageService.Delete(c.Request.Context(), h.storageBucket, originalObjectName)
		return nil, fmt.Errorf("failed to get thumbnail info: %w", err)
	}

	thumbnailURL, err := h.storageService.Upload(c.Request.Context(), h.storageBucket, thumbnailObjectName, thumbnailFile, thumbnailStat.Size(), "image/jpeg")
	if err != nil {
		// Cleanup original from storage
		h.storageService.Delete(c.Request.Context(), h.storageBucket, originalObjectName)
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return &uploadImageResult{
		OriginalURL:  originalURL,
		ThumbnailURL: thumbnailURL,
		CleanupPaths: cleanupPaths,
	}, nil
}

// cleanupStorageFiles removes files from storage
func (h *MediaHandler) cleanupStorageFiles(c *gin.Context, urls ...string) {
	for _, url := range urls {
		if url == "" {
			continue
		}
		// Extract object name from URL
		objectName := h.extractObjectName(url)
		if objectName != "" {
			if err := h.storageService.Delete(c.Request.Context(), h.storageBucket, objectName); err != nil {
				log.Printf("Failed to cleanup storage file %s: %v", objectName, err)
			}
		}
	}
}

// extractObjectName extracts the object key from a full URL
func (h *MediaHandler) extractObjectName(url string) string {
	// Expected format: http://localhost:9000/bucket-name/items/uuid/filename.jpg
	// We need to extract: items/uuid/filename.jpg
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return ""
	}
	// Find "items" and take everything after bucket name
	for i, part := range parts {
		if part == h.storageBucket && i+1 < len(parts) {
			return strings.Join(parts[i+1:], "/")
		}
	}
	return ""
}

// DeleteMedia handles DELETE /api/admin/items/:id/media/:mediaId
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Get media ID from URL parameter
	mediaIDStr := c.Param("mediaId")
	mediaID, err := strconv.ParseInt(mediaIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid media ID format",
		})
		return
	}

	// Verify media exists and belongs to the item
	media, err := h.mediaRepo.FindByIDAndItemID(mediaID, itemUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}
	if media == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Media not found",
		})
		return
	}

	// Delete from database first
	deletedMedia, err := h.mediaRepo.Delete(mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete media from database",
		})
		return
	}

	// Delete from storage (best effort, don't fail if storage deletion fails)
	if deletedMedia != nil {
		h.cleanupStorageFiles(c, deletedMedia.URL)
		if deletedMedia.ThumbnailURL != nil {
			h.cleanupStorageFiles(c, *deletedMedia.ThumbnailURL)
		}
	}

	c.Status(http.StatusNoContent)
}

// ReorderMedia handles PUT /api/admin/items/:id/media/reorder
func (h *MediaHandler) ReorderMedia(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Parse request body
	var req domain.ReorderMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Verify all media IDs belong to the item
	mediaList, err := h.mediaRepo.FindByIDsAndItemID(req.MediaIDs, itemUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	if len(mediaList) != len(req.MediaIDs) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "One or more media IDs do not belong to this item",
		})
		return
	}

	// Create update list
	updates := make([]struct {
		ID           int64
		DisplayOrder int
	}, len(req.MediaIDs))

	for i, mediaID := range req.MediaIDs {
		updates[i] = struct {
			ID           int64
			DisplayOrder int
		}{
			ID:           mediaID,
			DisplayOrder: i,
		}
	}

	// Update display orders
	if err := h.mediaRepo.UpdateDisplayOrder(updates); err != nil {
		if errors.Is(err, errors.New("record not found")) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "One or more media not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update display order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media reordered successfully"})
}

// GetMediaList handles GET /api/items/:id/media
func (h *MediaHandler) GetMediaList(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Get media list
	mediaList, err := h.mediaRepo.FindByItemID(itemUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	// Return response
	response := &domain.MediaListResponse{
		Items: mediaList,
		Total: len(mediaList),
	}

	c.JSON(http.StatusOK, response)
}
