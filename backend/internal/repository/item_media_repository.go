package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// ItemMediaRepository handles database operations for ItemMedia entities
type ItemMediaRepository struct {
	db *gorm.DB
}

// NewItemMediaRepository creates a new ItemMediaRepository instance
func NewItemMediaRepository(db *gorm.DB) *ItemMediaRepository {
	return &ItemMediaRepository{db: db}
}

// Create creates a new item media record
func (r *ItemMediaRepository) Create(media *domain.ItemMedia) error {
	return r.db.Create(media).Error
}

// FindByItemID retrieves all media for a specific item, ordered by display_order
func (r *ItemMediaRepository) FindByItemID(itemID uuid.UUID) ([]domain.ItemMedia, error) {
	var mediaList []domain.ItemMedia

	result := r.db.Where("item_id = ?", itemID).
		Order("display_order ASC, created_at ASC").
		Find(&mediaList)

	if result.Error != nil {
		return nil, result.Error
	}

	return mediaList, nil
}

// FindByID retrieves a single media record by ID
func (r *ItemMediaRepository) FindByID(id int64) (*domain.ItemMedia, error) {
	var media domain.ItemMedia

	result := r.db.First(&media, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &media, nil
}

// FindByIDAndItemID retrieves a media record by ID and verifies it belongs to the specified item
func (r *ItemMediaRepository) FindByIDAndItemID(id int64, itemID uuid.UUID) (*domain.ItemMedia, error) {
	var media domain.ItemMedia

	result := r.db.Where("id = ? AND item_id = ?", id, itemID).First(&media)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &media, nil
}

// Delete deletes a media record and returns the deleted record for cleanup
func (r *ItemMediaRepository) Delete(id int64) (*domain.ItemMedia, error) {
	var media domain.ItemMedia

	// First, retrieve the media record to get file paths
	result := r.db.First(&media, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Delete the record
	if err := r.db.Delete(&media).Error; err != nil {
		return nil, err
	}

	return &media, nil
}

// CountByItemIDAndType counts media records for an item filtered by type
func (r *ItemMediaRepository) CountByItemIDAndType(itemID uuid.UUID, mediaType domain.MediaType) (int64, error) {
	var count int64

	result := r.db.Model(&domain.ItemMedia{}).
		Where("item_id = ? AND media_type = ?", itemID, mediaType).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// CountByItemID counts all media records for an item
func (r *ItemMediaRepository) CountByItemID(itemID uuid.UUID) (int64, error) {
	var count int64

	result := r.db.Model(&domain.ItemMedia{}).
		Where("item_id = ?", itemID).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// UpdateDisplayOrder updates display_order for multiple media records
// This is typically used for reordering media items
func (r *ItemMediaRepository) UpdateDisplayOrder(updates []struct {
	ID           int64
	DisplayOrder int
}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			result := tx.Model(&domain.ItemMedia{}).
				Where("id = ?", update.ID).
				Update("display_order", update.DisplayOrder)

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
}

// GetNextDisplayOrder returns the next available display_order for an item
func (r *ItemMediaRepository) GetNextDisplayOrder(itemID uuid.UUID) (int, error) {
	var maxOrder int

	result := r.db.Model(&domain.ItemMedia{}).
		Where("item_id = ?", itemID).
		Select("COALESCE(MAX(display_order), -1)").
		Scan(&maxOrder)

	if result.Error != nil {
		return 0, result.Error
	}

	return maxOrder + 1, nil
}

// FindByIDsAndItemID retrieves multiple media records by IDs and verifies they all belong to the specified item
// This is useful for batch operations like reordering
func (r *ItemMediaRepository) FindByIDsAndItemID(ids []int64, itemID uuid.UUID) ([]domain.ItemMedia, error) {
	var mediaList []domain.ItemMedia

	result := r.db.Where("id IN ? AND item_id = ?", ids, itemID).Find(&mediaList)
	if result.Error != nil {
		return nil, result.Error
	}

	return mediaList, nil
}

// DeleteByItemID deletes all media for a specific item
// This is typically used when deleting an item
func (r *ItemMediaRepository) DeleteByItemID(itemID uuid.UUID) ([]domain.ItemMedia, error) {
	var mediaList []domain.ItemMedia

	// First, retrieve all media records to get file paths for cleanup
	result := r.db.Where("item_id = ?", itemID).Find(&mediaList)
	if result.Error != nil {
		return nil, result.Error
	}

	// Delete all records
	if len(mediaList) > 0 {
		if err := r.db.Where("item_id = ?", itemID).Delete(&domain.ItemMedia{}).Error; err != nil {
			return nil, err
		}
	}

	return mediaList, nil
}
