package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
)

// MockAdminRepositoryInterface is a mock implementation of AdminRepositoryInterface
type MockAdminRepositoryInterface struct {
	mock.Mock
}

func (m *MockAdminRepositoryInterface) FindByEmail(email string) (*domain.Admin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepositoryInterface) FindByID(id int64) (*domain.Admin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func (m *MockAdminRepositoryInterface) Create(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepositoryInterface) Update(admin *domain.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockAdminRepositoryInterface) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAdminRepositoryInterface) FindAdminsWithFilters(req *domain.AdminListRequest) ([]domain.Admin, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Admin), args.Error(1)
}

func (m *MockAdminRepositoryInterface) CountAdminsWithFilters(req *domain.AdminListRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAdminRepositoryInterface) UpdateAdminStatus(id int64, status domain.AdminStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func TestAdminService_GetAdminList(t *testing.T) {
	now := time.Now()

	t.Run("Success - Default parameters", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		expectedAdmins := []domain.Admin{
			{
				ID:          1,
				Email:       "admin1@example.com",
				DisplayName: "Admin 1",
				Role:        domain.RoleSystemAdmin,
				Status:      domain.StatusActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				Email:       "admin2@example.com",
				DisplayName: "Admin 2",
				Role:        domain.RoleAuctioneer,
				Status:      domain.StatusActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockRepo.On("CountAdminsWithFilters", req).Return(int64(50), nil)
		mockRepo.On("FindAdminsWithFilters", req).Return(expectedAdmins, nil)

		response, err := service.GetAdminList(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Admins, 2)
		assert.Equal(t, int64(50), response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.Limit)
		assert.Equal(t, 3, response.Pagination.TotalPages) // 50 / 20 = 2.5 -> 3

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - With filters", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:    1,
			Limit:   10,
			Keyword: "admin",
			Role:    domain.RoleSystemAdmin,
			Status:  []domain.AdminStatus{domain.StatusActive},
			Sort:    "email_asc",
		}

		expectedAdmins := []domain.Admin{
			{
				ID:          1,
				Email:       "admin@example.com",
				DisplayName: "Admin",
				Role:        domain.RoleSystemAdmin,
				Status:      domain.StatusActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockRepo.On("CountAdminsWithFilters", req).Return(int64(1), nil)
		mockRepo.On("FindAdminsWithFilters", req).Return(expectedAdmins, nil)

		response, err := service.GetAdminList(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Admins, 1)
		assert.Equal(t, int64(1), response.Pagination.Total)
		assert.Equal(t, 1, response.Pagination.TotalPages)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Empty result", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		expectedAdmins := []domain.Admin{}

		mockRepo.On("CountAdminsWithFilters", req).Return(int64(0), nil)
		mockRepo.On("FindAdminsWithFilters", req).Return(expectedAdmins, nil)

		response, err := service.GetAdminList(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Admins, 0)
		assert.Equal(t, int64(0), response.Pagination.Total)
		assert.Equal(t, 0, response.Pagination.TotalPages)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Set default page and limit", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  0, // Invalid page, should default to 1
			Limit: 0, // Invalid limit, should default to 20
		}

		// After validation, page and limit should be set to defaults
		expectedReq := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		mockRepo.On("CountAdminsWithFilters", expectedReq).Return(int64(0), nil)
		mockRepo.On("FindAdminsWithFilters", expectedReq).Return([]domain.Admin{}, nil)

		response, err := service.GetAdminList(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, req.Page)   // Verified that page was set to 1
		assert.Equal(t, 20, req.Limit) // Verified that limit was set to 20

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid limit (too large)", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 101, // Exceeds max limit of 100
		}

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, ErrInvalidLimit, err)
	})

	t.Run("Error - Invalid role", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
			Role:  domain.AdminRole("invalid_role"),
		}

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid role value")
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:   1,
			Limit:  20,
			Status: []domain.AdminStatus{domain.AdminStatus("invalid_status")},
		}

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, ErrInvalidStatus, err)
	})

	t.Run("Error - Invalid sort mode", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
			Sort:  "invalid_sort",
		}

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, ErrInvalidSortMode, err)
	})

	t.Run("Error - Count query failed", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		dbError := errors.New("database connection error")
		mockRepo.On("CountAdminsWithFilters", req).Return(int64(0), dbError)

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to count admins")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Find query failed", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		req := &domain.AdminListRequest{
			Page:  1,
			Limit: 20,
		}

		dbError := errors.New("database connection error")
		mockRepo.On("CountAdminsWithFilters", req).Return(int64(10), nil)
		mockRepo.On("FindAdminsWithFilters", req).Return(nil, dbError)

		response, err := service.GetAdminList(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to find admins")

		mockRepo.AssertExpectations(t)
	})
}

func TestAdminService_UpdateAdminStatus(t *testing.T) {
	now := time.Now()

	t.Run("Success - Update to suspended", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		newStatus := domain.StatusSuspended

		existingAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusActive,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		updatedAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusSuspended,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.On("FindByID", adminID).Return(existingAdmin, nil).Once()
		mockRepo.On("UpdateAdminStatus", adminID, newStatus).Return(nil)
		mockRepo.On("FindByID", adminID).Return(updatedAdmin, nil).Once()

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.StatusSuspended, result.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Update to active", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		newStatus := domain.StatusActive

		existingAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusSuspended,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		updatedAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusActive,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.On("FindByID", adminID).Return(existingAdmin, nil).Once()
		mockRepo.On("UpdateAdminStatus", adminID, newStatus).Return(nil)
		mockRepo.On("FindByID", adminID).Return(updatedAdmin, nil).Once()

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.StatusActive, result.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		invalidStatus := domain.AdminStatus("invalid_status")

		result, err := service.UpdateAdminStatus(adminID, invalidStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidStatus, err)
	})

	t.Run("Error - Admin not found", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(999)
		newStatus := domain.StatusSuspended

		mockRepo.On("FindByID", adminID).Return(nil, nil)

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrAdminNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - FindByID database error", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		newStatus := domain.StatusSuspended
		dbError := errors.New("database connection error")

		mockRepo.On("FindByID", adminID).Return(nil, dbError)

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find admin")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - UpdateAdminStatus failed", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		newStatus := domain.StatusSuspended

		existingAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusActive,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		updateError := errors.New("update failed")
		mockRepo.On("FindByID", adminID).Return(existingAdmin, nil)
		mockRepo.On("UpdateAdminStatus", adminID, newStatus).Return(updateError)

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to update admin status")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Fetch updated admin failed", func(t *testing.T) {
		mockRepo := new(MockAdminRepositoryInterface)
		service := NewAdminService(mockRepo)

		adminID := int64(1)
		newStatus := domain.StatusSuspended

		existingAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusActive,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		fetchError := errors.New("fetch failed")
		mockRepo.On("FindByID", adminID).Return(existingAdmin, nil).Once()
		mockRepo.On("UpdateAdminStatus", adminID, newStatus).Return(nil)
		mockRepo.On("FindByID", adminID).Return(nil, fetchError).Once()

		result, err := service.UpdateAdminStatus(adminID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to fetch updated admin")

		mockRepo.AssertExpectations(t)
	})
}

func TestValidationHelpers(t *testing.T) {
	t.Run("isValidRole", func(t *testing.T) {
		assert.True(t, isValidRole(domain.RoleSystemAdmin))
		assert.True(t, isValidRole(domain.RoleAuctioneer))
		assert.False(t, isValidRole(domain.AdminRole("invalid")))
	})

	t.Run("isValidStatus", func(t *testing.T) {
		assert.True(t, isValidStatus(domain.StatusActive))
		assert.True(t, isValidStatus(domain.StatusSuspended))
		assert.True(t, isValidStatus(domain.StatusDeleted))
		assert.False(t, isValidStatus(domain.AdminStatus("invalid")))
	})

	t.Run("isValidSortMode", func(t *testing.T) {
		assert.True(t, isValidSortMode("id_asc"))
		assert.True(t, isValidSortMode("id_desc"))
		assert.True(t, isValidSortMode("email_asc"))
		assert.True(t, isValidSortMode("email_desc"))
		assert.True(t, isValidSortMode("created_at_asc"))
		assert.True(t, isValidSortMode("created_at_desc"))
		assert.False(t, isValidSortMode("invalid_sort"))
		assert.False(t, isValidSortMode("name_asc"))
	})
}
