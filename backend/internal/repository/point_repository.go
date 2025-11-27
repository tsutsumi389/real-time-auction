package repository

import (
	"errors"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// PointRepository handles database operations for BidderPoints entities
type PointRepository struct {
	db *gorm.DB
}

// NewPointRepository creates a new PointRepository instance
func NewPointRepository(db *gorm.DB) *PointRepository {
	return &PointRepository{db: db}
}

// FindPointsByBidderID retrieves the points information for a bidder
func (r *PointRepository) FindPointsByBidderID(bidderID string) (*domain.BidderPoints, error) {
	var points domain.BidderPoints
	result := r.db.Where("bidder_id = ?", bidderID).First(&points)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &points, nil
}

// UpdatePoints updates the available and reserved points for a bidder
// availableDelta: change in available_points (can be negative)
// reservedDelta: change in reserved_points (can be negative)
func (r *PointRepository) UpdatePoints(bidderID string, availableDelta, reservedDelta int64, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	// Use raw SQL to update with incremental changes
	result := db.Exec(`
		UPDATE bidder_points
		SET available_points = available_points + ?,
		    reserved_points = reserved_points + ?,
		    updated_at = NOW()
		WHERE bidder_id = ?
	`, availableDelta, reservedDelta, bidderID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CreatePointHistory creates a new point history record
func (r *PointRepository) CreatePointHistory(history *domain.PointHistory, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.Create(history).Error
}

// GetCurrentPoints retrieves current points within a transaction (for consistency)
func (r *PointRepository) GetCurrentPoints(bidderID string, tx *gorm.DB) (*domain.BidderPoints, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	var points domain.BidderPoints
	result := db.Where("bidder_id = ?", bidderID).First(&points)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &points, nil
}
