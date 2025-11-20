package domain

import (
	"time"
)

// AdminRole represents the role of an admin user
type AdminRole string

const (
	RoleSystemAdmin AdminRole = "system_admin"
	RoleAuctioneer  AdminRole = "auctioneer"
)

// AdminStatus represents the status of an admin account
type AdminStatus string

const (
	StatusActive    AdminStatus = "active"
	StatusSuspended AdminStatus = "suspended"
	StatusDeleted   AdminStatus = "deleted"
)

// Admin represents an admin user in the system
type Admin struct {
	ID           int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string      `gorm:"type:varchar(255);not null" json:"-"` // Never include in JSON responses
	DisplayName  string      `gorm:"type:varchar(100)" json:"display_name"`
	Role         AdminRole   `gorm:"type:varchar(20);not null" json:"role"`
	Status       AdminStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Admin model
func (Admin) TableName() string {
	return "admins"
}

// IsActive checks if the admin account is active
func (a *Admin) IsActive() bool {
	return a.Status == StatusActive
}

// IsSuspended checks if the admin account is suspended
func (a *Admin) IsSuspended() bool {
	return a.Status == StatusSuspended
}

// IsDeleted checks if the admin account is deleted
func (a *Admin) IsDeleted() bool {
	return a.Status == StatusDeleted
}
