package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminService service.AdminServiceInterface
}

// NewAdminHandler creates a new AdminHandler instance
func NewAdminHandler(adminService service.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetAdminList handles GET /api/admins
func (h *AdminHandler) GetAdminList(c *gin.Context) {
	// Parse query parameters
	var req domain.AdminListRequest

	// Parse page (default: 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	req.Page = page

	// Parse limit (default: 20, max: 100)
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	req.Limit = limit

	// Parse keyword
	req.Keyword = strings.TrimSpace(c.Query("keyword"))

	// Parse role filter
	roleStr := c.Query("role")
	if roleStr != "" {
		req.Role = domain.AdminRole(roleStr)
	}

	// Parse status filter (comma-separated values)
	statusStr := c.Query("status")
	if statusStr != "" {
		statusValues := strings.Split(statusStr, ",")
		req.Status = make([]domain.AdminStatus, 0, len(statusValues))
		for _, s := range statusValues {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				req.Status = append(req.Status, domain.AdminStatus(trimmed))
			}
		}
	}

	// Parse sort mode
	req.Sort = c.DefaultQuery("sort", "id_asc")

	// Call service
	response, err := h.adminService.GetAdminList(&req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrInvalidPage),
			errors.Is(err, service.ErrInvalidLimit),
			errors.Is(err, service.ErrInvalidSortMode),
			errors.Is(err, service.ErrInvalidStatus),
			errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
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

// GetCurrentAdmin handles GET /api/admin/me
func (h *AdminHandler) GetCurrentAdmin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	// Convert to int64
	id, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	// Get admin from service
	admin, err := h.adminService.GetAdminByID(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAdminNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Admin not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	// Return admin data
	c.JSON(http.StatusOK, gin.H{
		"admin": admin,
	})
}

// RegisterAdmin handles POST /api/admins
func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	// Parse request body
	var req domain.AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	admin, err := h.adminService.RegisterAdmin(&req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already exists",
			})
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid role value",
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
	c.JSON(http.StatusCreated, admin)
}

// UpdateAdminStatus handles PATCH /api/admins/:id/status
func (h *AdminHandler) UpdateAdminStatus(c *gin.Context) {
	// Parse admin ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid admin ID",
		})
		return
	}

	// Parse request body
	var req domain.UpdateAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Call service
	updatedAdmin, err := h.adminService.UpdateAdminStatus(id, req.Status)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAdminNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Admin not found",
			})
		case errors.Is(err, service.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid status value",
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
	c.JSON(http.StatusOK, updatedAdmin)
}

// GetAdmin handles GET /api/admin/admins/:id
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	// Parse admin ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid admin ID",
		})
		return
	}

	// Get admin from service
	admin, err := h.adminService.GetAdminByID(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAdminNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Admin not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, domain.AdminDetailResponse{Admin: admin})
}

// UpdateAdmin handles PUT /api/admin/admins/:id
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	// Parse admin ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid admin ID",
		})
		return
	}

	// Get current user ID from context (set by auth middleware)
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	currentID, ok := currentUserID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req domain.AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	updatedAdmin, err := h.adminService.UpdateAdmin(id, &req, currentID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAdminNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Admin not found",
			})
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already exists",
			})
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid role value",
			})
		case errors.Is(err, service.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid status value",
			})
		case errors.Is(err, service.ErrCannotChangeOwnRole):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot change own role",
			})
		case errors.Is(err, service.ErrCannotSuspendSelf):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot suspend own account",
			})
		case errors.Is(err, service.ErrLastSystemAdmin):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot demote or suspend the last system admin",
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
	c.JSON(http.StatusOK, domain.AdminDetailResponse{Admin: updatedAdmin})
}
