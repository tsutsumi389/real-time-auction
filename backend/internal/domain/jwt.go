package domain

import "github.com/golang-jwt/jwt/v5"

// UserType represents the type of user (admin or bidder)
type UserType string

const (
	UserTypeAdmin  UserType = "admin"
	UserTypeBidder UserType = "bidder"
)

// JWTClaims represents the custom claims for JWT tokens
type JWTClaims struct {
	UserID      interface{} `json:"user_id"`          // int64 for admin, string for bidder
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	Role        AdminRole   `json:"role,omitempty"`   // Only for admin users
	UserType    UserType    `json:"user_type"`        // "admin" or "bidder"
	jwt.RegisteredClaims
}

// GetUserIDAsInt64 returns UserID as int64 (for admin users)
func (c *JWTClaims) GetUserIDAsInt64() (int64, bool) {
	if id, ok := c.UserID.(int64); ok {
		return id, true
	}
	if id, ok := c.UserID.(float64); ok {
		return int64(id), true
	}
	return 0, false
}

// GetUserIDAsString returns UserID as string (for bidder users)
func (c *JWTClaims) GetUserIDAsString() (string, bool) {
	if id, ok := c.UserID.(string); ok {
		return id, true
	}
	return "", false
}

// LoginResponse represents the response for login endpoints
type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

// UserInfo represents the user information in login response
type UserInfo struct {
	ID          interface{} `json:"id"`             // int64 for admin, string for bidder
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	Role        AdminRole   `json:"role,omitempty"` // Only for admin users
	UserType    UserType    `json:"user_type"`
	Points      *PointsInfo `json:"points,omitempty"` // Only for bidder users
}
