package service

import (
	"math"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

// AuctionService handles business logic for auction operations
type AuctionService struct {
	auctionRepo repository.AuctionRepositoryInterface
}

// NewAuctionService creates a new AuctionService instance
func NewAuctionService(auctionRepo repository.AuctionRepositoryInterface) *AuctionService {
	return &AuctionService{
		auctionRepo: auctionRepo,
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
		ItemCount:   itemCount,
		CreatedAt:   auction.CreatedAt,
		UpdatedAt:   auction.UpdatedAt,
	}, nil
}
