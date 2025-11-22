package repository

import (
	"errors"
	"fmt"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// BidderRepository handles database operations for Bidder entities
type BidderRepository struct {
	db *gorm.DB
}

// NewBidderRepository creates a new BidderRepository instance
func NewBidderRepository(db *gorm.DB) *BidderRepository {
	return &BidderRepository{db: db}
}

// FindByID finds a bidder by ID
func (r *BidderRepository) FindByID(id string) (*domain.Bidder, error) {
	var bidder domain.Bidder
	result := r.db.Where("id = ?", id).First(&bidder)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &bidder, nil
}

// FindByEmail finds a bidder by email address
func (r *BidderRepository) FindByEmail(email string) (*domain.Bidder, error) {
	var bidder domain.Bidder
	result := r.db.Where("email = ?", email).First(&bidder)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &bidder, nil
}

// FindBiddersWithFilters retrieves bidders with filters, sorting, and pagination
func (r *BidderRepository) FindBiddersWithFilters(req *domain.BidderListRequest) ([]domain.BidderWithPoints, error) {
	var results []domain.BidderWithPoints

	query := r.db.Table("bidders b").
		Select("b.id, b.email, b.display_name, b.status, b.created_at, b.updated_at, COALESCE(bp.total_points, 0) as points").
		Joins("LEFT JOIN bidder_points bp ON b.id = bp.bidder_id")

	// Apply keyword filter (email or display_name search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("b.email ILIKE ? OR b.display_name ILIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Apply status filter (multiple statuses allowed)
	if len(req.Status) > 0 {
		query = query.Where("b.status IN ?", req.Status)
	} else {
		// Default: show only active and suspended (exclude deleted)
		query = query.Where("b.status IN ?", []domain.BidderStatus{domain.BidderStatusActive, domain.BidderStatusSuspended})
	}

	// Apply sorting
	switch req.Sort {
	case "id_asc":
		query = query.Order("b.id ASC")
	case "id_desc":
		query = query.Order("b.id DESC")
	case "email_asc":
		query = query.Order("b.email ASC")
	case "email_desc":
		query = query.Order("b.email DESC")
	case "points_asc":
		query = query.Order("points ASC")
	case "points_desc":
		query = query.Order("points DESC")
	case "created_at_asc":
		query = query.Order("b.created_at ASC")
	case "created_at_desc":
		query = query.Order("b.created_at DESC")
	default:
		// Default sort: created_at ASC
		query = query.Order("b.created_at ASC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit).Offset(offset)

	result := query.Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

// CountBiddersWithFilters counts the total number of bidders matching the filters
func (r *BidderRepository) CountBiddersWithFilters(req *domain.BidderListRequest) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Bidder{})

	// Apply keyword filter (email or display_name search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("email ILIKE ? OR display_name ILIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Apply status filter (multiple statuses allowed)
	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	} else {
		// Default: show only active and suspended (exclude deleted)
		query = query.Where("status IN ?", []domain.BidderStatus{domain.BidderStatusActive, domain.BidderStatusSuspended})
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// GetBidderPoints retrieves the points information for a bidder
func (r *BidderRepository) GetBidderPoints(bidderID string) (*domain.BidderPoints, error) {
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

// GrantPoints grants points to a bidder (within a transaction)
func (r *BidderRepository) GrantPoints(bidderID string, points int64, adminID int64) (*domain.BidderWithPoints, *domain.PointHistory, error) {
	var result domain.BidderWithPoints
	var history domain.PointHistory

	// Execute within a transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get current bidder points
		var currentPoints domain.BidderPoints
		if err := tx.Where("bidder_id = ?", bidderID).First(&currentPoints).Error; err != nil {
			return fmt.Errorf("failed to get current points: %w", err)
		}

		// Calculate new values
		newTotalPoints := currentPoints.TotalPoints + points
		newAvailablePoints := currentPoints.AvailablePoints + points

		// Update bidder_points
		if err := tx.Model(&domain.BidderPoints{}).
			Where("bidder_id = ?", bidderID).
			Updates(map[string]interface{}{
				"total_points":     newTotalPoints,
				"available_points": newAvailablePoints,
			}).Error; err != nil {
			return fmt.Errorf("failed to update bidder points: %w", err)
		}

		// Create point history record
		history = domain.PointHistory{
			BidderID:       bidderID,
			Amount:         points,
			Type:           domain.PointHistoryTypeGrant,
			AdminID:        &adminID,
			BalanceBefore:  currentPoints.AvailablePoints,
			BalanceAfter:   newAvailablePoints,
			ReservedBefore: currentPoints.ReservedPoints,
			ReservedAfter:  currentPoints.ReservedPoints,
			TotalBefore:    currentPoints.TotalPoints,
			TotalAfter:     newTotalPoints,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to create point history: %w", err)
		}

		// Get updated bidder with points
		if err := tx.Table("bidders b").
			Select("b.id, b.email, b.display_name, b.status, b.created_at, b.updated_at, COALESCE(bp.total_points, 0) as points").
			Joins("LEFT JOIN bidder_points bp ON b.id = bp.bidder_id").
			Where("b.id = ?", bidderID).
			Scan(&result).Error; err != nil {
			return fmt.Errorf("failed to get updated bidder: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &result, &history, nil
}

// GetPointHistory retrieves the point history for a bidder with pagination
func (r *BidderRepository) GetPointHistory(bidderID string, page int, limit int) ([]domain.PointHistoryWithAuction, error) {
	var results []domain.PointHistoryWithAuction

	offset := (page - 1) * limit

	query := r.db.Table("point_history ph").
		Select("ph.*, a.title as auction_title").
		Joins("LEFT JOIN auctions a ON ph.related_auction_id = a.id").
		Where("ph.bidder_id = ?", bidderID).
		Order("ph.created_at DESC").
		Limit(limit).
		Offset(offset)

	result := query.Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

// CountPointHistory counts the total number of point history entries for a bidder
func (r *BidderRepository) CountPointHistory(bidderID string) (int64, error) {
	var count int64

	result := r.db.Model(&domain.PointHistory{}).
		Where("bidder_id = ?", bidderID).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// UpdateBidderStatus updates the status of a bidder account
func (r *BidderRepository) UpdateBidderStatus(id string, status domain.BidderStatus) error {
	return r.db.Model(&domain.Bidder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// CreateBidderWithPoints creates a new bidder with initial points (within a transaction)
func (r *BidderRepository) CreateBidderWithPoints(bidder *domain.Bidder, initialPoints int64, adminID int64) (*domain.BidderResponse, error) {
	var response *domain.BidderResponse

	// Execute within a transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Create bidder
		if err := tx.Create(bidder).Error; err != nil {
			return fmt.Errorf("failed to create bidder: %w", err)
		}

		// Create bidder_points record manually (trigger has been removed)
		bidderPoints := &domain.BidderPoints{
			BidderID:        bidder.ID,
			TotalPoints:     initialPoints,
			AvailablePoints: initialPoints,
			ReservedPoints:  0,
		}

		if err := tx.Create(bidderPoints).Error; err != nil {
			return fmt.Errorf("failed to create bidder points: %w", err)
		}

		// Create point_history record if initial points > 0 (trigger has been removed)
		if initialPoints > 0 {
			reason := "初期ポイント付与"
			pointHistory := &domain.PointHistory{
				BidderID:       bidder.ID,
				Amount:         initialPoints,
				Type:           domain.PointHistoryTypeGrant,
				Reason:         &reason,
				AdminID:        &adminID,
				BalanceBefore:  0,
				BalanceAfter:   initialPoints,
				ReservedBefore: 0,
				ReservedAfter:  0,
				TotalBefore:    0,
				TotalAfter:     initialPoints,
			}

			if err := tx.Create(pointHistory).Error; err != nil {
				return fmt.Errorf("failed to create point history: %w", err)
			}
		}

		// Build response
		response = &domain.BidderResponse{
			ID:          bidder.ID,
			Email:       bidder.Email,
			DisplayName: bidder.DisplayName,
			Status:      bidder.Status,
			Points: domain.PointsInfo{
				TotalPoints:     bidderPoints.TotalPoints,
				AvailablePoints: bidderPoints.AvailablePoints,
				ReservedPoints:  bidderPoints.ReservedPoints,
			},
			CreatedAt: bidder.CreatedAt,
			UpdatedAt: bidder.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
