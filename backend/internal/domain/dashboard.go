package domain

import "time"

// DashboardStats represents the dashboard statistics
type DashboardStats struct {
	ActiveAuctions int64 `json:"active_auctions"` // Number of active auctions
	TodayBids      int64 `json:"today_bids"`      // Number of bids today
	TotalBidders   int64 `json:"total_bidders"`   // Total number of active bidders
	TotalPoints    int64 `json:"total_points"`    // Total points in circulation
}

// RecentBid represents a recent bid with related information
type RecentBid struct {
	ItemID      string    `json:"item_id"`
	ItemName    string    `json:"item_name"`
	AuctionID   string    `json:"auction_id"`
	AuctionName string    `json:"auction_name"`
	BidderID    string    `json:"bidder_id"`
	BidderName  string    `json:"bidder_name"`
	Price       int64     `json:"price"`
	BidAt       time.Time `json:"bid_at"`
}

// NewBidder represents a newly registered bidder
type NewBidder struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName *string   `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// EndedAuction represents a recently ended auction
type EndedAuction struct {
	ItemID      string     `json:"item_id"`
	ItemName    string     `json:"item_name"`
	AuctionID   string     `json:"auction_id"`
	AuctionName string     `json:"auction_name"`
	WinnerID    *string    `json:"winner_id"`
	WinnerName  *string    `json:"winner_name"`
	FinalPrice  *int64     `json:"final_price"`
	EndedAt     *time.Time `json:"ended_at"`
}

// DashboardActivities represents recent activities on the dashboard
type DashboardActivities struct {
	RecentBids      []RecentBid    `json:"recent_bids"`
	NewBidders      []NewBidder    `json:"new_bidders,omitempty"` // Only for system_admin
	EndedAuctions   []EndedAuction `json:"ended_auctions"`
}

// DashboardStatsResponse represents the response for dashboard stats endpoint
type DashboardStatsResponse struct {
	Stats DashboardStats `json:"stats"`
}

// DashboardActivitiesResponse represents the response for dashboard activities endpoint
type DashboardActivitiesResponse struct {
	Activities DashboardActivities `json:"activities"`
}
