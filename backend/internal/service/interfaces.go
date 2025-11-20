package service

import "github.com/tsutsumi389/real-time-auction/internal/domain"

// AuthServiceInterface defines the interface for authentication service operations
type AuthServiceInterface interface {
	LoginAdmin(email, password string) (*domain.LoginResponse, error)
}

// JWTServiceInterface defines the interface for JWT service operations
type JWTServiceInterface interface {
	GenerateTokenForAdmin(admin *domain.Admin) (string, error)
	ValidateToken(tokenString string) (*domain.JWTClaims, error)
}
