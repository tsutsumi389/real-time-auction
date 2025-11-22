package domain

import (
	"time"
)

// BidderStatus represents the status of a bidder account
type BidderStatus string

const (
	BidderStatusActive    BidderStatus = "active"
	BidderStatusSuspended BidderStatus = "suspended"
	BidderStatusDeleted   BidderStatus = "deleted"
)

// Bidder represents a bidder user in the system
type Bidder struct {
	ID           string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string       `gorm:"type:varchar(255);not null" json:"email"`
	PasswordHash string       `gorm:"type:varchar(255);not null" json:"-"` // Never include in JSON responses
	DisplayName  *string      `gorm:"type:varchar(100)" json:"display_name"`
	Status       BidderStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Bidder model
func (Bidder) TableName() string {
	return "bidders"
}

// IsActive checks if the bidder account is active
func (b *Bidder) IsActive() bool {
	return b.Status == BidderStatusActive
}

// IsSuspended checks if the bidder account is suspended
func (b *Bidder) IsSuspended() bool {
	return b.Status == BidderStatusSuspended
}

// IsDeleted checks if the bidder account is deleted
func (b *Bidder) IsDeleted() bool {
	return b.Status == BidderStatusDeleted
}

// BidderPoints represents the points balance for a bidder
type BidderPoints struct {
	BidderID        string    `gorm:"type:uuid;primaryKey" json:"bidder_id"`
	TotalPoints     int64     `gorm:"not null;default:0" json:"total_points"`
	AvailablePoints int64     `gorm:"not null;default:0" json:"available_points"`
	ReservedPoints  int64     `gorm:"not null;default:0" json:"reserved_points"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for BidderPoints model
func (BidderPoints) TableName() string {
	return "bidder_points"
}

// PointHistoryType represents the type of point transaction
type PointHistoryType string

const (
	PointHistoryTypeGrant   PointHistoryType = "grant"
	PointHistoryTypeReserve PointHistoryType = "reserve"
	PointHistoryTypeRelease PointHistoryType = "release"
	PointHistoryTypeConsume PointHistoryType = "consume"
	PointHistoryTypeRefund  PointHistoryType = "refund"
)

// PointHistory represents a record of point transactions
type PointHistory struct {
	ID               int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	BidderID         string           `gorm:"type:uuid;not null" json:"bidder_id"`
	Amount           int64            `gorm:"not null" json:"amount"`
	Type             PointHistoryType `gorm:"type:varchar(50);not null" json:"type"`
	Reason           *string          `gorm:"type:text" json:"reason"`
	RelatedAuctionID *int64           `gorm:"type:bigint" json:"related_auction_id"`
	RelatedBidID     *int64           `gorm:"type:bigint" json:"related_bid_id"`
	AdminID          *int64           `gorm:"type:bigint" json:"admin_id"`
	BalanceBefore    int64            `gorm:"not null" json:"balance_before"`
	BalanceAfter     int64            `gorm:"not null" json:"balance_after"`
	ReservedBefore   int64            `gorm:"not null" json:"reserved_before"`
	ReservedAfter    int64            `gorm:"not null" json:"reserved_after"`
	TotalBefore      int64            `gorm:"not null" json:"total_before"`
	TotalAfter       int64            `gorm:"not null" json:"total_after"`
	CreatedAt        time.Time        `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for PointHistory model
func (PointHistory) TableName() string {
	return "point_history"
}

// BidderWithPoints represents a bidder with their points information
type BidderWithPoints struct {
	Bidder
	Points int64 `json:"points"` // total_points from bidder_points table
}

// BidderListRequest represents the request parameters for bidder list endpoint
type BidderListRequest struct {
	Page    int            `form:"page"`
	Limit   int            `form:"limit"`
	Keyword string         `form:"keyword"`
	Status  []BidderStatus `form:"status"`
	Sort    string         `form:"sort"`
}

// BidderListResponse represents the response for bidder list endpoint
type BidderListResponse struct {
	Bidders    []BidderWithPoints `json:"bidders"`
	Pagination Pagination         `json:"pagination"`
}

// GrantPointsRequest represents the request body for granting points to a bidder
type GrantPointsRequest struct {
	Points int64 `json:"points" binding:"required,min=1"`
}

// GrantPointsResponse represents the response for grant points endpoint
type GrantPointsResponse struct {
	Bidder  BidderWithPoints `json:"bidder"`
	History PointHistory     `json:"history"`
}

// PointHistoryWithAuction represents a point history entry with auction information
type PointHistoryWithAuction struct {
	PointHistory
	AuctionTitle *string `json:"auction_title"`
}

// PointHistoryListResponse represents the response for point history endpoint
type PointHistoryListResponse struct {
	Bidder     Bidder                    `json:"bidder"`
	History    []PointHistoryWithAuction `json:"history"`
	Pagination Pagination                `json:"pagination"`
}

// UpdateBidderStatusRequest represents the request body for status update
type UpdateBidderStatusRequest struct {
	Status BidderStatus `json:"status" binding:"required"`
}

// BidderCreateRequest represents the request body for creating a new bidder
type BidderCreateRequest struct {
	Email         string  `json:"email" binding:"required,email"`
	Password      string  `json:"password" binding:"required,min=8"`
	DisplayName   *string `json:"display_name"`
	InitialPoints *int64  `json:"initial_points" binding:"omitempty,min=0"`
}

// BidderResponse represents the response for a single bidder with points
type BidderResponse struct {
	ID          string        `json:"id"`
	Email       string        `json:"email"`
	DisplayName *string       `json:"display_name"`
	Status      BidderStatus  `json:"status"`
	Points      PointsInfo    `json:"points"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// PointsInfo represents the points information for a bidder
type PointsInfo struct {
	TotalPoints     int64 `json:"total_points"`
	AvailablePoints int64 `json:"available_points"`
	ReservedPoints  int64 `json:"reserved_points"`
}
