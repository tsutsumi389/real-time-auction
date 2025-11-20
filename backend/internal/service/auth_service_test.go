package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// MockAdminRepository is a mock implementation of AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) FindByEmail(email string) (*domain.Admin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepository) FindByID(id int64) (*domain.Admin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepository) Create(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepository) Update(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateTokenForAdmin(admin *domain.Admin) (string, error) {
	args := m.Called(admin)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JWTClaims), args.Error(1)
}

func TestAuthService_LoginAdmin(t *testing.T) {
	t.Run("Success - Valid credentials", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		admin := &domain.Admin{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
			DisplayName:  "Test Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		expectedToken := "mock.jwt.token"

		mockAdminRepo.On("FindByEmail", email).Return(admin, nil)
		mockJWTService.On("GenerateTokenForAdmin", admin).Return(expectedToken, nil)

		response, err := authService.LoginAdmin(email, password)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedToken, response.Token)
		assert.Equal(t, admin.ID, response.User.ID)
		assert.Equal(t, admin.Email, response.User.Email)
		assert.Equal(t, admin.DisplayName, response.User.DisplayName)
		assert.Equal(t, admin.Role, response.User.Role)
		assert.Equal(t, domain.UserTypeAdmin, response.User.UserType)

		mockAdminRepo.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})

	t.Run("Error - Invalid email", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "notfound@example.com"
		password := "password123"

		mockAdminRepo.On("FindByEmail", email).Return(nil, nil)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid password", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "wrongpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), 10)

		admin := &domain.Admin{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
			DisplayName:  "Test Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusActive,
		}

		mockAdminRepo.On("FindByEmail", email).Return(admin, nil)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error - Account suspended", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		admin := &domain.Admin{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
			DisplayName:  "Test Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusSuspended,
		}

		mockAdminRepo.On("FindByEmail", email).Return(admin, nil)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Equal(t, ErrAccountSuspended, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error - Account deleted", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		admin := &domain.Admin{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
			DisplayName:  "Test Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusDeleted,
		}

		mockAdminRepo.On("FindByEmail", email).Return(admin, nil)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Equal(t, ErrAccountDeleted, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "password123"
		dbError := errors.New("database connection error")

		mockAdminRepo.On("FindByEmail", email).Return(nil, dbError)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("Error - JWT generation failed", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		mockJWTService := new(MockJWTService)
		authService := NewAuthService(mockAdminRepo, mockJWTService)

		email := "admin@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

		admin := &domain.Admin{
			ID:           1,
			Email:        email,
			PasswordHash: string(hashedPassword),
			DisplayName:  "Test Admin",
			Role:         domain.RoleSystemAdmin,
			Status:       domain.StatusActive,
		}

		jwtError := errors.New("jwt generation failed")

		mockAdminRepo.On("FindByEmail", email).Return(admin, nil)
		mockJWTService.On("GenerateTokenForAdmin", admin).Return("", jwtError)

		response, err := authService.LoginAdmin(email, password)

		assert.Error(t, err)
		assert.Nil(t, response)

		mockAdminRepo.AssertExpectations(t)
		mockJWTService.AssertExpectations(t)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)

		// Verify the hash can be compared with the original password
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		assert.NoError(t, err)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("Valid password", func(t *testing.T) {
		err := ValidatePassword("password123")
		assert.NoError(t, err)
	})

	t.Run("Invalid password - too short", func(t *testing.T) {
		err := ValidatePassword("pass")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 8 characters")
	})

	t.Run("Valid password - minimum length", func(t *testing.T) {
		err := ValidatePassword("12345678")
		assert.NoError(t, err)
	})
}
