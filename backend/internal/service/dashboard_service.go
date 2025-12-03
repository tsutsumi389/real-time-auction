package service

import (
	"fmt"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

// DashboardService handles dashboard-related business logic
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

// NewDashboardService creates a new DashboardService instance
func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

// GetStats retrieves dashboard statistics
func (s *DashboardService) GetStats() (*domain.DashboardStats, error) {
	stats, err := s.dashboardRepo.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return stats, nil
}

// GetActivities retrieves recent activities with role-based filtering
func (s *DashboardService) GetActivities(role domain.AdminRole) (*domain.DashboardActivities, error) {
	activities := &domain.DashboardActivities{}

	// Get recent bids (all admins can see)
	recentBids, err := s.dashboardRepo.GetRecentBids(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent bids: %w", err)
	}
	activities.RecentBids = recentBids

	// Get ended auctions (all admins can see)
	endedAuctions, err := s.dashboardRepo.GetEndedAuctions(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get ended auctions: %w", err)
	}
	activities.EndedAuctions = endedAuctions

	// Get new bidders (only system_admin can see)
	if role == domain.RoleSystemAdmin {
		newBidders, err := s.dashboardRepo.GetNewBidders(5)
		if err != nil {
			return nil, fmt.Errorf("failed to get new bidders: %w", err)
		}
		activities.NewBidders = newBidders
	}

	return activities, nil
}
