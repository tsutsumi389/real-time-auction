package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// BidRepository handles database operations for Bid entities
type BidRepository struct {
	db *gorm.DB
}

// NewBidRepository creates a new BidRepository instance
func NewBidRepository(db *gorm.DB) *BidRepository {
	return &BidRepository{db: db}
}

// CreateBid creates a new bid record
func (r *BidRepository) CreateBid(bid *domain.Bid) error {
	return r.db.Create(bid).Error
}

// FindBidsByItemID retrieves bids for a specific item with pagination
func (r *BidRepository) FindBidsByItemID(itemID uuid.UUID, limit, offset int) ([]domain.BidWithBidderInfo, error) {
	var bids []domain.BidWithBidderInfo

	query := r.db.Table("bids b").
		Select("b.id, b.item_id, b.bidder_id, COALESCE(bd.display_name, bd.email) as bidder_name, b.price, b.is_winning, b.bid_at").
		Joins("LEFT JOIN bidders bd ON b.bidder_id = bd.id").
		Where("b.item_id = ?", itemID).
		Order("b.bid_at DESC").
		Limit(limit).
		Offset(offset)

	result := query.Scan(&bids)
	if result.Error != nil {
		return nil, result.Error
	}

	return bids, nil
}

// CountBidsByItemID counts the total number of bids for a specific item
func (r *BidRepository) CountBidsByItemID(itemID uuid.UUID) (int64, error) {
	var count int64

	result := r.db.Model(&domain.Bid{}).
		Where("item_id = ?", itemID).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// FindWinningBidByItemID retrieves the current winning bid for an item
func (r *BidRepository) FindWinningBidByItemID(itemID uuid.UUID) (*domain.Bid, error) {
	var bid domain.Bid
	result := r.db.Where("item_id = ? AND is_winning = ?", itemID, true).First(&bid)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &bid, nil
}

// UpdateBidWinningStatus updates the is_winning flag for all bids of an item
// This sets is_winning=false for all bids except the specified bidID
func (r *BidRepository) UpdateBidWinningStatus(itemID uuid.UUID, bidID int64, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	// Set all bids to is_winning=false first
	if err := db.Model(&domain.Bid{}).
		Where("item_id = ?", itemID).
		Update("is_winning", false).Error; err != nil {
		return err
	}

	// Set the specified bid to is_winning=true
	if err := db.Model(&domain.Bid{}).
		Where("id = ?", bidID).
		Update("is_winning", true).Error; err != nil {
		return err
	}

	return nil
}
