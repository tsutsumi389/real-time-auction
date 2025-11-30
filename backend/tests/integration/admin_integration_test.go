package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/handler"
	"github.com/tsutsumi389/real-time-auction/internal/middleware"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"github.com/tsutsumi389/real-time-auction/internal/service"
	"gorm.io/gorm"
)

// setupAdminTestRouter sets up a test Gin router with admin endpoints
func setupAdminTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup repositories
	adminRepo := repository.NewAdminRepository(db)
	bidderRepo := repository.NewBidderRepository(db)

	// Setup services
	jwtService := service.NewJWTService("test-secret")
	authService := service.NewAuthService(adminRepo, bidderRepo, jwtService)
	adminService := service.NewAdminService(adminRepo)

	// Setup handlers
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Setup routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/admin/login", authHandler.AdminLogin)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			systemAdmin := protected.Group("")
			systemAdmin.Use(middleware.RequireSystemAdmin())
			{
				systemAdmin.GET("/admin/admins", adminHandler.GetAdminList)
				systemAdmin.PATCH("/admin/admins/:id/status", adminHandler.UpdateAdminStatus)
			}
		}
	}

	return router
}

// getSystemAdminToken logs in and returns a valid system admin token
func getSystemAdminToken(t *testing.T, router *gin.Engine) string {
	loginReq := map[string]string{
		"email":    "admin@example.com",
		"password": "password123",
	}

	reqBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get system admin token: %d - %s", w.Code, w.Body.String())
	}

	var response domain.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal login response: %v", err)
	}

	return response.Token
}

// getAuctioneerToken logs in and returns a valid auctioneer token
func getAuctioneerToken(t *testing.T, router *gin.Engine) string {
	loginReq := map[string]string{
		"email":    "auctioneer@example.com",
		"password": "password123",
	}

	reqBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get auctioneer token: %d - %s", w.Code, w.Body.String())
	}

	var response domain.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal login response: %v", err)
	}

	return response.Token
}

func TestGetAdminListIntegration(t *testing.T) {
	db := setupTestDB(t)
	router := setupAdminTestRouter(db)

	t.Run("Success - Get admin list with default parameters", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify response structure
		assert.NotNil(t, response.Admins)
		assert.NotNil(t, response.Pagination)
		assert.Greater(t, len(response.Admins), 0) // Should have seed data
		assert.Greater(t, response.Pagination.Total, int64(0))
	})

	t.Run("Success - Filter by role (system_admin)", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?role=system_admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// All admins should be system_admin
		for _, admin := range response.Admins {
			assert.Equal(t, domain.RoleSystemAdmin, admin.Role)
		}
	})

	t.Run("Success - Filter by role (auctioneer)", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?role=auctioneer", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// All admins should be auctioneer
		for _, admin := range response.Admins {
			assert.Equal(t, domain.RoleAuctioneer, admin.Role)
		}
	})

	t.Run("Success - Filter by status (active)", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=active", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// All admins should be active
		for _, admin := range response.Admins {
			assert.Equal(t, domain.StatusActive, admin.Status)
		}
	})

	t.Run("Success - Filter by status (suspended)", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=suspended", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// All admins should be suspended
		for _, admin := range response.Admins {
			assert.Equal(t, domain.StatusSuspended, admin.Status)
		}
	})

	t.Run("Success - Search by email", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?search=admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// At least the admin@example.com should be included
		assert.Greater(t, len(response.Admins), 0)
	})

	t.Run("Success - Sort by email ascending", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?sort=email_asc", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify sort order
		if len(response.Admins) > 1 {
			for i := 0; i < len(response.Admins)-1; i++ {
				assert.True(t, response.Admins[i].Email <= response.Admins[i+1].Email)
			}
		}
	})

	t.Run("Success - Sort by email descending", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?sort=email_desc", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify sort order
		if len(response.Admins) > 1 {
			for i := 0; i < len(response.Admins)-1; i++ {
				assert.True(t, response.Admins[i].Email >= response.Admins[i+1].Email)
			}
		}
	})

	t.Run("Success - Pagination", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		// Get first page
		req, _ := http.NewRequest(http.MethodGet, "/api/admins?page=1&limit=2", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.AdminListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.LessOrEqual(t, len(response.Admins), 2)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.Limit)
	})

	t.Run("Error - Invalid role", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?role=invalid", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Invalid status", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=invalid", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Invalid sort field", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?sort=invalid", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Invalid sort mode", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins?sort=invalid_mode", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - No authentication", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Auctioneer role forbidden", func(t *testing.T) {
		token := getAuctioneerToken(t, router)

		req, _ := http.NewRequest(http.MethodGet, "/api/admins", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestUpdateAdminStatusIntegration(t *testing.T) {
	db := setupTestDB(t)
	router := setupAdminTestRouter(db)

	t.Run("Success - Suspend active admin", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		// First, get an active admin ID from the list
		req, _ := http.NewRequest(http.MethodGet, "/api/admins?status=active&limit=1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var listResponse domain.AdminListResponse
		json.Unmarshal(w.Body.Bytes(), &listResponse)

		if len(listResponse.Admins) == 0 {
			t.Skip("No active admin found to test")
		}

		adminID := listResponse.Admins[0].ID

		// Update status
		updateReq := map[string]string{
			"status": "suspended",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admins/%d/status", adminID), bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Admin
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, adminID, response.ID)
		assert.Equal(t, domain.StatusSuspended, response.Status)

		// Revert the status for other tests
		updateReq = map[string]string{
			"status": "active",
		}
		reqBody, _ = json.Marshal(updateReq)
		req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admins/%d/status", adminID), bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	t.Run("Success - Activate suspended admin", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		// Use suspended admin from seed data
		req, _ := http.NewRequest(http.MethodGet, "/api/admins?search=suspended", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var listResponse domain.AdminListResponse
		json.Unmarshal(w.Body.Bytes(), &listResponse)

		if len(listResponse.Admins) == 0 {
			t.Skip("No suspended admin found to test")
		}

		adminID := listResponse.Admins[0].ID

		// Update status
		updateReq := map[string]string{
			"status": "active",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admins/%d/status", adminID), bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.Admin
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, adminID, response.ID)
		assert.Equal(t, domain.StatusActive, response.Status)

		// Revert the status
		updateReq = map[string]string{
			"status": "suspended",
		}
		reqBody, _ = json.Marshal(updateReq)
		req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/admins/%d/status", adminID), bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	t.Run("Error - Invalid admin ID", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		updateReq := map[string]string{
			"status": "suspended",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/99999/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error - Invalid status value", func(t *testing.T) {
		token := getSystemAdminToken(t, router)

		updateReq := map[string]string{
			"status": "invalid",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - No authentication", func(t *testing.T) {
		updateReq := map[string]string{
			"status": "suspended",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Auctioneer role forbidden", func(t *testing.T) {
		token := getAuctioneerToken(t, router)

		updateReq := map[string]string{
			"status": "suspended",
		}
		reqBody, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPatch, "/api/admins/1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
