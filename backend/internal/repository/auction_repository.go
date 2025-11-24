package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// AuctionRepository handles database operations for Auction entities
type AuctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository creates a new AuctionRepository instance
func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

// FindByID finds an auction by ID
func (r *AuctionRepository) FindByID(id string) (*domain.Auction, error) {
	var auction domain.Auction
	auctionID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result := r.db.First(&auction, "id = ?", auctionID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &auction, nil
}

// FindAuctionsWithFilters retrieves auctions with item counts, filters, sorting, and pagination
func (r *AuctionRepository) FindAuctionsWithFilters(req *domain.AuctionListRequest) ([]domain.AuctionWithItemCount, error) {
	var results []domain.AuctionWithItemCount

	query := r.db.Model(&domain.Auction{}).
		Select(`auctions.id, auctions.title, auctions.description, auctions.status,
			auctions.created_at, auctions.updated_at, COUNT(items.id) as item_count`).
		Joins("LEFT JOIN items ON items.auction_id = auctions.id")

	// Apply keyword filter (title and description search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("(auctions.title ILIKE ? OR auctions.description ILIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Apply status filter
	if req.Status != "" {
		query = query.Where("auctions.status = ?", req.Status)
	}

	// Apply date range filters
	if req.CreatedAfter != "" {
		createdAfter, err := time.Parse("2006-01-02", req.CreatedAfter)
		if err == nil {
			query = query.Where("auctions.created_at >= ?", createdAfter)
		}
	}

	if req.UpdatedBefore != "" {
		updatedBefore, err := time.Parse("2006-01-02", req.UpdatedBefore)
		if err == nil {
			// Add 1 day to include the entire day
			updatedBefore = updatedBefore.Add(24 * time.Hour)
			query = query.Where("auctions.updated_at < ?", updatedBefore)
		}
	}

	// Group by auction fields (required for COUNT aggregate)
	query = query.Group("auctions.id, auctions.title, auctions.description, auctions.status, auctions.created_at, auctions.updated_at")

	// Apply sorting
	switch req.Sort {
	case "created_at_asc":
		query = query.Order("auctions.created_at ASC")
	case "created_at_desc":
		query = query.Order("auctions.created_at DESC")
	case "updated_at_asc":
		query = query.Order("auctions.updated_at ASC")
	case "updated_at_desc":
		query = query.Order("auctions.updated_at DESC")
	case "id_asc":
		query = query.Order("auctions.id ASC")
	case "id_desc":
		query = query.Order("auctions.id DESC")
	default:
		// Default: created_at descending (newest first)
		query = query.Order("auctions.created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	// Execute query
	result := query.Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

// CountAuctionsWithFilters counts auctions matching the filters
func (r *AuctionRepository) CountAuctionsWithFilters(req *domain.AuctionListRequest) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Auction{})

	// Apply keyword filter
	if req.Keyword != "" {
		query = query.Where("(title ILIKE ? OR description ILIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Apply status filter
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Apply date range filters
	if req.CreatedAfter != "" {
		createdAfter, err := time.Parse("2006-01-02", req.CreatedAfter)
		if err == nil {
			query = query.Where("created_at >= ?", createdAfter)
		}
	}

	if req.UpdatedBefore != "" {
		updatedBefore, err := time.Parse("2006-01-02", req.UpdatedBefore)
		if err == nil {
			// Add 1 day to include the entire day
			updatedBefore = updatedBefore.Add(24 * time.Hour)
			query = query.Where("updated_at < ?", updatedBefore)
		}
	}

	// Execute count query
	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// UpdateAuctionStatus updates the status of an auction
func (r *AuctionRepository) UpdateAuctionStatus(id string, status domain.AuctionStatus) error {
	auctionID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.Model(&domain.Auction{}).
		Where("id = ?", auctionID).
		Update("status", status).Error
}

// CountItemsByAuctionID counts items in an auction
func (r *AuctionRepository) CountItemsByAuctionID(auctionID string) (int64, error) {
	var count int64
	id, err := uuid.Parse(auctionID)
	if err != nil {
		return 0, err
	}

	result := r.db.Model(&domain.Item{}).
		Where("auction_id = ?", id).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// FindItemsByAuctionID retrieves all items for an auction
func (r *AuctionRepository) FindItemsByAuctionID(auctionID string) ([]domain.Item, error) {
	var items []domain.Item
	id, err := uuid.Parse(auctionID)
	if err != nil {
		return nil, err
	}

	result := r.db.Where("auction_id = ?", id).
		Order("lot_number ASC").
		Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// CreateAuction creates a new auction
func (r *AuctionRepository) CreateAuction(auction *domain.Auction) error {
	return r.db.Create(auction).Error
}

// CreateItems creates multiple items for an auction
func (r *AuctionRepository) CreateItems(items []domain.Item) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

// CreateAuctionWithItems creates an auction with items in a transaction
func (r *AuctionRepository) CreateAuctionWithItems(auction *domain.Auction, items []domain.Item) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create auction
		if err := tx.Create(auction).Error; err != nil {
			return err
		}

		// Create items if provided
		if len(items) > 0 {
			// Set auction_id for all items
			for i := range items {
				items[i].AuctionID = auction.ID
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindPublicAuctionsWithFilters retrieves public auctions (non-pending) with filters, sorting, and offset/limit pagination
func (r *AuctionRepository) FindPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) ([]domain.BidderAuctionSummary, error) {
	var results []domain.BidderAuctionSummary

	query := r.db.Model(&domain.Auction{}).
		Select(`auctions.id, auctions.title, auctions.description, auctions.status,
			auctions.started_at, auctions.created_at, auctions.updated_at,
			COUNT(items.id) as item_count,
			(
				SELECT im.thumbnail_url
				FROM items i2
				LEFT JOIN item_media im ON im.item_id = i2.id AND im.media_type = 'image'
				WHERE i2.auction_id = auctions.id
				ORDER BY i2.lot_number ASC, im.display_order ASC
				LIMIT 1
			) AS thumbnail_url`).
		Joins("LEFT JOIN items ON items.auction_id = auctions.id").
		Where("auctions.status IN ?", []domain.AuctionStatus{
			domain.AuctionStatusActive,
			domain.AuctionStatusEnded,
			domain.AuctionStatusCancelled,
		})

	// Apply keyword filter (title search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("auctions.title ILIKE ?", "%"+req.Keyword+"%")
	}

	// Apply status filter
	if req.Status != "" {
		query = query.Where("auctions.status = ?", req.Status)
	}

	// Group by auction fields (required for COUNT aggregate)
	query = query.Group("auctions.id, auctions.title, auctions.description, auctions.status, auctions.started_at, auctions.created_at, auctions.updated_at")

	// Apply sorting
	switch req.Sort {
	case "started_at_asc":
		query = query.Order("auctions.started_at ASC")
	case "started_at_desc":
		query = query.Order("auctions.started_at DESC")
	case "updated_at_asc":
		query = query.Order("auctions.updated_at ASC")
	case "updated_at_desc":
		query = query.Order("auctions.updated_at DESC")
	default:
		// Default: started_at descending (newest first)
		query = query.Order("auctions.started_at DESC")
	}

	// Apply offset/limit pagination
	query = query.Offset(req.Offset).Limit(req.Limit)

	// Execute query
	result := query.Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

// CountPublicAuctionsWithFilters counts public auctions (non-pending) matching the filters
func (r *AuctionRepository) CountPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Auction{}).
		Where("status IN ?", []domain.AuctionStatus{
			domain.AuctionStatusActive,
			domain.AuctionStatusEnded,
			domain.AuctionStatusCancelled,
		})

	// Apply keyword filter
	if req.Keyword != "" {
		query = query.Where("title ILIKE ?", "%"+req.Keyword+"%")
	}

	// Apply status filter
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Execute count query
	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
