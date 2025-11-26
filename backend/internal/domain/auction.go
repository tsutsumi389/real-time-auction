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
	StartedAt   *time.Time    `gorm:"type:timestamptz" json:"started_at"`
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
	StartedAt   *time.Time    `json:"started_at"`
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

// CreateItemRequest represents an item to be created with an auction
type CreateItemRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Description   string `json:"description" binding:"max=2000"`
	LotNumber     int    `json:"lot_number" binding:"required,min=1"`
	StartingPrice *int64 `json:"starting_price" binding:"omitempty,min=1"`
}

// CreateAuctionRequest represents the request body for creating an auction
type CreateAuctionRequest struct {
	Title       string              `json:"title" binding:"required,max=200"`
	Description string              `json:"description" binding:"max=2000"`
	StartedAt   *time.Time          `json:"started_at"`
	Items       []CreateItemRequest `json:"items" binding:"required,min=1,dive"`
}

// CreateAuctionResponse represents the response for creating an auction
type CreateAuctionResponse struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      AuctionStatus `json:"status"`
	StartedAt   *time.Time    `json:"started_at"`
	ItemCount   int           `json:"item_count"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// BidderPagination represents pagination information for bidder auction list (offset/limit style)
type BidderPagination struct {
	Total   int64 `json:"total"`
	Offset  int   `json:"offset"`
	Limit   int   `json:"limit"`
	HasMore bool  `json:"has_more"`
}

// BidderAuctionSummary represents auction summary information for bidders
type BidderAuctionSummary struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Status       AuctionStatus `json:"status"`
	ItemCount    int64         `json:"item_count"`
	ThumbnailURL *string       `json:"thumbnail_url"`
	StartedAt    *time.Time    `json:"started_at"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// BidderAuctionListRequest represents the request parameters for bidder auction list endpoint
type BidderAuctionListRequest struct {
	Offset  int           `form:"offset"`
	Limit   int           `form:"limit"`
	Keyword string        `form:"keyword"`
	Status  AuctionStatus `form:"status"`
	Sort    string        `form:"sort"`
}

// BidderAuctionListResponse represents the response for bidder auction list endpoint
type BidderAuctionListResponse struct {
	Auctions   []BidderAuctionSummary `json:"auctions"`
	Pagination BidderPagination       `json:"pagination"`
}

// GetAuctionDetailResponse represents the response for auction detail endpoint
type GetAuctionDetailResponse struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      AuctionStatus `json:"status"`
	StartedAt   *time.Time    `json:"started_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Items       []Item        `json:"items"`
}

// ParticipantInfo represents a participant in an auction
type ParticipantInfo struct {
	BidderID    uuid.UUID  `json:"bidder_id"`
	DisplayName string     `json:"display_name"`
	BidCount    int64      `json:"bid_count"`
	IsOnline    bool       `json:"is_online"`
	LastBidAt   *time.Time `json:"last_bid_at"`
}

// ParticipantsResponse represents the response for participants endpoint
type ParticipantsResponse struct {
	Total        int64             `json:"total"`
	Participants []ParticipantInfo `json:"participants"`
}
