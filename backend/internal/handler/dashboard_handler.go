package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler instance
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStats handles GET /api/admin/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	// Get stats from service
	stats, err := h.dashboardService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve dashboard stats",
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, domain.DashboardStatsResponse{
		Stats: *stats,
	})
}

// GetActivities handles GET /api/admin/dashboard/activities
func (h *DashboardHandler) GetActivities(c *gin.Context) {
	// Get role from context (set by auth middleware)
	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	role, ok := roleValue.(domain.AdminRole)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid role format",
		})
		return
	}

	// Get activities from service with role-based filtering
	activities, err := h.dashboardService.GetActivities(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve dashboard activities",
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, domain.DashboardActivitiesResponse{
		Activities: *activities,
	})
}
