package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWTService instance
func NewJWTService(secretKey string) *JWTService {
	if secretKey == "" {
		// Use environment variable if no secret provided
		secretKey = os.Getenv("JWT_SECRET")
		if secretKey == "" {
			// Default secret for development (should be overridden in production)
			secretKey = "your-secure-random-secret-key-min-32-chars"
		}
	}
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateTokenForAdmin generates a JWT token for an admin user
func (s *JWTService) GenerateTokenForAdmin(admin *domain.Admin) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &domain.JWTClaims{
		UserID:      admin.ID,
		Email:       admin.Email,
		DisplayName: admin.DisplayName,
		Role:        admin.Role,
		UserType:    domain.UserTypeAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateTokenForBidder generates a JWT token for a bidder user
func (s *JWTService) GenerateTokenForBidder(bidder *domain.Bidder) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Get display name (use empty string if nil)
	displayName := ""
	if bidder.DisplayName != nil {
		displayName = *bidder.DisplayName
	}

	// Create claims
	claims := &domain.JWTClaims{
		UserID:      bidder.ID, // UUID string
		Email:       bidder.Email,
		DisplayName: displayName,
		UserType:    domain.UserTypeBidder,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
