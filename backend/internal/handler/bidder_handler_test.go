package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// MockBidderService is a mock implementation of BidderService
type MockBidderService struct {
	mock.Mock
}

func (m *MockBidderService) GetBidderByID(id string) (*domain.Bidder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bidder), args.Error(1)
}

func (m *MockBidderService) GetBidderList(req *domain.BidderListRequest) (*domain.BidderListResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BidderListResponse), args.Error(1)
}

func (m *MockBidderService) GrantPoints(bidderID string, points int64, adminID int64) (*domain.GrantPointsResponse, error) {
	args := m.Called(bidderID, points, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GrantPointsResponse), args.Error(1)
}

func (m *MockBidderService) GetPointHistory(bidderID string, page int, limit int) (*domain.PointHistoryListResponse, error) {
	args := m.Called(bidderID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointHistoryListResponse), args.Error(1)
}

func (m *MockBidderService) UpdateBidderStatus(id string, status domain.BidderStatus) (*domain.Bidder, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bidder), args.Error(1)
}

func TestBidderHandler_GetBidderList(t *testing.T) {
	t.Run("Success - Default parameters", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		displayName := "Bidder 1"
		expectedResponse := &domain.BidderListResponse{
			Bidders: []domain.BidderWithPoints{
				{
					Bidder: domain.Bidder{
						ID:          "bidder-1",
						Email:       "bidder1@example.com",
						DisplayName: &displayName,
						Status:      domain.BidderStatusActive,
					},
					Points: 1000,
				},
			},
			Pagination: domain.Pagination{
				Total:      1,
				Page:       1,
				Limit:      20,
				TotalPages: 1,
			},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return req.Page == 1 && req.Limit == 20
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.BidderListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Bidders))
		assert.Equal(t, "bidder1@example.com", response.Bidders[0].Email)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - With query parameters", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		displayName := "Test Bidder"
		expectedResponse := &domain.BidderListResponse{
			Bidders: []domain.BidderWithPoints{
				{
					Bidder: domain.Bidder{
						ID:          "bidder-1",
						Email:       "test@example.com",
						DisplayName: &displayName,
						Status:      domain.BidderStatusActive,
					},
					Points: 500,
				},
			},
			Pagination: domain.Pagination{
				Total:      1,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return req.Page == 1 && req.Limit == 10 && req.Keyword == "test" && req.Sort == "points_desc"
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?page=1&limit=10&keyword=test&sort=points_desc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - With status filter", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		displayName := "Suspended Bidder"
		expectedResponse := &domain.BidderListResponse{
			Bidders: []domain.BidderWithPoints{
				{
					Bidder: domain.Bidder{
						ID:          "bidder-1",
						Email:       "suspended@example.com",
						DisplayName: &displayName,
						Status:      domain.BidderStatusSuspended,
					},
					Points: 0,
				},
			},
			Pagination: domain.Pagination{
				Total:      1,
				Page:       1,
				Limit:      20,
				TotalPages: 1,
			},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return len(req.Status) == 1 && req.Status[0] == domain.BidderStatusSuspended
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?status=suspended", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - With multiple status filter", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		expectedResponse := &domain.BidderListResponse{
			Bidders:    []domain.BidderWithPoints{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 20, TotalPages: 0},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return len(req.Status) == 2
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?status=active,suspended", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - Invalid page defaults to 1", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		expectedResponse := &domain.BidderListResponse{
			Bidders:    []domain.BidderWithPoints{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 20, TotalPages: 0},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return req.Page == 1
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?page=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - Limit exceeds max, capped to 100", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		expectedResponse := &domain.BidderListResponse{
			Bidders:    []domain.BidderWithPoints{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 100, TotalPages: 0},
		}

		mockBidderService.On("GetBidderList", mock.MatchedBy(func(req *domain.BidderListRequest) bool {
			return req.Limit == 100
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?limit=200", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		mockBidderService.On("GetBidderList", mock.Anything).
			Return(nil, service.ErrInvalidBidderStatus)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders?status=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "invalid")

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Internal server error", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders", handler.GetBidderList)

		mockBidderService.On("GetBidderList", mock.Anything).
			Return(nil, errors.New("database connection error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Error)

		mockBidderService.AssertExpectations(t)
	})
}

func TestBidderHandler_GrantPoints(t *testing.T) {
	t.Run("Success - Valid request", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()

		// Add middleware to set claims
		router.Use(func(c *gin.Context) {
			claims := &domain.JWTClaims{
				UserID:   1,
				Email:    "admin@example.com",
				Role:     domain.RoleSystemAdmin,
				UserType: domain.UserTypeAdmin,
			}
			c.Set("claims", claims)
			c.Next()
		})

		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		bidderID := "test-bidder-id"
		points := int64(500)

		grantReq := domain.GrantPointsRequest{
			Points: points,
		}

		adminID := int64(1)
		displayName := "Test Bidder"
		expectedResponse := &domain.GrantPointsResponse{
			Bidder: domain.BidderWithPoints{
				Bidder: domain.Bidder{
					ID:          bidderID,
					Email:       "bidder@example.com",
					DisplayName: &displayName,
					Status:      domain.BidderStatusActive,
				},
				Points: 1500,
			},
			History: domain.PointHistory{
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
			},
		}

		mockBidderService.On("GrantPoints", bidderID, points, adminID).
			Return(expectedResponse, nil)

		reqBody, _ := json.Marshal(grantReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/"+bidderID+"/points", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.GrantPointsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1500), response.Bidder.Points)
		assert.Equal(t, points, response.History.Amount)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/test-id/points", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Error - No JWT claims", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		grantReq := domain.GrantPointsRequest{
			Points: 500,
		}

		reqBody, _ := json.Marshal(grantReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/test-id/points", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()

		router.Use(func(c *gin.Context) {
			claims := &domain.JWTClaims{
				UserID:   1,
				Email:    "admin@example.com",
				Role:     domain.RoleSystemAdmin,
				UserType: domain.UserTypeAdmin,
			}
			c.Set("claims", claims)
			c.Next()
		})

		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		bidderID := "non-existent-id"
		points := int64(500)

		grantReq := domain.GrantPointsRequest{
			Points: points,
		}

		mockBidderService.On("GrantPoints", bidderID, points, int64(1)).
			Return(nil, service.ErrBidderNotFound)

		reqBody, _ := json.Marshal(grantReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/"+bidderID+"/points", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bidder not found", response.Error)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Invalid points (caught by binding validation)", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()

		router.Use(func(c *gin.Context) {
			claims := &domain.JWTClaims{
				UserID:   1,
				Email:    "admin@example.com",
				Role:     domain.RoleSystemAdmin,
				UserType: domain.UserTypeAdmin,
			}
			c.Set("claims", claims)
			c.Next()
		})

		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		bidderID := "test-bidder-id"
		points := int64(-100)

		grantReq := domain.GrantPointsRequest{
			Points: points,
		}

		// Note: Service won't be called because Gin validation fails first
		reqBody, _ := json.Marshal(grantReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/"+bidderID+"/points", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Gin binding validation returns "Invalid request body"
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Error - Points exceed maximum", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()

		router.Use(func(c *gin.Context) {
			claims := &domain.JWTClaims{
				UserID:   1,
				Email:    "admin@example.com",
				Role:     domain.RoleSystemAdmin,
				UserType: domain.UserTypeAdmin,
			}
			c.Set("claims", claims)
			c.Next()
		})

		router.POST("/api/admin/bidders/:id/points", handler.GrantPoints)

		bidderID := "test-bidder-id"
		points := int64(2000000)

		grantReq := domain.GrantPointsRequest{
			Points: points,
		}

		mockBidderService.On("GrantPoints", bidderID, points, int64(1)).
			Return(nil, service.ErrPointsExceedMaximum)

		reqBody, _ := json.Marshal(grantReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/admin/bidders/"+bidderID+"/points", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Points exceed maximum limit", response.Error)

		mockBidderService.AssertExpectations(t)
	})
}

func TestBidderHandler_GetPointHistory(t *testing.T) {
	t.Run("Success - Valid request", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders/:id/points/history", handler.GetPointHistory)

		bidderID := "test-bidder-id"
		adminID := int64(1)
		displayName := "Test Bidder"

		expectedResponse := &domain.PointHistoryListResponse{
			Bidder: domain.Bidder{
				ID:          bidderID,
				Email:       "bidder@example.com",
				DisplayName: &displayName,
				Status:      domain.BidderStatusActive,
			},
			History: []domain.PointHistoryWithAuction{
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
			},
			Pagination: domain.Pagination{
				Total:      1,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
		}

		mockBidderService.On("GetPointHistory", bidderID, 1, 10).
			Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders/"+bidderID+"/points/history", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.PointHistoryListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.History))
		assert.Equal(t, bidderID, response.Bidder.ID)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - With pagination parameters", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders/:id/points/history", handler.GetPointHistory)

		bidderID := "test-bidder-id"
		displayName := "Test Bidder"

		expectedResponse := &domain.PointHistoryListResponse{
			Bidder: domain.Bidder{
				ID:          bidderID,
				Email:       "bidder@example.com",
				DisplayName: &displayName,
				Status:      domain.BidderStatusActive,
			},
			History: []domain.PointHistoryWithAuction{},
			Pagination: domain.Pagination{
				Total:      0,
				Page:       2,
				Limit:      20,
				TotalPages: 0,
			},
		}

		mockBidderService.On("GetPointHistory", bidderID, 2, 20).
			Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders/"+bidderID+"/points/history?page=2&limit=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - Limit exceeds max, capped to 50", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders/:id/points/history", handler.GetPointHistory)

		bidderID := "test-bidder-id"

		expectedResponse := &domain.PointHistoryListResponse{
			Bidder:     domain.Bidder{ID: bidderID},
			History:    []domain.PointHistoryWithAuction{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 50, TotalPages: 0},
		}

		mockBidderService.On("GetPointHistory", bidderID, 1, 50).
			Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders/"+bidderID+"/points/history?limit=100", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.GET("/api/admin/bidders/:id/points/history", handler.GetPointHistory)

		bidderID := "non-existent-id"

		mockBidderService.On("GetPointHistory", bidderID, 1, 10).
			Return(nil, service.ErrBidderNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/admin/bidders/"+bidderID+"/points/history", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bidder not found", response.Error)

		mockBidderService.AssertExpectations(t)
	})
}

func TestBidderHandler_UpdateBidderStatus(t *testing.T) {
	t.Run("Success - Update to suspended", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.PATCH("/api/admin/bidders/:id/status", handler.UpdateBidderStatus)

		bidderID := "test-bidder-id"
		statusReq := domain.UpdateBidderStatusRequest{
			Status: domain.BidderStatusSuspended,
		}

		displayName := "Test Bidder"
		expectedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusSuspended,
		}

		mockBidderService.On("UpdateBidderStatus", bidderID, domain.BidderStatusSuspended).
			Return(expectedBidder, nil)

		reqBody, _ := json.Marshal(statusReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admin/bidders/"+bidderID+"/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Bidder
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.BidderStatusSuspended, response.Status)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Success - Update to active", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.PATCH("/api/admin/bidders/:id/status", handler.UpdateBidderStatus)

		bidderID := "test-bidder-id"
		statusReq := domain.UpdateBidderStatusRequest{
			Status: domain.BidderStatusActive,
		}

		displayName := "Test Bidder"
		expectedBidder := &domain.Bidder{
			ID:          bidderID,
			Email:       "bidder@example.com",
			DisplayName: &displayName,
			Status:      domain.BidderStatusActive,
		}

		mockBidderService.On("UpdateBidderStatus", bidderID, domain.BidderStatusActive).
			Return(expectedBidder, nil)

		reqBody, _ := json.Marshal(statusReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admin/bidders/"+bidderID+"/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Bidder
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.BidderStatusActive, response.Status)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.PATCH("/api/admin/bidders/:id/status", handler.UpdateBidderStatus)

		req, _ := http.NewRequest(http.MethodPatch, "/api/admin/bidders/test-id/status", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Error - Bidder not found", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.PATCH("/api/admin/bidders/:id/status", handler.UpdateBidderStatus)

		bidderID := "non-existent-id"
		statusReq := domain.UpdateBidderStatusRequest{
			Status: domain.BidderStatusSuspended,
		}

		mockBidderService.On("UpdateBidderStatus", bidderID, domain.BidderStatusSuspended).
			Return(nil, service.ErrBidderNotFound)

		reqBody, _ := json.Marshal(statusReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admin/bidders/"+bidderID+"/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bidder not found", response.Error)

		mockBidderService.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockBidderService := new(MockBidderService)
		handler := NewBidderHandler(mockBidderService)
		router := setupTestRouter()
		router.PATCH("/api/admin/bidders/:id/status", handler.UpdateBidderStatus)

		bidderID := "test-bidder-id"
		statusReq := domain.UpdateBidderStatusRequest{
			Status: domain.BidderStatus("invalid_status"),
		}

		mockBidderService.On("UpdateBidderStatus", bidderID, domain.BidderStatus("invalid_status")).
			Return(nil, service.ErrInvalidBidderStatus)

		reqBody, _ := json.Marshal(statusReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admin/bidders/"+bidderID+"/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid status value", response.Error)

		mockBidderService.AssertExpectations(t)
	})
}
