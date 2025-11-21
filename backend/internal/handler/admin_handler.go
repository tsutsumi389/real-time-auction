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
			errors.Is(err, service.ErrInvalidStatus):
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
