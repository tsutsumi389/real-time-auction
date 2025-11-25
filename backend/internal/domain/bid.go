package domain

import (
	"time"

	"github.com/google/uuid"
)

// Bid represents a bid placed on an auction item
type Bid struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_bids_item" json:"item_id"`
	BidderID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_bids_bidder" json:"bidder_id"`
	Price     int64      `gorm:"type:bigint;not null" json:"price"`
	IsWinning bool       `gorm:"default:false;not null" json:"is_winning"`
	BidAt     time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"bid_at"`
}

// TableName specifies the table name for Bid model
func (Bid) TableName() string {
	return "bids"
}

// BidWithBidderInfo represents a bid with bidder information
type BidWithBidderInfo struct {
	ID          int64     `json:"id"`
	ItemID      uuid.UUID `json:"item_id"`
	BidderID    uuid.UUID `json:"bidder_id"`
	BidderName  string    `json:"bidder_name"`
	Price       int64     `json:"price"`
	IsWinning   bool      `json:"is_winning"`
	BidAt       time.Time `json:"bid_at"`
}

// BidHistoryResponse represents the response for bid history endpoint
type BidHistoryResponse struct {
	Total int64               `json:"total"`
	Bids  []BidWithBidderInfo `json:"bids"`
}
