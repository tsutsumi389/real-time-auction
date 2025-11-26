package service

import (
	"errors"
	"fmt"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountSuspended   = errors.New("account is suspended")
	ErrAccountDeleted     = errors.New("account is deleted")
)

// AuthService handles authentication logic
type AuthService struct {
	adminRepo   repository.AdminRepositoryInterface
	bidderRepo  repository.BidderRepositoryInterface
	jwtService  JWTServiceInterface
}

// NewAuthService creates a new AuthService instance
func NewAuthService(adminRepo repository.AdminRepositoryInterface, bidderRepo repository.BidderRepositoryInterface, jwtService JWTServiceInterface) *AuthService {
	return &AuthService{
		adminRepo:   adminRepo,
		bidderRepo:  bidderRepo,
		jwtService:  jwtService,
	}
}

// LoginAdmin authenticates an admin user and returns a JWT token
func (s *AuthService) LoginAdmin(email, password string) (*domain.LoginResponse, error) {
	// Find admin by email
	admin, err := s.adminRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	// Check if admin exists
	if admin == nil {
		return nil, ErrInvalidCredentials
	}

	// Check account status
	if admin.IsDeleted() {
		return nil, ErrAccountDeleted
	}

	if admin.IsSuspended() {
		return nil, ErrAccountSuspended
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateTokenForAdmin(admin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	response := &domain.LoginResponse{
		Token: token,
		User: &domain.UserInfo{
			ID:          admin.ID,
			Email:       admin.Email,
			DisplayName: admin.DisplayName,
			Role:        admin.Role,
			UserType:    domain.UserTypeAdmin,
		},
	}

	return response, nil
}

// LoginBidder authenticates a bidder user and returns a JWT token
func (s *AuthService) LoginBidder(email, password string) (*domain.LoginResponse, error) {
	// Find bidder by email
	bidder, err := s.bidderRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find bidder: %w", err)
	}

	// Check if bidder exists
	if bidder == nil {
		return nil, ErrInvalidCredentials
	}

	// Check account status
	if bidder.IsDeleted() {
		return nil, ErrAccountDeleted
	}

	if bidder.IsSuspended() {
		return nil, ErrAccountSuspended
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(bidder.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Get bidder points
	points, err := s.bidderRepo.GetBidderPoints(bidder.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bidder points: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateTokenForBidder(bidder)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Get display name (use empty string if nil)
	displayName := ""
	if bidder.DisplayName != nil {
		displayName = *bidder.DisplayName
	}

	// Build points info
	var pointsInfo *domain.PointsInfo
	if points != nil {
		pointsInfo = &domain.PointsInfo{
			TotalPoints:     points.TotalPoints,
			AvailablePoints: points.AvailablePoints,
			ReservedPoints:  points.ReservedPoints,
		}
	} else {
		// If no points record exists, return zero points
		pointsInfo = &domain.PointsInfo{
			TotalPoints:     0,
			AvailablePoints: 0,
			ReservedPoints:  0,
		}
	}

	// Build response
	response := &domain.LoginResponse{
		Token: token,
		User: &domain.UserInfo{
			ID:          bidder.ID, // UUID string
			Email:       bidder.Email,
			DisplayName: displayName,
			UserType:    domain.UserTypeBidder,
			Points:      pointsInfo,
		},
	}

	return response, nil
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// Use bcrypt cost of 10 (balance between security and performance)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ValidatePassword checks if a password meets the minimum requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
