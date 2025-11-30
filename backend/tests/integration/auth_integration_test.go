package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/handler"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"github.com/tsutsumi389/real-time-auction/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates a test database connection
// Note: This requires a test database to be running
func setupTestDB(t *testing.T) *gorm.DB {
	// Use environment variable or default to test database
	dsn := "host=postgres user=auction_user password=auction_pass_dev_only dbname=auction_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v", err)
	}

	return db
}

// setupTestRouter sets up a test Gin router with authentication endpoints
func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup repositories
	adminRepo := repository.NewAdminRepository(db)
	bidderRepo := repository.NewBidderRepository(db)

	// Setup services
	jwtService := service.NewJWTService("")
	authService := service.NewAuthService(adminRepo, bidderRepo, jwtService)

	// Setup handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/admin/login", authHandler.AdminLogin)
		}
	}

	return router
}

func TestAdminLoginIntegration(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	t.Run("Success - Login with valid credentials from seed data", func(t *testing.T) {
		// Use seed data credentials (password is password123 from migration)
		loginReq := map[string]string{
			"email":    "admin@example.com",
			"password": "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "admin@example.com", response.User.Email)
		assert.Equal(t, domain.RoleSystemAdmin, response.User.Role)
	})

	t.Run("Error - Invalid credentials", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "admin@example.com",
			"password": "wrongpassword",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response handler.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid email or password", response.Error)
	})

	t.Run("Error - Invalid email format", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "invalid-email",
			"password": "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Password too short", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "admin@example.com",
			"password": "short",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Suspended account", func(t *testing.T) {
		// Use suspended account from seed data
		loginReq := map[string]string{
			"email":    "suspended@example.com",
			"password": "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response handler.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Account is suspended", response.Error)
	})

	t.Run("Error - Non-existent account", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/admin/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
