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

// BidderHandler handles bidder-related HTTP requests
type BidderHandler struct {
	bidderService service.BidderServiceInterface
}

// NewBidderHandler creates a new BidderHandler instance
func NewBidderHandler(bidderService service.BidderServiceInterface) *BidderHandler {
	return &BidderHandler{
		bidderService: bidderService,
	}
}

// RegisterBidder handles POST /api/admin/bidders
func (h *BidderHandler) RegisterBidder(c *gin.Context) {
	// Parse request body
	var req domain.BidderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Get admin ID from JWT claims
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

	// Get admin ID from claims
	adminID, ok := jwtClaims.GetUserIDAsInt64()
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid admin ID in token",
		})
		return
	}

	// Call service
	response, err := h.bidderService.RegisterBidder(&req, adminID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already exists",
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
	c.JSON(http.StatusCreated, response)
}

// GetBidderList handles GET /api/admin/bidders
func (h *BidderHandler) GetBidderList(c *gin.Context) {
	// Parse query parameters
	var req domain.BidderListRequest

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

	// Parse status filter (comma-separated values)
	statusStr := c.Query("status")
	if statusStr != "" {
		statusValues := strings.Split(statusStr, ",")
		req.Status = make([]domain.BidderStatus, 0, len(statusValues))
		for _, s := range statusValues {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				req.Status = append(req.Status, domain.BidderStatus(trimmed))
			}
		}
	}

	// Parse sort mode
	req.Sort = c.DefaultQuery("sort", "created_at_asc")

	// Call service
	response, err := h.bidderService.GetBidderList(&req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrInvalidPage),
			errors.Is(err, service.ErrInvalidLimit),
			errors.Is(err, service.ErrInvalidBidderSortMode),
			errors.Is(err, service.ErrInvalidBidderStatus):
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

// GrantPoints handles POST /api/admin/bidders/:id/points
func (h *BidderHandler) GrantPoints(c *gin.Context) {
	// Get bidder ID from URL parameter
	bidderID := c.Param("id")

	// Parse request body
	var req domain.GrantPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Get admin ID from JWT claims
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

	// Get admin ID from claims
	adminID, ok := jwtClaims.GetUserIDAsInt64()
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid admin ID in token",
		})
		return
	}

	// Call service
	response, err := h.bidderService.GrantPoints(bidderID, req.Points, adminID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrBidderNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Bidder not found",
			})
		case errors.Is(err, service.ErrInvalidPoints):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid points value",
			})
		case errors.Is(err, service.ErrPointsExceedMaximum):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Points exceed maximum limit",
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

// GetPointHistory handles GET /api/admin/bidders/:id/points/history
func (h *BidderHandler) GetPointHistory(c *gin.Context) {
	// Get bidder ID from URL parameter
	bidderID := c.Param("id")

	// Parse page (default: 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse limit (default: 10, max: 50)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Call service
	response, err := h.bidderService.GetPointHistory(bidderID, page, limit)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrBidderNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Bidder not found",
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

// UpdateBidderStatus handles PATCH /api/admin/bidders/:id/status
func (h *BidderHandler) UpdateBidderStatus(c *gin.Context) {
	// Get bidder ID from URL parameter
	bidderID := c.Param("id")

	// Parse request body
	var req domain.UpdateBidderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Call service
	updatedBidder, err := h.bidderService.UpdateBidderStatus(bidderID, req.Status)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrBidderNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Bidder not found",
			})
		case errors.Is(err, service.ErrInvalidBidderStatus):
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
	c.JSON(http.StatusOK, updatedBidder)
}

// GetBidderByID handles GET /api/admin/bidders/:id
func (h *BidderHandler) GetBidderByID(c *gin.Context) {
	// Get bidder ID from URL parameter
	bidderID := c.Param("id")

	// Call service
	response, err := h.bidderService.GetBidderDetail(bidderID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrBidderNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Bidder not found",
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

// UpdateBidder handles PUT /api/admin/bidders/:id
func (h *BidderHandler) UpdateBidder(c *gin.Context) {
	// Get bidder ID from URL parameter
	bidderID := c.Param("id")

	// Parse request body
	var req domain.BidderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Call service
	response, err := h.bidderService.UpdateBidder(bidderID, &req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrBidderNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Bidder not found",
			})
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already exists",
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
