package service

import (
	"errors"
	"fmt"

	"github.com/tsutsumi389/real-time-auction/backend/internal/domain"
	"github.com/tsutsumi389/real-time-auction/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountSuspended   = errors.New("account is suspended")
	ErrAccountDeleted     = errors.New("account is deleted")
)

// AuthService handles authentication logic
type AuthService struct {
	adminRepo  *repository.AdminRepository
	jwtService *JWTService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(adminRepo *repository.AdminRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		adminRepo:  adminRepo,
		jwtService: jwtService,
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
