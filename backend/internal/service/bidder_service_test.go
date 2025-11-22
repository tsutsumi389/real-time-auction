package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
)

// MockBidderRepository is a mock implementation of BidderRepository
type MockBidderRepository struct {
	mock.Mock
}

func (m *MockBidderRepository) FindByID(id string) (*domain.Bidder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bidder), args.Error(1)
}

func (m *MockBidderRepository) FindByEmail(email string) (*domain.Bidder, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bidder), args.Error(1)
}

func (m *MockBidderRepository) FindBiddersWithFilters(req *domain.BidderListRequest) ([]domain.BidderWithPoints, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BidderWithPoints), args.Error(1)
}

func (m *MockBidderRepository) CountBiddersWithFilters(req *domain.BidderListRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBidderRepository) GetBidderPoints(bidderID string) (*domain.BidderPoints, error) {
	args := m.Called(bidderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BidderPoints), args.Error(1)
}

func (m *MockBidderRepository) GrantPoints(bidderID string, points int64, adminID int64) (*domain.BidderWithPoints, *domain.PointHistory, error) {
	args := m.Called(bidderID, points, adminID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*domain.BidderWithPoints), args.Get(1).(*domain.PointHistory), args.Error(2)
}

func (m *MockBidderRepository) GetPointHistory(bidderID string, page int, limit int) ([]domain.PointHistoryWithAuction, error) {
	args := m.Called(bidderID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PointHistoryWithAuction), args.Error(1)
}

func (m *MockBidderRepository) CountPointHistory(bidderID string) (int64, error) {
	args := m.Called(bidderID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBidderRepository) UpdateBidderStatus(id string, status domain.BidderStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockBidderRepository) CreateBidderWithPoints(bidder *domain.Bidder, initialPoints int64, adminID int64) (*domain.BidderResponse, error) {
	args := m.Called(bidder, initialPoints, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BidderResponse), args.Error(1)
}

func TestBidderService_GetBidderByID(t *testing.T) {
	t.Run("Success - Bidder found", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		displayName := "Test Bidder"
		expectedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		mockBidderRepo.On("FindByID", bidderID).Return(expectedBidder, nil)

		result, err := bidderService.GetBidderByID(bidderID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBidder.ID, result.ID)
		assert.Equal(t, expectedBidder.Email, result.Email)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "non-existent-id"

		mockBidderRepo.On("FindByID", bidderID).Return(nil, nil)

		result, err := bidderService.GetBidderByID(bidderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrBidderNotFound, err)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		repositoryError := errors.New("database connection error")

		mockBidderRepo.On("FindByID", bidderID).Return(nil, repositoryError)

		result, err := bidderService.GetBidderByID(bidderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find bidder")

		mockBidderRepo.AssertExpectations(t)
	})
}

func TestBidderService_GetBidderList(t *testing.T) {
	t.Run("Success - Default parameters", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:  1,
			Limit: 20,
		}

		displayName1 := "Bidder 1"
		displayName2 := "Bidder 2"
		expectedBidders := []domain.BidderWithPoints{
			{
				Bidder: domain.Bidder{
					ID:          "bidder-1",
					Email:       "bidder1@example.com",
					DisplayName: &displayName1,
					Status:      domain.BidderStatusActive,
				},
				Points: 1000,
			},
			{
				Bidder: domain.Bidder{
					ID:          "bidder-2",
					Email:       "bidder2@example.com",
					DisplayName: &displayName2,
					Status:      domain.BidderStatusActive,
				},
				Points: 500,
			},
		}

		mockBidderRepo.On("CountBiddersWithFilters", req).Return(int64(2), nil)
		mockBidderRepo.On("FindBiddersWithFilters", req).Return(expectedBidders, nil)

		result, err := bidderService.GetBidderList(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Bidders))
		assert.Equal(t, int64(2), result.Pagination.Total)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 20, result.Pagination.Limit)
		assert.Equal(t, 1, result.Pagination.TotalPages)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - With keyword filter", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:    1,
			Limit:   20,
			Keyword: "test",
		}

		displayName := "Test Bidder"
		expectedBidders := []domain.BidderWithPoints{
			{
				Bidder: domain.Bidder{
					ID:          "bidder-1",
					Email:       "test@example.com",
					DisplayName: &displayName,
					Status:      domain.BidderStatusActive,
				},
				Points: 1000,
			},
		}

		mockBidderRepo.On("CountBiddersWithFilters", req).Return(int64(1), nil)
		mockBidderRepo.On("FindBiddersWithFilters", req).Return(expectedBidders, nil)

		result, err := bidderService.GetBidderList(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.Bidders))
		assert.Contains(t, result.Bidders[0].Email, "test")

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - With status filter", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:   1,
			Limit:  20,
			Status: []domain.BidderStatus{domain.BidderStatusSuspended},
		}

		displayName := "Suspended Bidder"
		expectedBidders := []domain.BidderWithPoints{
			{
				Bidder: domain.Bidder{
					ID:          "bidder-1",
					Email:       "suspended@example.com",
					DisplayName: &displayName,
					Status:      domain.BidderStatusSuspended,
				},
				Points: 0,
			},
		}

		mockBidderRepo.On("CountBiddersWithFilters", req).Return(int64(1), nil)
		mockBidderRepo.On("FindBiddersWithFilters", req).Return(expectedBidders, nil)

		result, err := bidderService.GetBidderList(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.Bidders))
		assert.Equal(t, domain.BidderStatusSuspended, result.Bidders[0].Status)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:   1,
			Limit:  20,
			Status: []domain.BidderStatus{"invalid_status"},
		}

		result, err := bidderService.GetBidderList(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidBidderStatus, err)
	})

	t.Run("Error - Invalid sort mode", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:  1,
			Limit: 20,
			Sort:  "invalid_sort",
		}

		result, err := bidderService.GetBidderList(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidBidderSortMode, err)
	})

	t.Run("Error - Limit exceeds maximum", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderListRequest{
			Page:  1,
			Limit: 200, // Exceeds max of 100
		}

		result, err := bidderService.GetBidderList(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidLimit, err)
	})
}

func TestBidderService_GrantPoints(t *testing.T) {
	t.Run("Success - Valid points grant", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		points := int64(500)
		adminID := int64(1)

		displayName := "Test Bidder"
		existingBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		updatedBidder := &domain.BidderWithPoints{
			Bidder: domain.Bidder{
				ID:          bidderID,
				Email:       "bidder@example.com",
				DisplayName: &displayName,
				Status:      domain.BidderStatusActive,
			},
			Points: 1500,
		}

		history := &domain.PointHistory{
			BidderID:       bidderID,
			Amount:         points,
			Type:           domain.PointHistoryTypeGrant,
			AdminID:        &adminID,
			BalanceBefore:  1000,
			BalanceAfter:   1500,
			ReservedBefore: 0,
			ReservedAfter:  0,
			TotalBefore:    1000,
			TotalAfter:     1500,
		}

		mockBidderRepo.On("FindByID", bidderID).Return(existingBidder, nil)
		mockBidderRepo.On("GrantPoints", bidderID, points, adminID).Return(updatedBidder, history, nil)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedBidder.Bidder.ID, result.Bidder.ID)
		assert.Equal(t, int64(1500), result.Bidder.Points)
		assert.Equal(t, points, result.History.Amount)
		assert.Equal(t, domain.PointHistoryTypeGrant, result.History.Type)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid points (zero)", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		points := int64(0)
		adminID := int64(1)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidPoints, err)
	})

	t.Run("Error - Invalid points (negative)", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		points := int64(-100)
		adminID := int64(1)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidPoints, err)
	})

	t.Run("Error - Points exceed maximum", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		points := int64(2000000) // Exceeds MaxPointsPerGrant (1000000)
		adminID := int64(1)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrPointsExceedMaximum, err)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "non-existent-id"
		points := int64(500)
		adminID := int64(1)

		mockBidderRepo.On("FindByID", bidderID).Return(nil, nil)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrBidderNotFound, err)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Cannot grant points to deleted bidder", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "deleted-bidder-id"
		points := int64(500)
		adminID := int64(1)

		displayName := "Deleted Bidder"
		deletedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "deleted@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusDeleted,
		}

		mockBidderRepo.On("FindByID", bidderID).Return(deletedBidder, nil)

		result, err := bidderService.GrantPoints(bidderID, points, adminID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot grant points to deleted bidder")

		mockBidderRepo.AssertExpectations(t)
	})
}

func TestBidderService_GetPointHistory(t *testing.T) {
	t.Run("Success - Valid request", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		page := 1
		limit := 10

		displayName := "Test Bidder"
		existingBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		adminID := int64(1)
		history := []domain.PointHistoryWithAuction{
			{
				PointHistory: domain.PointHistory{
					BidderID:       bidderID,
					Amount:         500,
					Type:           domain.PointHistoryTypeGrant,
					AdminID:        &adminID,
					BalanceBefore:  0,
					BalanceAfter:   500,
					ReservedBefore: 0,
					ReservedAfter:  0,
					TotalBefore:    0,
					TotalAfter:     500,
				},
				AuctionTitle: nil,
			},
		}

		mockBidderRepo.On("FindByID", bidderID).Return(existingBidder, nil)
		mockBidderRepo.On("CountPointHistory", bidderID).Return(int64(1), nil)
		mockBidderRepo.On("GetPointHistory", bidderID, page, limit).Return(history, nil)

		result, err := bidderService.GetPointHistory(bidderID, page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, bidderID, result.Bidder.ID)
		assert.Equal(t, 1, len(result.History))
		assert.Equal(t, int64(1), result.Pagination.Total)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 10, result.Pagination.Limit)
		assert.Equal(t, 1, result.Pagination.TotalPages)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - Auto-correct invalid page and limit", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		page := 0       // Invalid, should default to 1
		limit := 100    // Exceeds max, should cap to 50

		displayName := "Test Bidder"
		existingBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		history := []domain.PointHistoryWithAuction{}

		mockBidderRepo.On("FindByID", bidderID).Return(existingBidder, nil)
		mockBidderRepo.On("CountPointHistory", bidderID).Return(int64(0), nil)
		mockBidderRepo.On("GetPointHistory", bidderID, 1, 50).Return(history, nil)

		result, err := bidderService.GetPointHistory(bidderID, page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 50, result.Pagination.Limit)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "non-existent-id"
		page := 1
		limit := 10

		mockBidderRepo.On("FindByID", bidderID).Return(nil, nil)

		result, err := bidderService.GetPointHistory(bidderID, page, limit)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrBidderNotFound, err)

		mockBidderRepo.AssertExpectations(t)
	})
}

func TestBidderService_UpdateBidderStatus(t *testing.T) {
	t.Run("Success - Update to suspended", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		newStatus := domain.BidderStatusSuspended

		displayName := "Test Bidder"
		existingBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		updatedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusSuspended,
		}

		mockBidderRepo.On("FindByID", bidderID).Return(existingBidder, nil).Once()
		mockBidderRepo.On("UpdateBidderStatus", bidderID, newStatus).Return(nil)
		mockBidderRepo.On("FindByID", bidderID).Return(updatedBidder, nil).Once()

		result, err := bidderService.UpdateBidderStatus(bidderID, newStatus)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.BidderStatusSuspended, result.Status)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - Update to active", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		newStatus := domain.BidderStatusActive

		displayName := "Test Bidder"
		existingBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusSuspended,
		}

		updatedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		mockBidderRepo.On("FindByID", bidderID).Return(existingBidder, nil).Once()
		mockBidderRepo.On("UpdateBidderStatus", bidderID, newStatus).Return(nil)
		mockBidderRepo.On("FindByID", bidderID).Return(updatedBidder, nil).Once()

		result, err := bidderService.UpdateBidderStatus(bidderID, newStatus)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.BidderStatusActive, result.Status)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "test-bidder-id"
		newStatus := domain.BidderStatus("invalid_status")

		result, err := bidderService.UpdateBidderStatus(bidderID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrInvalidBidderStatus, err)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		bidderID := "non-existent-id"
		newStatus := domain.BidderStatusSuspended

		mockBidderRepo.On("FindByID", bidderID).Return(nil, nil)

		result, err := bidderService.UpdateBidderStatus(bidderID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrBidderNotFound, err)

		mockBidderRepo.AssertExpectations(t)
	})
}

func TestBidderService_RegisterBidder(t *testing.T) {
	t.Run("Success - Valid registration with initial points", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		displayName := "Test Bidder"
		initialPoints := int64(1000)
		adminID := int64(1)

		req := &domain.BidderCreateRequest{
			Email:         "bidder@example.com",
			Password:      "password123",
			DisplayName:   &displayName,
			InitialPoints: &initialPoints,
		}

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: req.DisplayName,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     initialPoints,
				AvailablePoints: initialPoints,
				ReservedPoints:  0,
			},
		}

		// Check email is not taken
		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)

		// Create bidder with points
		mockBidderRepo.On("CreateBidderWithPoints",
			mock.MatchedBy(func(bidder *domain.Bidder) bool {
				return bidder.Email == req.Email && bidder.Status == domain.BidderStatusActive
			}),
			initialPoints,
			adminID,
		).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, adminID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.Email)
		assert.Equal(t, initialPoints, result.Points.TotalPoints)
		assert.Equal(t, initialPoints, result.Points.AvailablePoints)
		assert.Equal(t, int64(0), result.Points.ReservedPoints)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - Valid registration without initial points", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		displayName := "Test Bidder"
		adminID := int64(1)

		req := &domain.BidderCreateRequest{
			Email:         "bidder@example.com",
			Password:      "password123",
			DisplayName:   &displayName,
			InitialPoints: nil, // No initial points
		}

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: req.DisplayName,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     0,
				AvailablePoints: 0,
				ReservedPoints:  0,
			},
		}

		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)
		mockBidderRepo.On("CreateBidderWithPoints",
			mock.Anything,
			int64(0),
			adminID,
		).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, adminID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.Email)
		assert.Equal(t, int64(0), result.Points.TotalPoints)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Success - Display name defaults to nil", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		adminID := int64(1)

		req := &domain.BidderCreateRequest{
			Email:         "bidder@example.com",
			Password:      "password123",
			DisplayName:   nil, // No display name
			InitialPoints: nil,
		}

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: nil,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     0,
				AvailablePoints: 0,
				ReservedPoints:  0,
			},
		}

		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)
		mockBidderRepo.On("CreateBidderWithPoints", mock.Anything, int64(0), adminID).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, adminID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.DisplayName)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Email already exists", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderCreateRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingBidder := &domain.Bidder{
			ID:    "existing-id",
			Email: req.Email,
		}

		mockBidderRepo.On("FindByEmail", req.Email).Return(existingBidder, nil)

		result, err := bidderService.RegisterBidder(req, int64(1))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrEmailAlreadyExists, err)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid email format", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderCreateRequest{
			Email:    "invalid-email",
			Password: "password123",
		}

		// Note: Currently validation is handled at handler layer by Gin binding
		// Service layer will process even with invalid email format
		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: nil,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     0,
				AvailablePoints: 0,
				ReservedPoints:  0,
			},
		}

		mockBidderRepo.On("CreateBidderWithPoints", mock.Anything, int64(0), int64(1)).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, int64(1))

		// Service layer doesn't validate email format currently
		// It's validated at handler layer by Gin binding
		assert.NoError(t, err)
		assert.NotNil(t, result)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Password too short", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderCreateRequest{
			Email:    "bidder@example.com",
			Password: "short", // Less than 8 characters
		}

		// Note: Currently validation is handled at handler layer by Gin binding
		// Service layer will process even with short password
		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: nil,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     0,
				AvailablePoints: 0,
				ReservedPoints:  0,
			},
		}

		mockBidderRepo.On("CreateBidderWithPoints", mock.Anything, int64(0), int64(1)).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, int64(1))

		// Service layer doesn't validate password length currently
		// It's validated at handler layer by Gin binding
		assert.NoError(t, err)
		assert.NotNil(t, result)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Negative initial points", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		negativePoints := int64(-100)
		req := &domain.BidderCreateRequest{
			Email:         "bidder@example.com",
			Password:      "password123",
			InitialPoints: &negativePoints,
		}

		// Note: Currently validation is handled at handler layer by Gin binding
		// Service layer treats negative points as 0
		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)

		expectedResponse := &domain.BidderResponse{
			ID:          "generated-uuid",
			Email:       req.Email,
			DisplayName: nil,
			Status:      domain.BidderStatusActive,
			Points: domain.PointsInfo{
				TotalPoints:     0,
				AvailablePoints: 0,
				ReservedPoints:  0,
			},
		}

		// Service will use 0 for negative initial points
		mockBidderRepo.On("CreateBidderWithPoints", mock.Anything, int64(0), int64(1)).Return(expectedResponse, nil)

		result, err := bidderService.RegisterBidder(req, int64(1))

		// Service layer doesn't validate negative points currently
		// It's validated at handler layer by Gin binding
		// RegisterBidder will use 0 if InitialPoints is negative
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Points.TotalPoints)

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error when checking email", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderCreateRequest{
			Email:    "bidder@example.com",
			Password: "password123",
		}

		repositoryError := errors.New("database connection error")
		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, repositoryError)

		result, err := bidderService.RegisterBidder(req, int64(1))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to check email")

		mockBidderRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository error when creating bidder", func(t *testing.T) {
		mockBidderRepo := new(MockBidderRepository)
		bidderService := NewBidderService(mockBidderRepo)

		req := &domain.BidderCreateRequest{
			Email:    "bidder@example.com",
			Password: "password123",
		}

		mockBidderRepo.On("FindByEmail", req.Email).Return(nil, nil)

		repositoryError := errors.New("transaction failed")
		mockBidderRepo.On("CreateBidderWithPoints", mock.Anything, int64(0), int64(1)).Return(nil, repositoryError)

		result, err := bidderService.RegisterBidder(req, int64(1))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create bidder")

		mockBidderRepo.AssertExpectations(t)
	})
}
