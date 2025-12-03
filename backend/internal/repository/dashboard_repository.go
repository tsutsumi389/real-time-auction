package repository

import (
	"time"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// DashboardRepository handles database operations for dashboard
type DashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new DashboardRepository instance
func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetStats retrieves dashboard statistics
func (r *DashboardRepository) GetStats() (*domain.DashboardStats, error) {
	var stats domain.DashboardStats

	// Count active auctions
	if err := r.db.Model(&domain.Auction{}).
		Where("status = ?", domain.AuctionStatusActive).
		Count(&stats.ActiveAuctions).Error; err != nil {
		return nil, err
	}

	// Count today's bids
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.Model(&domain.Bid{}).
		Where("bid_at >= ?", today).
		Count(&stats.TodayBids).Error; err != nil {
		return nil, err
	}

	// Count total active bidders
	if err := r.db.Model(&domain.Bidder{}).
		Where("status = ?", domain.BidderStatusActive).
		Count(&stats.TotalBidders).Error; err != nil {
		return nil, err
	}

	// Sum total points in circulation
	var totalPoints *int64
	if err := r.db.Model(&domain.BidderPoints{}).
		Select("COALESCE(SUM(total_points), 0)").
		Scan(&totalPoints).Error; err != nil {
		return nil, err
	}
	if totalPoints != nil {
		stats.TotalPoints = *totalPoints
	}

	return &stats, nil
}

// GetRecentBids retrieves the most recent bids (limit 5)
func (r *DashboardRepository) GetRecentBids(limit int) ([]domain.RecentBid, error) {
	var results []domain.RecentBid

	query := r.db.Table("bids b").
		Select(`b.item_id,
			i.name as item_name,
			a.id as auction_id,
			a.title as auction_name,
			b.bidder_id,
			bd.display_name as bidder_name,
			b.price,
			b.bid_at`).
		Joins("JOIN items i ON b.item_id = i.id").
		Joins("JOIN auctions a ON i.auction_id = a.id").
		Joins("JOIN bidders bd ON b.bidder_id = bd.id").
		Order("b.bid_at DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetNewBidders retrieves the most recently registered bidders (limit 5)
func (r *DashboardRepository) GetNewBidders(limit int) ([]domain.NewBidder, error) {
	var results []domain.NewBidder

	query := r.db.Model(&domain.Bidder{}).
		Select("id, email, display_name, created_at").
		Where("status = ?", domain.BidderStatusActive).
		Order("created_at DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetEndedAuctions retrieves recently ended items (limit 5)
func (r *DashboardRepository) GetEndedAuctions(limit int) ([]domain.EndedAuction, error) {
	var results []domain.EndedAuction

	query := r.db.Table("items i").
		Select(`i.id as item_id,
			i.name as item_name,
			a.id as auction_id,
			a.title as auction_name,
			i.winner_id,
			bd.display_name as winner_name,
			i.current_price as final_price,
			i.ended_at`).
		Joins("JOIN auctions a ON i.auction_id = a.id").
		Joins("LEFT JOIN bidders bd ON i.winner_id = bd.id").
		Where("i.ended_at IS NOT NULL").
		Order("i.ended_at DESC").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
