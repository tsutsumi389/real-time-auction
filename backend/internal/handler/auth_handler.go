package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest represents the request body for admin login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// AdminLogin handles admin login requests
// POST /api/auth/admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req LoginRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Trim whitespace from email
	req.Email = strings.TrimSpace(req.Email)

	// Call auth service
	response, err := h.authService.LoginAdmin(req.Email, req.Password)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid email or password",
			})
		case errors.Is(err, service.ErrAccountSuspended):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Account is suspended",
			})
		case errors.Is(err, service.ErrAccountDeleted):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Account is deleted",
			})
		default:
			// Log internal errors but don't expose details to client
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, response)
}
