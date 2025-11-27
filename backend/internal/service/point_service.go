package service

import (
	"errors"
	"fmt"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

var (
	ErrPointsNotFound = errors.New("points not found for bidder")
)

// PointService handles point-related business logic
type PointService struct {
	pointRepo *repository.PointRepository
}

// NewPointService creates a new PointService instance
func NewPointService(pointRepo *repository.PointRepository) *PointService {
	return &PointService{
		pointRepo: pointRepo,
	}
}

// GetPoints retrieves the points information for a bidder
func (s *PointService) GetPoints(bidderID string) (*domain.GetPointsResponse, error) {
	// Get points from database
	points, err := s.pointRepo.FindPointsByBidderID(bidderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get points: %w", err)
	}

	if points == nil {
		return nil, ErrPointsNotFound
	}

	// Build response
	response := &domain.GetPointsResponse{
		BidderID:        points.BidderID,
		TotalPoints:     points.TotalPoints,
		AvailablePoints: points.AvailablePoints,
		ReservedPoints:  points.ReservedPoints,
		UpdatedAt:       points.UpdatedAt,
	}

	return response, nil
}
