package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

// ItemMedia represents media (images/videos) associated with an auction item
type ItemMedia struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID       uuid.UUID `gorm:"type:uuid;not null;index:idx_item_media_item" json:"item_id"`
	MediaType    MediaType `gorm:"type:varchar(20);not null;index:idx_item_media_type" json:"media_type"`
	URL          string    `gorm:"type:varchar(500);not null" json:"url"`
	ThumbnailURL *string   `gorm:"type:varchar(500)" json:"thumbnail_url,omitempty"`
	DisplayOrder int       `gorm:"not null;default:0;index:idx_item_media_item" json:"display_order"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Item *Item `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for ItemMedia model
func (ItemMedia) TableName() string {
	return "item_media"
}

// Validate performs validation on ItemMedia fields
func (m *ItemMedia) Validate() error {
	// Validate media type
	if m.MediaType != MediaTypeImage && m.MediaType != MediaTypeVideo {
		return ErrInvalidMediaType
	}

	// Validate URL
	if m.URL == "" {
		return ErrInvalidMediaURL
	}

	// Validate display order
	if m.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}

	return nil
}

// IsImage checks if the media is an image
func (m *ItemMedia) IsImage() bool {
	return m.MediaType == MediaTypeImage
}

// IsVideo checks if the media is a video
func (m *ItemMedia) IsVideo() bool {
	return m.MediaType == MediaTypeVideo
}

// UploadMediaRequest represents the request to upload media
type UploadMediaRequest struct {
	ItemID    uuid.UUID `json:"item_id"`
	MediaType MediaType `json:"media_type" binding:"required,oneof=image video"`
}

// UploadMediaResponse represents the response for media upload
type UploadMediaResponse struct {
	ID           int64     `json:"id"`
	ItemID       uuid.UUID `json:"item_id"`
	MediaType    MediaType `json:"media_type"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// ReorderMediaRequest represents the request to reorder media items
type ReorderMediaRequest struct {
	MediaIDs []int64 `json:"media_ids" binding:"required,min=1"`
}

// MediaListResponse represents the response for media list
type MediaListResponse struct {
	Items []ItemMedia `json:"items"`
	Total int         `json:"total"`
}

// Error definitions for media operations
var (
	ErrInvalidMediaType    = errors.New("invalid media type: must be 'image' or 'video'")
	ErrInvalidMediaURL     = errors.New("invalid media URL: cannot be empty")
	ErrInvalidDisplayOrder = errors.New("invalid display order: must be non-negative")
	ErrMediaLimitExceeded  = errors.New("media limit exceeded")
	ErrMediaNotFound       = errors.New("media not found")
	ErrFileTooLarge        = errors.New("file size exceeds limit")
	ErrInvalidFileType     = errors.New("invalid file type")
)
