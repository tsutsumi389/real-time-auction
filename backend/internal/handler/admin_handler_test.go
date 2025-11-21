package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// MockAdminService is a mock implementation of AdminService
type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) GetAdminList(req *domain.AdminListRequest) (*domain.AdminListResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AdminListResponse), args.Error(1)
}

func (m *MockAdminService) UpdateAdminStatus(id int64, status domain.AdminStatus) (*domain.Admin, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Admin), args.Error(1)
}

func TestAdminHandler_GetAdminList(t *testing.T) {
	now := time.Now()

	t.Run("Success - Default parameters", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		expectedResponse := &domain.AdminListResponse{
			Admins: []domain.Admin{
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
			},
			Pagination: domain.Pagination{
				Total:      50,
				Page:       1,
				Limit:      20,
				TotalPages: 3,
			},
		}

		mockService.On("GetAdminList", mock.MatchedBy(func(req *domain.AdminListRequest) bool {
			return req.Page == 1 && req.Limit == 20
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Admins, 2)
		assert.Equal(t, int64(50), response.Pagination.Total)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - With query parameters", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		expectedResponse := &domain.AdminListResponse{
			Admins: []domain.Admin{
				{
					ID:          1,
					Email:       "admin@example.com",
					DisplayName: "Admin",
					Role:        domain.RoleSystemAdmin,
					Status:      domain.StatusActive,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			Pagination: domain.Pagination{
				Total:      1,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
		}

		mockService.On("GetAdminList", mock.MatchedBy(func(req *domain.AdminListRequest) bool {
			return req.Page == 1 &&
				req.Limit == 10 &&
				req.Keyword == "admin" &&
				req.Role == domain.RoleSystemAdmin &&
				len(req.Status) == 1 &&
				req.Status[0] == domain.StatusActive &&
				req.Sort == "email_asc"
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?page=1&limit=10&keyword=admin&role=system_admin&status=active&sort=email_asc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Admins, 1)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Multiple status filters", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		expectedResponse := &domain.AdminListResponse{
			Admins: []domain.Admin{
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
					Status:      domain.StatusSuspended,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			Pagination: domain.Pagination{
				Total:      2,
				Page:       1,
				Limit:      20,
				TotalPages: 1,
			},
		}

		mockService.On("GetAdminList", mock.MatchedBy(func(req *domain.AdminListRequest) bool {
			return len(req.Status) == 2 &&
				req.Status[0] == domain.StatusActive &&
				req.Status[1] == domain.StatusSuspended
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=active,suspended", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Invalid page defaults to 1", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		expectedResponse := &domain.AdminListResponse{
			Admins:     []domain.Admin{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 20, TotalPages: 0},
		}

		mockService.On("GetAdminList", mock.MatchedBy(func(req *domain.AdminListRequest) bool {
			return req.Page == 1
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?page=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Limit exceeds max, capped at 100", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		expectedResponse := &domain.AdminListResponse{
			Admins:     []domain.Admin{},
			Pagination: domain.Pagination{Total: 0, Page: 1, Limit: 100, TotalPages: 0},
		}

		mockService.On("GetAdminList", mock.MatchedBy(func(req *domain.AdminListRequest) bool {
			return req.Limit == 100
		})).Return(expectedResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?limit=200", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid page parameter (service error)", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		mockService.On("GetAdminList", mock.Anything).Return(nil, service.ErrInvalidPage)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "page")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid limit parameter (service error)", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		mockService.On("GetAdminList", mock.Anything).Return(nil, service.ErrInvalidLimit)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "limit")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid sort mode", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		mockService.On("GetAdminList", mock.Anything).Return(nil, service.ErrInvalidSortMode)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?sort=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "sort")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		mockService.On("GetAdminList", mock.Anything).Return(nil, service.ErrInvalidStatus)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "status")

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Internal server error", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.GET("/api/admins", handler.GetAdminList)

		mockService.On("GetAdminList", mock.Anything).Return(nil, errors.New("unexpected error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Error)

		mockService.AssertExpectations(t)
	})
}

func TestAdminHandler_UpdateAdminStatus(t *testing.T) {
	now := time.Now()

	t.Run("Success - Update to suspended", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		adminID := int64(1)
		newStatus := domain.StatusSuspended

		updateReq := domain.UpdateAdminStatusRequest{
			Status: newStatus,
		}

		expectedAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusSuspended,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockService.On("UpdateAdminStatus", adminID, newStatus).Return(expectedAdmin, nil)

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Admin
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.StatusSuspended, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("Success - Update to active", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		adminID := int64(1)
		newStatus := domain.StatusActive

		updateReq := domain.UpdateAdminStatusRequest{
			Status: newStatus,
		}

		expectedAdmin := &domain.Admin{
			ID:          adminID,
			Email:       "admin@example.com",
			DisplayName: "Test Admin",
			Role:        domain.RoleSystemAdmin,
			Status:      domain.StatusActive,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockService.On("UpdateAdminStatus", adminID, newStatus).Return(expectedAdmin, nil)

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Admin
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, domain.StatusActive, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid admin ID", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		updateReq := domain.UpdateAdminStatusRequest{
			Status: domain.StatusSuspended,
		}

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/invalid/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid admin ID", response.Error)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Error - Missing status field", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		reqBody := []byte(`{}`)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Error - Admin not found", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		adminID := int64(999)
		newStatus := domain.StatusSuspended

		updateReq := domain.UpdateAdminStatusRequest{
			Status: newStatus,
		}

		mockService.On("UpdateAdminStatus", adminID, newStatus).Return(nil, service.ErrAdminNotFound)

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/999/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Admin not found", response.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Invalid status value", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		adminID := int64(1)
		invalidStatus := domain.AdminStatus("invalid_status")

		updateReq := domain.UpdateAdminStatusRequest{
			Status: invalidStatus,
		}

		mockService.On("UpdateAdminStatus", adminID, invalidStatus).Return(nil, service.ErrInvalidStatus)

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid status value", response.Error)

		mockService.AssertExpectations(t)
	})

	t.Run("Error - Internal server error", func(t *testing.T) {
		mockService := new(MockAdminService)
		handler := NewAdminHandler(mockService)
		router := setupTestRouter()
		router.PATCH("/api/admins/:id/status", handler.UpdateAdminStatus)

		adminID := int64(1)
		newStatus := domain.StatusSuspended

		updateReq := domain.UpdateAdminStatusRequest{
			Status: newStatus,
		}

		mockService.On("UpdateAdminStatus", adminID, newStatus).Return(nil, errors.New("unexpected error"))

		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Error)

		mockService.AssertExpectations(t)
	})
}
