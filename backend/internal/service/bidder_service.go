package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

var (
	ErrBidderNotFound       = errors.New("bidder not found")
	ErrInvalidBidderStatus  = errors.New("invalid bidder status value")
	ErrInvalidPoints        = errors.New("points must be greater than 0")
	ErrPointsExceedMaximum  = errors.New("points exceed maximum limit")
	ErrInvalidBidderSortMode = errors.New("invalid sort mode for bidders")
)

const (
	MaxPointsPerGrant = 1000000 // Maximum points that can be granted at once
)

// BidderService handles bidder-related business logic
type BidderService struct {
	bidderRepo repository.BidderRepositoryInterface
}

// NewBidderService creates a new BidderService instance
func NewBidderService(bidderRepo repository.BidderRepositoryInterface) *BidderService {
	return &BidderService{
		bidderRepo: bidderRepo,
	}
}

// GetBidderByID retrieves a single bidder by ID
func (s *BidderService) GetBidderByID(id string) (*domain.Bidder, error) {
	bidder, err := s.bidderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidder: %w", err)
	}

	if bidder == nil {
		return nil, ErrBidderNotFound
	}

	return bidder, nil
}

// GetBidderList retrieves a paginated list of bidders with filters
func (s *BidderService) GetBidderList(req *domain.BidderListRequest) (*domain.BidderListResponse, error) {
	// Validate and set defaults
	if err := s.validateBidderListRequest(req); err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.bidderRepo.CountBiddersWithFilters(req)
	if err != nil {
		return nil, fmt.Errorf("failed to count bidders: %w", err)
	}

	// Get bidders
	bidders, err := s.bidderRepo.FindBiddersWithFilters(req)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidders: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	// Build response
	response := &domain.BidderListResponse{
		Bidders: bidders,
		Pagination: domain.Pagination{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// GrantPoints grants points to a bidder
func (s *BidderService) GrantPoints(bidderID string, points int64, adminID int64) (*domain.GrantPointsResponse, error) {
	// Validate points
	if points <= 0 {
		return nil, ErrInvalidPoints
	}

	if points > MaxPointsPerGrant {
		return nil, ErrPointsExceedMaximum
	}

	// Check if bidder exists
	bidder, err := s.bidderRepo.FindByID(bidderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidder: %w", err)
	}

	if bidder == nil {
		return nil, ErrBidderNotFound
	}

	// Check if bidder is deleted
	if bidder.IsDeleted() {
		return nil, errors.New("cannot grant points to deleted bidder")
	}

	// Grant points (within a transaction)
	updatedBidder, history, err := s.bidderRepo.GrantPoints(bidderID, points, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to grant points: %w", err)
	}

	// Build response
	response := &domain.GrantPointsResponse{
		Bidder:  *updatedBidder,
		History: *history,
	}

	return response, nil
}

// GetPointHistory retrieves the point history for a bidder
func (s *BidderService) GetPointHistory(bidderID string, page int, limit int) (*domain.PointHistoryListResponse, error) {
	// Validate and set defaults
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	// Check if bidder exists
	bidder, err := s.bidderRepo.FindByID(bidderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidder: %w", err)
	}

	if bidder == nil {
		return nil, ErrBidderNotFound
	}

	// Get total count
	total, err := s.bidderRepo.CountPointHistory(bidderID)
	if err != nil {
		return nil, fmt.Errorf("failed to count point history: %w", err)
	}

	// Get point history
	history, err := s.bidderRepo.GetPointHistory(bidderID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get point history: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Build response
	response := &domain.PointHistoryListResponse{
		Bidder:  *bidder,
		History: history,
		Pagination: domain.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// UpdateBidderStatus updates the status of a bidder account
func (s *BidderService) UpdateBidderStatus(id string, status domain.BidderStatus) (*domain.Bidder, error) {
	// Validate status value
	if !isValidBidderStatus(status) {
		return nil, ErrInvalidBidderStatus
	}

	// Check if bidder exists
	bidder, err := s.bidderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidder: %w", err)
	}

	if bidder == nil {
		return nil, ErrBidderNotFound
	}

	// Update status
	if err := s.bidderRepo.UpdateBidderStatus(id, status); err != nil {
		return nil, fmt.Errorf("failed to update bidder status: %w", err)
	}

	// Fetch updated bidder
	updatedBidder, err := s.bidderRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated bidder: %w", err)
	}

	return updatedBidder, nil
}

// validateBidderListRequest validates and sets defaults for BidderListRequest
func (s *BidderService) validateBidderListRequest(req *domain.BidderListRequest) error {
	// Set default page if not specified
	if req.Page <= 0 {
		req.Page = 1
	}

	// Set default limit if not specified
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// Validate limit range
	if req.Limit > 100 {
		return ErrInvalidLimit
	}

	// Validate statuses if specified
	for _, status := range req.Status {
		if !isValidBidderStatus(status) {
			return ErrInvalidBidderStatus
		}
	}

	// Validate sort mode if specified
	if req.Sort != "" && !isValidBidderSortMode(req.Sort) {
		return ErrInvalidBidderSortMode
	}

	return nil
}

// isValidBidderStatus checks if the bidder status is valid
func isValidBidderStatus(status domain.BidderStatus) bool {
	return status == domain.BidderStatusActive ||
		status == domain.BidderStatusSuspended ||
		status == domain.BidderStatusDeleted
}

// isValidBidderSortMode checks if the sort mode is valid for bidders
func isValidBidderSortMode(sort string) bool {
	validSorts := []string{
		"id_asc", "id_desc",
		"email_asc", "email_desc",
		"points_asc", "points_desc",
		"created_at_asc", "created_at_desc",
	}

	for _, validSort := range validSorts {
		if sort == validSort {
			return true
		}
	}

	return false
}
