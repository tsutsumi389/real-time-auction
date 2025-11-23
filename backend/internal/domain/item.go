package domain

import (
	"time"

	"github.com/google/uuid"
)

// Item represents an auction item
type Item struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AuctionID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_items_auction" json:"auction_id"`
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
