package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"gorm.io/gorm"
)

// AuctionService handles business logic for auction operations
type AuctionService struct {
	db          *gorm.DB
	auctionRepo repository.AuctionRepositoryInterface
	bidRepo     *repository.BidRepository
	pointRepo   *repository.PointRepository
	redisClient *redis.Client
	ctx         context.Context
}

// NewAuctionService creates a new AuctionService instance
func NewAuctionService(
	db *gorm.DB,
	auctionRepo repository.AuctionRepositoryInterface,
	bidRepo *repository.BidRepository,
	pointRepo *repository.PointRepository,
	redisClient *redis.Client,
) *AuctionService {
	return &AuctionService{
		db:          db,
		auctionRepo: auctionRepo,
		bidRepo:     bidRepo,
		pointRepo:   pointRepo,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// GetAuctionList retrieves a list of auctions with filters, sorting, and pagination
func (s *AuctionService) GetAuctionList(req *domain.AuctionListRequest) (*domain.AuctionListResponse, error) {
	// Validate and set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Sort == "" {
		req.Sort = "created_at_desc"
	}

	// Validate sort mode
	validSortModes := map[string]bool{
		"created_at_asc":  true,
		"created_at_desc": true,
		"updated_at_asc":  true,
		"updated_at_desc": true,
		"id_asc":          true,
		"id_desc":         true,
	}
	if !validSortModes[req.Sort] {
		return nil, ErrInvalidSortMode
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := map[domain.AuctionStatus]bool{
			domain.AuctionStatusPending:   true,
			domain.AuctionStatusActive:    true,
			domain.AuctionStatusEnded:     true,
			domain.AuctionStatusCancelled: true,
		}
		if !validStatuses[req.Status] {
			return nil, ErrInvalidStatus
		}
	}

	// Get auctions from repository
	auctions, err := s.auctionRepo.FindAuctionsWithFilters(req)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.auctionRepo.CountAuctionsWithFilters(req)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	// Build response
	response := &domain.AuctionListResponse{
		Auctions: auctions,
		Pagination: domain.Pagination{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// StartAuction starts an auction by changing its status to active
func (s *AuctionService) StartAuction(id string) (*domain.AuctionWithItemCount, error) {
	// Find auction
	auction, err := s.auctionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, ErrAuctionNotFound
	}

	// Check if auction is in pending status
	if auction.Status != domain.AuctionStatusPending {
		return nil, ErrAuctionNotPending
	}

	// Check if auction has items
	itemCount, err := s.auctionRepo.CountItemsByAuctionID(id)
	if err != nil {
		return nil, err
	}
	if itemCount == 0 {
		return nil, ErrNoItemsInAuction
	}

	// Check if all items have starting_price set
	items, err := s.auctionRepo.FindItemsByAuctionID(id)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.StartingPrice == nil {
			return nil, ErrItemsMissingStartingPrice
		}
	}

	// Update auction status to active
	err = s.auctionRepo.UpdateAuctionStatus(id, domain.AuctionStatusActive)
	if err != nil {
		return nil, err
	}

	// Return updated auction with item count
	return &domain.AuctionWithItemCount{
		ID:          auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		Status:      domain.AuctionStatusActive,
		StartedAt:   auction.StartedAt,
		ItemCount:   itemCount,
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
	}, nil
}

// EndAuction ends an auction by changing its status to ended
func (s *AuctionService) EndAuction(id string) (*domain.AuctionWithItemCount, error) {
	// Find auction
	auction, err := s.auctionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, ErrAuctionNotFound
	}

	// Check if auction is in active status
	if auction.Status != domain.AuctionStatusActive {
		return nil, ErrAuctionNotActive
	}

	// Get item count
	itemCount, err := s.auctionRepo.CountItemsByAuctionID(id)
	if err != nil {
		return nil, err
	}

	// Update auction status to ended
	err = s.auctionRepo.UpdateAuctionStatus(id, domain.AuctionStatusEnded)
	if err != nil {
		return nil, err
	}

	// TODO: Set ended_at for all items and finalize winners
	// This will be implemented when we add item-level operations

	// Return updated auction with item count
	return &domain.AuctionWithItemCount{
		ID:          auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		Status:      domain.AuctionStatusEnded,
		StartedAt:   auction.StartedAt,
		ItemCount:   itemCount,
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
	}, nil
}

// CancelAuction cancels an auction by changing its status to cancelled
func (s *AuctionService) CancelAuction(id string) (*domain.AuctionWithItemCount, error) {
	// Find auction
	auction, err := s.auctionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, ErrAuctionNotFound
	}

	// Check if auction is in active status
	if auction.Status != domain.AuctionStatusActive {
		return nil, ErrAuctionNotActive
	}

	// Get item count
	itemCount, err := s.auctionRepo.CountItemsByAuctionID(id)
	if err != nil {
		return nil, err
	}

	// Update auction status to cancelled
	err = s.auctionRepo.UpdateAuctionStatus(id, domain.AuctionStatusCancelled)
	if err != nil {
		return nil, err
	}

	// TODO: Invalidate bids and refund reserved points
	// This will be implemented when we add bid and point operations

	// Return updated auction with item count
	return &domain.AuctionWithItemCount{
		ID:          auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		Status:      domain.AuctionStatusCancelled,
		StartedAt:   auction.StartedAt,
		ItemCount:   itemCount,
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
	}, nil
}

// CreateAuction creates a new auction with items
func (s *AuctionService) CreateAuction(req *domain.CreateAuctionRequest) (*domain.CreateAuctionResponse, error) {
	// Create auction entity
	auction := &domain.Auction{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.AuctionStatusPending,
		StartedAt:   req.StartedAt,
	}

	// Create item entities
	items := make([]domain.Item, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = domain.Item{
			Name:          itemReq.Name,
			Description:   itemReq.Description,
			LotNumber:     itemReq.LotNumber,
			StartingPrice: itemReq.StartingPrice,
		}
	}

	// Create auction with items in transaction
	err := s.auctionRepo.CreateAuctionWithItems(auction, items)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &domain.CreateAuctionResponse{
		ID:          auction.ID,
		Title:       auction.Title,
		Description: auction.Description,
		Status:      auction.Status,
		StartedAt:   auction.StartedAt,
		ItemCount:   len(items),
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
	}

	return response, nil
}

// GetBidderAuctionList retrieves public auctions for bidders with filters, sorting, and offset/limit pagination
func (s *AuctionService) GetBidderAuctionList(req *domain.BidderAuctionListRequest) (*domain.BidderAuctionListResponse, error) {
	// Validate and set defaults
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Sort == "" {
		req.Sort = "started_at_desc"
	}

	// Validate sort mode
	validSortModes := map[string]bool{
		"started_at_asc":  true,
		"started_at_desc": true,
		"updated_at_asc":  true,
		"updated_at_desc": true,
	}
	if !validSortModes[req.Sort] {
		return nil, ErrInvalidSortMode
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := map[domain.AuctionStatus]bool{
			domain.AuctionStatusActive:    true,
			domain.AuctionStatusEnded:     true,
			domain.AuctionStatusCancelled: true,
		}
		if !validStatuses[req.Status] {
			return nil, ErrInvalidStatus
		}
	}

	// Get public auctions from repository
	auctions, err := s.auctionRepo.FindPublicAuctionsWithFilters(req)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.auctionRepo.CountPublicAuctionsWithFilters(req)
	if err != nil {
		return nil, err
	}

	// Calculate has_more flag
	hasMore := total > int64(req.Offset+req.Limit)

	// Build response
	response := &domain.BidderAuctionListResponse{
		Auctions: auctions,
		Pagination: domain.BidderPagination{
			Total:   total,
			Offset:  req.Offset,
			Limit:   req.Limit,
			HasMore: hasMore,
		},
	}

	return response, nil
}

// GetAuctionDetail retrieves auction details with all items
func (s *AuctionService) GetAuctionDetail(id string) (*domain.GetAuctionDetailResponse, error) {
	auction, err := s.auctionRepo.FindAuctionWithItems(id)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, ErrAuctionNotFound
	}
	return auction, nil
}

// StartItem starts an item auction
func (s *AuctionService) StartItem(itemID string) (*domain.StartItemResponse, error) {
	// Start the item
	item, err := s.auctionRepo.StartItem(itemID)
	if err != nil {
		if err.Error() == "item already started" {
			return nil, ErrItemAlreadyStarted
		}
		if err.Error() == "starting price not set" {
			return nil, ErrStartingPriceNotSet
		}
		return nil, err
	}

	if item == nil {
		return nil, ErrItemNotFound
	}

	// Publish WebSocket event to Redis Pub/Sub
	if s.redisClient != nil {
		event := map[string]interface{}{
			"type":    "item:started",
			"item_id": item.ID.String(),
			"payload": map[string]interface{}{
				"item": map[string]interface{}{
					"id":            item.ID.String(),
					"auction_id":    item.AuctionID.String(),
					"name":          item.Name,
					"current_price": item.CurrentPrice,
					"started_at":    item.StartedAt,
					"status":        "active",
				},
			},
		}
		eventJSON, err := json.Marshal(event)
		if err == nil {
			channel := "auction:item_started"
			_ = s.redisClient.Publish(s.ctx, channel, eventJSON).Err()
		}
	}

	// Build response
	return &domain.StartItemResponse{
		ItemID:       item.ID,
		AuctionID:    item.AuctionID,
		CurrentPrice: *item.CurrentPrice,
		StartedAt:    *item.StartedAt,
	}, nil
}

// OpenPrice opens a new price for an item
func (s *AuctionService) OpenPrice(itemID string, newPrice int64, adminID int64) (*domain.OpenPriceResponse, error) {
	// Find the item
	item, err := s.auctionRepo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if item has been started
	if item.StartedAt == nil {
		return nil, ErrItemNotStarted
	}

	// Check if item has already ended
	if item.EndedAt != nil {
		return nil, ErrItemAlreadyEnded
	}

	// Check if new price is higher than current price
	if item.CurrentPrice != nil && newPrice <= *item.CurrentPrice {
		return nil, ErrPriceTooLow
	}

	// Get previous price
	previousPrice := int64(0)
	if item.CurrentPrice != nil {
		previousPrice = *item.CurrentPrice
	}

	// Execute transaction to update price and release reserved points
	var priceHistory *domain.PriceHistory
	var hadBid bool

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Check if there was a bid at the previous price
		var winningBid *domain.Bid
		if previousPrice > 0 {
			winningBid, err = s.bidRepo.FindWinningBidByItemID(item.ID)
			if err != nil {
				return fmt.Errorf("failed to find winning bid: %w", err)
			}
			hadBid = winningBid != nil && winningBid.Price == previousPrice
		}

		// If there is a winning bid, release the reserved points (price has changed, old bid is invalid)
		if winningBid != nil {
			bidderIDStr := winningBid.BidderID.String()

			// Get bidder's current points
			currentPoints, err := s.pointRepo.GetCurrentPoints(bidderIDStr, tx)
			if err != nil {
				return fmt.Errorf("failed to get current points: %w", err)
			}
			if currentPoints == nil {
				return ErrPointsNotFound
			}

			// Release reserved points
			if err := s.pointRepo.UpdatePoints(bidderIDStr, winningBid.Price, -winningBid.Price, tx); err != nil {
				return fmt.Errorf("failed to release points: %w", err)
			}

			// Create point history for release
			releaseHistory := &domain.PointHistory{
				BidderID:       bidderIDStr,
				Amount:         winningBid.Price,
				Type:           domain.PointHistoryTypeRelease,
				RelatedBidID:   &winningBid.ID,
				BalanceBefore:  currentPoints.AvailablePoints,
				BalanceAfter:   currentPoints.AvailablePoints + winningBid.Price,
				ReservedBefore: currentPoints.ReservedPoints,
				ReservedAfter:  currentPoints.ReservedPoints - winningBid.Price,
				TotalBefore:    currentPoints.TotalPoints,
				TotalAfter:     currentPoints.TotalPoints,
			}
			if err := s.pointRepo.CreatePointHistory(releaseHistory, tx); err != nil {
				return fmt.Errorf("failed to create release history: %w", err)
			}

			// Mark all bids as not winning
			if err := s.bidRepo.UpdateBidWinningStatus(item.ID, 0, tx); err != nil {
				return fmt.Errorf("failed to update winning status: %w", err)
			}
		}

		// Update item current price (using tx for transaction consistency)
		if err := tx.Model(&domain.Item{}).Where("id = ?", itemID).Update("current_price", newPrice).Error; err != nil {
			return fmt.Errorf("failed to update item price: %w", err)
		}

		// Create price history record
		now := time.Now()
		priceHistory = &domain.PriceHistory{
			ItemID:      item.ID,
			Price:       newPrice,
			DisclosedBy: adminID,
			HadBid:      hadBid,
			DisclosedAt: now,
		}
		if err := tx.Create(priceHistory).Error; err != nil {
			return fmt.Errorf("failed to create price history: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish WebSocket event to Redis Pub/Sub
	if s.redisClient != nil {
		event := map[string]interface{}{
			"type":    "price:opened",
			"item_id": item.ID.String(),
			"payload": map[string]interface{}{
				"item_id": item.ID.String(),
				"price":   newPrice,
				"price_history": map[string]interface{}{
					"id":           priceHistory.ID,
					"item_id":      priceHistory.ItemID.String(),
					"price":        priceHistory.Price,
					"disclosed_by": priceHistory.DisclosedBy,
					"had_bid":      priceHistory.HadBid,
					"disclosed_at": priceHistory.DisclosedAt,
				},
			},
		}
		eventJSON, err := json.Marshal(event)
		if err == nil {
			channel := "auction:price_open"
			_ = s.redisClient.Publish(s.ctx, channel, eventJSON).Err()
		}
	}

	// Build response
	return &domain.OpenPriceResponse{
		ItemID:        item.ID,
		CurrentPrice:  newPrice,
		PreviousPrice: previousPrice,
		DisclosedAt:   now,
		PriceHistory:  priceHistory,
	}, nil
}

// EndItem ends an item auction
func (s *AuctionService) EndItem(itemID string) (*domain.EndItemResponse, error) {
	// Find the item
	item, err := s.auctionRepo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if item has been started
	if item.StartedAt == nil {
		return nil, ErrItemNotStarted
	}

	// Check if item has already ended
	if item.EndedAt != nil {
		return nil, ErrItemAlreadyEnded
	}

	// Find the winning bid
	winningBid, err := s.auctionRepo.FindWinningBidByItemID(itemID)
	if err != nil {
		return nil, err
	}
	if winningBid == nil {
		return nil, ErrNoBidsFound
	}

	// End the item with winner information
	finalPrice := winningBid.Price
	endedItem, err := s.auctionRepo.EndItem(itemID, winningBid.BidderID, finalPrice)
	if err != nil {
		return nil, err
	}

	// Get winner name
	var winnerName *string
	bids, err := s.auctionRepo.FindBidsByItemID(itemID, 1, 0)
	if err == nil && len(bids) > 0 && bids[0].IsWinning {
		winnerName = &bids[0].BidderName
	}

	// TODO: Handle point consumption and release for winners and non-winners
	// This will be implemented when we add point transaction logic

	// Publish WebSocket event to Redis Pub/Sub
	if s.redisClient != nil {
		event := map[string]interface{}{
			"type":    "item:ended",
			"item_id": endedItem.ID.String(),
			"payload": map[string]interface{}{
				"item": map[string]interface{}{
					"id":          endedItem.ID.String(),
					"auction_id":  endedItem.AuctionID.String(),
					"name":        endedItem.Name,
					"final_price": finalPrice,
					"winner_id":   endedItem.WinnerID,
					"ended_at":    endedItem.EndedAt,
					"status":      "ended",
				},
			},
		}
		eventJSON, err := json.Marshal(event)
		if err == nil {
			channel := "auction:item_ended"
			_ = s.redisClient.Publish(s.ctx, channel, eventJSON).Err()
		}
	}

	// Build response
	return &domain.EndItemResponse{
		ItemID:     endedItem.ID,
		WinnerID:   endedItem.WinnerID,
		WinnerName: winnerName,
		FinalPrice: finalPrice,
		EndedAt:    *endedItem.EndedAt,
	}, nil
}

// GetBidHistory retrieves bid history for an item
func (s *AuctionService) GetBidHistory(itemID string, limit int, offset int) (*domain.BidHistoryResponse, error) {
	// Set defaults
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	// Get bids from repository
	bids, err := s.auctionRepo.FindBidsByItemID(itemID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.auctionRepo.CountBidsByItemID(itemID)
	if err != nil {
		return nil, err
	}

	// Build response
	return &domain.BidHistoryResponse{
		Total: total,
		Bids:  bids,
	}, nil
}

// GetPriceHistory retrieves price disclosure history for an item
func (s *AuctionService) GetPriceHistory(itemID string) (*domain.PriceHistoryResponse, error) {
	// Get price history from repository
	history, err := s.auctionRepo.FindPriceHistoryByItemID(itemID)
	if err != nil {
		return nil, err
	}

	// Build response
	return &domain.PriceHistoryResponse{
		Total:   int64(len(history)),
		History: history,
	}, nil
}

// GetParticipants retrieves participants for an auction
func (s *AuctionService) GetParticipants(auctionID string) (*domain.ParticipantsResponse, error) {
	// Get participants from repository
	participants, err := s.auctionRepo.FindParticipantsByAuctionID(auctionID)
	if err != nil {
		return nil, err
	}

	// Build response
	return &domain.ParticipantsResponse{
		Total:        int64(len(participants)),
		Participants: participants,
	}, nil
}

// CancelAuctionWithReason cancels an auction with a reason
func (s *AuctionService) CancelAuctionWithReason(auctionID string, reason string) (*domain.CancelAuctionResponse, error) {
	// Find auction
	auction, err := s.auctionRepo.FindByID(auctionID)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, ErrAuctionNotFound
	}

	// Check if auction is active
	if auction.Status != domain.AuctionStatusActive {
		return nil, ErrAuctionNotActive
	}

	// Cancel auction with refunds
	response, err := s.auctionRepo.CancelAuctionWithRefunds(auctionID, reason)
	if err != nil {
		return nil, err
	}

	return response, nil
}
