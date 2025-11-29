package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"gorm.io/gorm"
)

var (
	// Bid-specific errors (reuse existing errors from errors.go where possible)
	ErrInsufficientPoints   = errors.New("insufficient available points")
	ErrPriceMismatch        = errors.New("price does not match current price")
	ErrBidLockFailed        = errors.New("failed to acquire bid lock, please try again")
	ErrAlreadyWinningBidder = errors.New("you are already the winning bidder")
)

const (
	BidLockTimeout = 5 * time.Second // Redis lock timeout
	BidLockRetries = 0               // No retries, fail fast
)

// BidService handles bid-related business logic
type BidService struct {
	db              *gorm.DB
	redisClient     *redis.Client
	bidRepo         *repository.BidRepository
	pointRepo       *repository.PointRepository
	auctionRepo     *repository.AuctionRepository
	ctx             context.Context
}

// NewBidService creates a new BidService instance
func NewBidService(
	db *gorm.DB,
	redisClient *redis.Client,
	bidRepo *repository.BidRepository,
	pointRepo *repository.PointRepository,
	auctionRepo *repository.AuctionRepository,
) *BidService {
	return &BidService{
		db:          db,
		redisClient: redisClient,
		bidRepo:     bidRepo,
		pointRepo:   pointRepo,
		auctionRepo: auctionRepo,
		ctx:         context.Background(),
	}
}

// PlaceBidRequest represents the request to place a bid
type PlaceBidRequest struct {
	ItemID   string
	BidderID string
	Price    int64
}

// PlaceBidResponse represents the response after placing a bid
type PlaceBidResponse struct {
	Bid    *domain.Bid           `json:"bid"`
	Points *domain.BidderPoints  `json:"points"`
}

// PlaceBid executes the bid placement with distributed locking and transaction
func (s *BidService) PlaceBid(req *PlaceBidRequest) (*PlaceBidResponse, error) {
	// Parse item ID
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	// Parse bidder ID
	bidderID, err := uuid.Parse(req.BidderID)
	if err != nil {
		return nil, fmt.Errorf("invalid bidder ID: %w", err)
	}

	// Step 1: Validate item exists and is eligible for bidding
	item, err := s.auctionRepo.FindItemByID(req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find item: %w", err)
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if item has started
	if item.StartedAt == nil {
		return nil, ErrItemNotStarted
	}

	// Check if item has ended
	if item.EndedAt != nil {
		return nil, ErrItemAlreadyEnded
	}

	// Check if current price is set
	if item.CurrentPrice == nil {
		return nil, ErrPriceMismatch
	}

	// Validate price matches current price
	if req.Price != *item.CurrentPrice {
		return nil, ErrPriceMismatch
	}

	// Step 2: Acquire distributed lock for this item
	lockKey := fmt.Sprintf("bid:lock:item:%s", req.ItemID)
	lockValue := fmt.Sprintf("%s:%d", req.BidderID, time.Now().UnixNano())

	// Try to acquire lock with SET NX EX
	acquired, err := s.redisClient.SetNX(s.ctx, lockKey, lockValue, BidLockTimeout).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !acquired {
		return nil, ErrBidLockFailed
	}

	// Ensure lock is released
	defer func() {
		// Only delete lock if it's still ours
		script := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		s.redisClient.Eval(s.ctx, script, []string{lockKey}, lockValue).Result()
	}()

	// Step 3: Check if bidder is already the winning bidder
	winningBid, err := s.bidRepo.FindWinningBidByItemID(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check winning bid: %w", err)
	}
	if winningBid != nil && winningBid.BidderID == bidderID {
		return nil, ErrAlreadyWinningBidder
	}

	// Step 4: Execute transaction
	var bid *domain.Bid
	var updatedPoints *domain.BidderPoints

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Get current points within transaction (for consistency)
		currentPoints, err := s.pointRepo.GetCurrentPoints(req.BidderID, tx)
		if err != nil {
			return fmt.Errorf("failed to get current points: %w", err)
		}
		if currentPoints == nil {
			return ErrPointsNotFound
		}

		// Check if bidder has sufficient available points
		if currentPoints.AvailablePoints < req.Price {
			return ErrInsufficientPoints
		}

		// If there is a previous winning bid, release those reserved points
		// (Note: We already checked that the bidder is not the current winning bidder in Step 3)
		if winningBid != nil {
			// Get the previous bidder's points
			previousBidderIDStr := winningBid.BidderID.String()
			previousPoints, err := s.pointRepo.GetCurrentPoints(previousBidderIDStr, tx)
			if err != nil {
				return fmt.Errorf("failed to get previous bidder points: %w", err)
			}
			if previousPoints == nil {
				return ErrPointsNotFound
			}

			// Release previous bidder's reserved points
			if err := s.pointRepo.UpdatePoints(previousBidderIDStr, winningBid.Price, -winningBid.Price, tx); err != nil {
				return fmt.Errorf("failed to release previous points: %w", err)
			}

			// Create point history for release
			releaseHistory := &domain.PointHistory{
				BidderID:       previousBidderIDStr,
				Amount:         winningBid.Price,
				Type:           domain.PointHistoryTypeRelease,
				RelatedBidID:   &winningBid.ID,
				BalanceBefore:  previousPoints.AvailablePoints,
				BalanceAfter:   previousPoints.AvailablePoints + winningBid.Price,
				ReservedBefore: previousPoints.ReservedPoints,
				ReservedAfter:  previousPoints.ReservedPoints - winningBid.Price,
				TotalBefore:    previousPoints.TotalPoints,
				TotalAfter:     previousPoints.TotalPoints,
			}
			if err := s.pointRepo.CreatePointHistory(releaseHistory, tx); err != nil {
				return fmt.Errorf("failed to create release history: %w", err)
			}
		}

		// Create bid record
		bid = &domain.Bid{
			ItemID:    itemID,
			BidderID:  bidderID,
			Price:     req.Price,
			IsWinning: true,
			BidAt:     time.Now(),
		}
		if err := s.bidRepo.CreateBid(bid); err != nil {
			return fmt.Errorf("failed to create bid: %w", err)
		}

		// Update points: available -> reserved
		if err := s.pointRepo.UpdatePoints(req.BidderID, -req.Price, req.Price, tx); err != nil {
			return fmt.Errorf("failed to update points: %w", err)
		}

		// Create point history for reserve
		reserveHistory := &domain.PointHistory{
			BidderID:       req.BidderID,
			Amount:         req.Price,
			Type:           domain.PointHistoryTypeReserve,
			RelatedBidID:   &bid.ID,
			BalanceBefore:  currentPoints.AvailablePoints,
			BalanceAfter:   currentPoints.AvailablePoints - req.Price,
			ReservedBefore: currentPoints.ReservedPoints,
			ReservedAfter:  currentPoints.ReservedPoints + req.Price,
			TotalBefore:    currentPoints.TotalPoints,
			TotalAfter:     currentPoints.TotalPoints,
		}
		if err := s.pointRepo.CreatePointHistory(reserveHistory, tx); err != nil {
			return fmt.Errorf("failed to create reserve history: %w", err)
		}

		// Update is_winning flags (set all to false except this bid)
		if err := s.bidRepo.UpdateBidWinningStatus(itemID, bid.ID, tx); err != nil {
			return fmt.Errorf("failed to update winning status: %w", err)
		}

		// Get updated points
		updatedPoints, err = s.pointRepo.GetCurrentPoints(req.BidderID, tx)
		if err != nil {
			return fmt.Errorf("failed to get updated points: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Step 5: Publish bid event to Redis Pub/Sub
	if err := s.publishBidEvent(bid, item); err != nil {
		// Log error but don't fail the bid
		fmt.Printf("Warning: failed to publish bid event: %v\n", err)
	}

	// Return response
	return &PlaceBidResponse{
		Bid:    bid,
		Points: updatedPoints,
	}, nil
}

// publishBidEvent publishes a bid event to Redis Pub/Sub
func (s *BidService) publishBidEvent(bid *domain.Bid, item *domain.Item) error {
	event := map[string]interface{}{
		"type":       "bid:placed",
		"auction_id": item.AuctionID.String(),
		"item_id":    bid.ItemID.String(),
		"bid": map[string]interface{}{
			"id":         bid.ID,
			"bidder_id":  bid.BidderID.String(),
			"price":      bid.Price,
			"is_winning": bid.IsWinning,
			"bid_at":     bid.BidAt.Format(time.RFC3339),
		},
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	channel := "auction:bid"
	if err := s.redisClient.Publish(s.ctx, channel, eventJSON).Err(); err != nil {
		return fmt.Errorf("failed to publish to Redis: %w", err)
	}

	return nil
}

// GetBidHistory retrieves the bid history for an item with bidder info
func (s *BidService) GetBidHistory(itemID string, bidderID string, limit, offset int) (*domain.BidHistoryResponse, error) {
	// Parse item ID
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	// Get bids with bidder info
	bids, err := s.bidRepo.FindBidsByItemID(itemUUID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid history: %w", err)
	}

	// Get total count
	total, err := s.bidRepo.CountBidsByItemID(itemUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to count bids: %w", err)
	}

	// Note: We're reusing BidWithBidderInfo but ideally we'd create a new type with is_own_bid field
	// For now, we'll return the standard response and let the frontend compare bidder_id

	return &domain.BidHistoryResponse{
		Total: total,
		Bids:  bids,
	}, nil
}
