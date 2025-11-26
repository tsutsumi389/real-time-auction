package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// FindAuctionWithItems finds an auction by ID with all its items
func (r *AuctionRepository) FindAuctionWithItems(id string) (*domain.GetAuctionDetailResponse, error) {
	auctionID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var auction domain.Auction
	if err := r.db.First(&auction, "id = ?", auctionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var items []domain.Item
	if err := r.db.Where("auction_id = ?", auctionID).
		Order("lot_number ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &domain.GetAuctionDetailResponse{
		ID:          auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		Status:      auction.Status,
		StartedAt:   auction.StartedAt,
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
		Items:       items,
	}, nil
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

// FindItemByID finds an item by ID
func (r *AuctionRepository) FindItemByID(itemID string) (*domain.Item, error) {
	var item domain.Item
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	result := r.db.First(&item, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &item, nil
}

// StartItem starts an item by setting its current_price to starting_price and recording started_at
func (r *AuctionRepository) StartItem(itemID string) (*domain.Item, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	var item domain.Item
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Find the item with FOR UPDATE lock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&item, "id = ?", id).Error; err != nil {
			return err
		}

		// Check if already started
		if item.StartedAt != nil {
			return errors.New("item already started")
		}

		// Check if starting_price is set
		if item.StartingPrice == nil {
			return errors.New("starting price not set")
		}

		// Update item
		now := time.Now()
		currentPrice := *item.StartingPrice
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"current_price": currentPrice,
			"started_at":    now,
		}).Error; err != nil {
			return err
		}

		// Update the in-memory object
		item.CurrentPrice = &currentPrice
		item.StartedAt = &now

		// Update auction status to active if this is the first item
		var auction domain.Auction
		if err := tx.First(&auction, "id = ?", item.AuctionID).Error; err != nil {
			return err
		}
		if auction.Status == domain.AuctionStatusPending {
			if err := tx.Model(&auction).Update("status", domain.AuctionStatusActive).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItemCurrentPrice updates the current price of an item
func (r *AuctionRepository) UpdateItemCurrentPrice(itemID string, price int64) error {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return err
	}

	return r.db.Model(&domain.Item{}).
		Where("id = ?", id).
		Update("current_price", price).Error
}

// EndItem ends an item by setting winner_id and ended_at
func (r *AuctionRepository) EndItem(itemID string, winnerID uuid.UUID, finalPrice int64) (*domain.Item, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	var item domain.Item
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Find the item
		if err := tx.First(&item, "id = ?", id).Error; err != nil {
			return err
		}

		// Update item
		now := time.Now()
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"winner_id":     winnerID,
			"current_price": finalPrice,
			"ended_at":      now,
		}).Error; err != nil {
			return err
		}

		// Update in-memory object
		item.WinnerID = &winnerID
		item.CurrentPrice = &finalPrice
		item.EndedAt = &now

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// FindBidsByItemID retrieves bids for an item with pagination
func (r *AuctionRepository) FindBidsByItemID(itemID string, limit int, offset int) ([]domain.BidWithBidderInfo, error) {
	var results []domain.BidWithBidderInfo
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	query := r.db.Table("bids b").
		Select("b.id, b.item_id, b.bidder_id, bd.display_name as bidder_name, b.price, b.is_winning, b.bid_at").
		Joins("LEFT JOIN bidders bd ON b.bidder_id = bd.id").
		Where("b.item_id = ?", id).
		Order("b.bid_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// CountBidsByItemID counts bids for an item
func (r *AuctionRepository) CountBidsByItemID(itemID string) (int64, error) {
	var count int64
	id, err := uuid.Parse(itemID)
	if err != nil {
		return 0, err
	}

	if err := r.db.Model(&domain.Bid{}).
		Where("item_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindWinningBidByItemID finds the winning bid for an item
func (r *AuctionRepository) FindWinningBidByItemID(itemID string) (*domain.Bid, error) {
	var bid domain.Bid
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	result := r.db.Where("item_id = ? AND is_winning = ?", id, true).First(&bid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &bid, nil
}

// CreateBid creates a new bid
func (r *AuctionRepository) CreateBid(bid *domain.Bid) error {
	return r.db.Create(bid).Error
}

// UpdateBidWinningStatus updates the is_winning flag for bids
func (r *AuctionRepository) UpdateBidWinningStatus(itemID string, winningBidID int64) error {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Set all bids for this item to is_winning = false
		if err := tx.Model(&domain.Bid{}).
			Where("item_id = ?", id).
			Update("is_winning", false).Error; err != nil {
			return err
		}

		// Set the winning bid to is_winning = true
		if err := tx.Model(&domain.Bid{}).
			Where("id = ?", winningBidID).
			Update("is_winning", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindPriceHistoryByItemID retrieves price history for an item
func (r *AuctionRepository) FindPriceHistoryByItemID(itemID string) ([]domain.PriceHistoryWithAdmin, error) {
	var results []domain.PriceHistoryWithAdmin
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	query := r.db.Table("price_history ph").
		Select("ph.id, ph.item_id, ph.price, a.display_name as disclosed_by_name, ph.had_bid, ph.disclosed_at").
		Joins("LEFT JOIN admins a ON ph.disclosed_by = a.id").
		Where("ph.item_id = ?", id).
		Order("ph.disclosed_at DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// CreatePriceHistory creates a new price history record
func (r *AuctionRepository) CreatePriceHistory(history *domain.PriceHistory) error {
	return r.db.Create(history).Error
}

// FindParticipantsByAuctionID retrieves participants for an auction
func (r *AuctionRepository) FindParticipantsByAuctionID(auctionID string) ([]domain.ParticipantInfo, error) {
	var results []domain.ParticipantInfo
	id, err := uuid.Parse(auctionID)
	if err != nil {
		return nil, err
	}

	query := r.db.Table("bidders bd").
		Select(`bd.id as bidder_id,
			bd.display_name,
			COUNT(b.id) as bid_count,
			false as is_online,
			MAX(b.bid_at) as last_bid_at`).
		Joins("LEFT JOIN bids b ON bd.id = b.bidder_id").
		Joins("LEFT JOIN items i ON b.item_id = i.id").
		Where("i.auction_id = ?", id).
		Group("bd.id, bd.display_name").
		Order("bid_count DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// CancelAuctionWithRefunds cancels an auction and refunds all reserved points
func (r *AuctionRepository) CancelAuctionWithRefunds(auctionID string, reason string) (*domain.CancelAuctionResponse, error) {
	id, err := uuid.Parse(auctionID)
	if err != nil {
		return nil, err
	}

	var response domain.CancelAuctionResponse
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Update auction status to cancelled
		if err := tx.Model(&domain.Auction{}).
			Where("id = ?", id).
			Update("status", domain.AuctionStatusCancelled).Error; err != nil {
			return err
		}

		// Set ended_at for all items in this auction
		now := time.Now()
		if err := tx.Model(&domain.Item{}).
			Where("auction_id = ? AND ended_at IS NULL", id).
			Update("ended_at", now).Error; err != nil {
			return err
		}

		// Get all bidders who have reserved points for this auction
		var bidderPoints []struct {
			BidderID      uuid.UUID `gorm:"column:bidder_id"`
			ReservedPoints int64     `gorm:"column:reserved_points"`
		}

		if err := tx.Table("bidder_points bp").
			Select("bp.bidder_id, bp.reserved_points").
			Joins("JOIN bids b ON bp.bidder_id = b.bidder_id").
			Joins("JOIN items i ON b.item_id = i.id").
			Where("i.auction_id = ? AND bp.reserved_points > 0", id).
			Group("bp.bidder_id, bp.reserved_points").
			Scan(&bidderPoints).Error; err != nil {
			return err
		}

		var totalRefunded int64
		refundedCount := int64(len(bidderPoints))

		// Refund reserved points for each bidder
		for _, bp := range bidderPoints {
			// Get current point balance
			var currentPoints domain.BidderPoints
			if err := tx.Where("bidder_id = ?", bp.BidderID).First(&currentPoints).Error; err != nil {
				return err
			}

			// Calculate new balances
			newAvailable := currentPoints.AvailablePoints + bp.ReservedPoints
			newReserved := int64(0)

			// Update bidder_points
			if err := tx.Model(&domain.BidderPoints{}).
				Where("bidder_id = ?", bp.BidderID).
				Updates(map[string]interface{}{
					"available_points": newAvailable,
					"reserved_points":  newReserved,
				}).Error; err != nil {
				return err
			}

			// Create point_history record
			reasonStr := fmt.Sprintf("Auction cancelled: %s", reason)
			pointHistory := &domain.PointHistory{
				BidderID:         bp.BidderID.String(),
				Amount:           bp.ReservedPoints,
				Type:             "refund",
				Reason:           &reasonStr,
				BalanceBefore:    currentPoints.AvailablePoints,
				BalanceAfter:     newAvailable,
				ReservedBefore:   currentPoints.ReservedPoints,
				ReservedAfter:    newReserved,
				TotalBefore:      currentPoints.TotalPoints,
				TotalAfter:       currentPoints.TotalPoints,
				// Note: related_auction_id is BIGINT in DB but auction.id is UUID
				// We can't set this field until schema is fixed
				RelatedAuctionID: nil,
			}
			if err := tx.Create(pointHistory).Error; err != nil {
				return err
			}

			totalRefunded += bp.ReservedPoints
		}

		// Build response
		response = domain.CancelAuctionResponse{
			AuctionID:           id,
			Status:              "cancelled",
			RefundedBidders:     refundedCount,
			TotalRefundedPoints: totalRefunded,
			CancelledAt:         now,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}
