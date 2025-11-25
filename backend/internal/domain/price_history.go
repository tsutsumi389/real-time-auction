package domain

import (
	"time"

	"github.com/google/uuid"
)

// PriceHistory represents the history of price disclosures for an auction item
type PriceHistory struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID      uuid.UUID `gorm:"type:uuid;not null;index:idx_price_history_item" json:"item_id"`
	Price       int64     `gorm:"type:bigint;not null" json:"price"`
	DisclosedBy int64     `gorm:"not null;index:idx_price_history_disclosed_by" json:"disclosed_by"`
	HadBid      bool      `gorm:"default:false;not null" json:"had_bid"`
	DisclosedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"disclosed_at"`
}

// TableName specifies the table name for PriceHistory model
func (PriceHistory) TableName() string {
	return "price_history"
}

// PriceHistoryWithAdmin represents price history with admin information
type PriceHistoryWithAdmin struct {
	ID              int64     `json:"id"`
	ItemID          uuid.UUID `json:"item_id"`
	Price           int64     `json:"price"`
	DisclosedByName string    `json:"disclosed_by_name"`
	HadBid          bool      `json:"had_bid"`
	DisclosedAt     time.Time `json:"disclosed_at"`
}

// PriceHistoryResponse represents the response for price history endpoint
type PriceHistoryResponse struct {
	Total   int64                   `json:"total"`
	History []PriceHistoryWithAdmin `json:"history"`
}
