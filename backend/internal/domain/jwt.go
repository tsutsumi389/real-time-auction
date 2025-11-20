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
	UserID      int64     `json:"user_id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        AdminRole `json:"role,omitempty"`   // Only for admin users
	UserType    UserType  `json:"user_type"`        // "admin" or "bidder"
	jwt.RegisteredClaims
}

// LoginResponse represents the response for login endpoints
type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

// UserInfo represents the user information in login response
type UserInfo struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        AdminRole `json:"role,omitempty"` // Only for admin users
	UserType    UserType  `json:"user_type"`
}
