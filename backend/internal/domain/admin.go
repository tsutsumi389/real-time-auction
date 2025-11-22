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

// Pagination represents pagination information for list responses
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// AdminListRequest represents the request parameters for admin list endpoint
type AdminListRequest struct {
	Page    int           `form:"page"`
	Limit   int           `form:"limit"`
	Keyword string        `form:"keyword"`
	Role    AdminRole     `form:"role"`
	Status  []AdminStatus `form:"status"`
	Sort    string        `form:"sort"`
}

// AdminListResponse represents the response for admin list endpoint
type AdminListResponse struct {
	Admins     []Admin    `json:"admins"`
	Pagination Pagination `json:"pagination"`
}

// AdminCreateRequest represents the request body for creating a new admin
type AdminCreateRequest struct {
	Email       string    `json:"email" binding:"required,email,max=255"`
	Password    string    `json:"password" binding:"required,min=8"`
	DisplayName string    `json:"display_name" binding:"omitempty,max=100"`
	Role        AdminRole `json:"role" binding:"required,oneof=system_admin auctioneer"`
}

// UpdateAdminStatusRequest represents the request body for status update
type UpdateAdminStatusRequest struct {
	Status AdminStatus `json:"status" binding:"required"`
}
