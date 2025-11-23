package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuctionStatus represents the status of an auction
type AuctionStatus string

const (
	AuctionStatusPending   AuctionStatus = "pending"   // Non-public (not visible to bidders)
	AuctionStatusActive    AuctionStatus = "active"    // Public and accepting bids
	AuctionStatusEnded     AuctionStatus = "ended"     // Ended (visible to bidders, read-only)
	AuctionStatusCancelled AuctionStatus = "cancelled" // Cancelled (visible to bidders, read-only)
)

// Auction represents an auction container that groups multiple items
type Auction struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title       string        `gorm:"type:varchar(200);not null" json:"title"`
	Description string        `gorm:"type:text" json:"description"`
	Status      AuctionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Auction model
func (Auction) TableName() string {
	return "auctions"
}

// AuctionWithItemCount represents an auction with item count
type AuctionWithItemCount struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      AuctionStatus `json:"status"`
	ItemCount   int64         `json:"item_count"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// AuctionListRequest represents the request parameters for auction list endpoint
type AuctionListRequest struct {
	Page          int           `form:"page"`
	Limit         int           `form:"limit"`
	Keyword       string        `form:"keyword"`
	Status        AuctionStatus `form:"status"`
	CreatedAfter  string        `form:"created_after"`  // YYYY-MM-DD format
	UpdatedBefore string        `form:"updated_before"` // YYYY-MM-DD format
	Sort          string        `form:"sort"`
}

// AuctionListResponse represents the response for auction list endpoint
type AuctionListResponse struct {
	Auctions   []AuctionWithItemCount `json:"auctions"`
	Pagination Pagination             `json:"pagination"`
}
