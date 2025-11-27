package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// BidHandler handles bid-related HTTP requests
type BidHandler struct {
	pointService *service.PointService
}

// NewBidHandler creates a new BidHandler instance
func NewBidHandler(pointService *service.PointService) *BidHandler {
	return &BidHandler{
		pointService: pointService,
	}
}

// GetPoints handles GET /api/bidder/points
func (h *BidHandler) GetPoints(c *gin.Context) {
	// Get bidder ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	jwtClaims, ok := claims.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid token claims",
		})
		return
	}

	// Verify user type is bidder
	if jwtClaims.UserType != "bidder" {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error: "Forbidden: Only bidders can access this endpoint",
		})
		return
	}

	// Get bidder ID from claims
	bidderID, ok := jwtClaims.GetUserIDAsString()
	if !ok || bidderID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid bidder ID in token",
		})
		return
	}

	// Call service to get points
	response, err := h.pointService.GetPoints(bidderID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrPointsNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Points not found",
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
