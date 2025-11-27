package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// BidHandler handles bid-related HTTP requests
type BidHandler struct {
	pointService *service.PointService
	bidService   *service.BidService
}

// NewBidHandler creates a new BidHandler instance
func NewBidHandler(pointService *service.PointService, bidService *service.BidService) *BidHandler {
	return &BidHandler{
		pointService: pointService,
		bidService:   bidService,
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

// PlaceBidRequest represents the request body for placing a bid
type PlaceBidRequest struct {
	Price int64 `json:"price" binding:"required,min=1"`
}

// PlaceBid handles POST /api/bidder/items/:id/bid
func (h *BidHandler) PlaceBid(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	// Parse request body
	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

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

	// Get bidder ID from claims
	bidderID, ok := jwtClaims.GetUserIDAsString()
	if !ok || bidderID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid bidder ID in token",
		})
		return
	}

	// Call service to place bid
	response, err := h.bidService.PlaceBid(&service.PlaceBidRequest{
		ItemID:   itemID,
		BidderID: bidderID,
		Price:    req.Price,
	})

	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotStarted):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item has not started yet",
			})
		case errors.Is(err, service.ErrItemAlreadyEnded):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item has already ended",
			})
		case errors.Is(err, service.ErrInsufficientPoints):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Insufficient points",
			})
		case errors.Is(err, service.ErrPriceMismatch):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Price has changed. Please check the latest price",
			})
		case errors.Is(err, service.ErrBidLockFailed):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Another bidder placed a bid first. Please try again",
			})
		case errors.Is(err, service.ErrAlreadyWinningBidder):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "You are already the winning bidder",
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

// GetBidHistory handles GET /api/bidder/items/:id/bids
func (h *BidHandler) GetBidHistory(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Item ID is required",
		})
		return
	}

	// Get query parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 200 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

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

	// Get bidder ID from claims
	bidderID, ok := jwtClaims.GetUserIDAsString()
	if !ok || bidderID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid bidder ID in token",
		})
		return
	}

	// Call service to get bid history
	response, err := h.bidService.GetBidHistory(itemID, bidderID, limit, offset)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
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
