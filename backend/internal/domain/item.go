package domain

import (
	"time"

	"github.com/google/uuid"
)

// Item represents an auction item
type Item struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AuctionID     *uuid.UUID `gorm:"type:uuid;index:idx_items_auction" json:"auction_id"`
	Name          string     `gorm:"type:varchar(200);not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	StartingPrice *int64     `gorm:"type:bigint" json:"starting_price"`
	CurrentPrice  *int64     `gorm:"type:bigint" json:"current_price"`
	WinnerID      *uuid.UUID `gorm:"type:uuid;index:idx_items_winner" json:"winner_id"`
	StartedAt     *time.Time `gorm:"index:idx_items_started_at" json:"started_at"`
	EndedAt       *time.Time `gorm:"index:idx_items_ended_at" json:"ended_at"`
	LotNumber     int        `gorm:"not null" json:"lot_number"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Item model
func (Item) TableName() string {
	return "items"
}

// ItemStatus represents the status of an auction item
type ItemStatus string

const (
	ItemStatusPending ItemStatus = "pending" // Not started yet
	ItemStatusActive  ItemStatus = "active"  // Currently active (started but not ended)
	ItemStatusEnded   ItemStatus = "ended"   // Ended
)

// ItemWithStatus represents an item with computed status field
type ItemWithStatus struct {
	Item
	Status ItemStatus  `json:"status"`
	Media  []ItemMedia `json:"media,omitempty"`
}

// GetStatus computes the status of an item based on timestamps
func (i *Item) GetStatus() ItemStatus {
	if i.EndedAt != nil {
		return ItemStatusEnded
	}
	if i.StartedAt != nil {
		return ItemStatusActive
	}
	return ItemStatusPending
}

// ToItemWithStatus converts an Item to ItemWithStatus
func (i *Item) ToItemWithStatus() ItemWithStatus {
	return ItemWithStatus{
		Item:   *i,
		Status: i.GetStatus(),
	}
}

// StartItemResponse represents the response for starting an item
type StartItemResponse struct {
	ItemID       uuid.UUID  `json:"item_id"`
	AuctionID    uuid.UUID  `json:"auction_id"`
	CurrentPrice int64      `json:"current_price"`
	StartedAt    time.Time  `json:"started_at"`
}

// OpenPriceRequest represents the request to open a new price
type OpenPriceRequest struct {
	NewPrice int64 `json:"new_price" binding:"required,min=1"`
}

// OpenPriceResponse represents the response for opening a new price
type OpenPriceResponse struct {
	ItemID        uuid.UUID     `json:"item_id"`
	CurrentPrice  int64         `json:"current_price"`
	PreviousPrice int64         `json:"previous_price"`
	DisclosedAt   time.Time     `json:"disclosed_at"`
	PriceHistory  *PriceHistory `json:"price_history"`
}

// EndItemResponse represents the response for ending an item
type EndItemResponse struct {
	ItemID     uuid.UUID  `json:"item_id"`
	WinnerID   *uuid.UUID `json:"winner_id"`
	WinnerName *string    `json:"winner_name"`
	FinalPrice int64      `json:"final_price"`
	EndedAt    time.Time  `json:"ended_at"`
}

// CancelAuctionRequest represents the request to cancel an auction
type CancelAuctionRequest struct {
	Reason string `json:"reason" binding:"max=500"`
}

// CancelAuctionResponse represents the response for cancelling an auction
type CancelAuctionResponse struct {
	AuctionID           uuid.UUID `json:"auction_id"`
	Status              string    `json:"status"`
	RefundedBidders     int64     `json:"refunded_bidders"`
	TotalRefundedPoints int64     `json:"total_refunded_points"`
	CancelledAt         time.Time `json:"cancelled_at"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Description *string `json:"description"`
}

// AddItemRequest represents the request to add a new item to an auction
type AddItemRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Description   string `json:"description" binding:"max=2000"`
	StartingPrice *int64 `json:"starting_price" binding:"omitempty,min=1"`
}

// ReorderItemsRequest represents the request to reorder items in an auction
type ReorderItemsRequest struct {
	ItemIDs []uuid.UUID `json:"item_ids" binding:"required,min=1"`
}

// ItemEditInfo represents an item with edit permissions
type ItemEditInfo struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	LotNumber     int        `json:"lot_number"`
	StartingPrice *int64     `json:"starting_price"`
	CurrentPrice  *int64     `json:"current_price"`
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	CanEdit       bool       `json:"can_edit"`
	CanDelete     bool       `json:"can_delete"`
	BidCount      int64      `json:"bid_count"`
}

// StandaloneItemRequest represents the request to create/update an item without auction assignment
type StandaloneItemRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Description   string `json:"description" binding:"max=2000"`
	StartingPrice *int64 `json:"starting_price" binding:"omitempty,min=1"`
}

// AssignItemsRequest represents the request to assign items to an auction
type AssignItemsRequest struct {
	ItemIDs []uuid.UUID `json:"item_ids" binding:"required,min=1"`
}

// UnassignItemRequest represents the request to unassign an item from an auction
// Note: This is typically used with URL path parameters, not request body
type UnassignItemRequest struct {
	ItemID uuid.UUID `json:"item_id"`
}

// ItemListResponse represents the response for item list with pagination
type ItemListResponse struct {
	Items []ItemListItem `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// ItemListItem represents an item in the list response
type ItemListItem struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	StartingPrice *int64     `json:"starting_price"`
	AuctionID     *uuid.UUID `json:"auction_id"`
	AuctionTitle  *string    `json:"auction_title"`
	LotNumber     int        `json:"lot_number"`
	BidCount      int64      `json:"bid_count"`
	CanDelete     bool       `json:"can_delete"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ItemDetailResponse represents the detailed response for a single item
type ItemDetailResponse struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	StartingPrice *int64     `json:"starting_price"`
	CurrentPrice  *int64     `json:"current_price"`
	AuctionID     *uuid.UUID `json:"auction_id"`
	AuctionTitle  *string    `json:"auction_title"`
	LotNumber     int        `json:"lot_number"`
	WinnerID      *uuid.UUID `json:"winner_id"`
	StartedAt     *time.Time `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	BidCount      int64      `json:"bid_count"`
	CanEdit       bool       `json:"can_edit"`
	CanDelete     bool       `json:"can_delete"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
