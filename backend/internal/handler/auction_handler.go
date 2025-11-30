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

// AuctionHandler handles auction-related HTTP requests
type AuctionHandler struct {
	auctionService service.AuctionServiceInterface
}

// NewAuctionHandler creates a new AuctionHandler instance
func NewAuctionHandler(auctionService service.AuctionServiceInterface) *AuctionHandler {
	return &AuctionHandler{
		auctionService: auctionService,
	}
}

// GetAuctionList handles GET /api/admin/auctions
func (h *AuctionHandler) GetAuctionList(c *gin.Context) {
	// Parse query parameters
	var req domain.AuctionListRequest

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

	// Parse status filter
	statusStr := c.Query("status")
	if statusStr != "" {
		req.Status = domain.AuctionStatus(statusStr)
	}

	// Parse date filters
	req.CreatedAfter = strings.TrimSpace(c.Query("created_after"))
	req.UpdatedBefore = strings.TrimSpace(c.Query("updated_before"))

	// Parse sort mode
	req.Sort = c.DefaultQuery("sort", "created_at_desc")

	// Call service
	response, err := h.auctionService.GetAuctionList(&req)
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

// StartAuction handles POST /api/admin/auctions/:id/start
func (h *AuctionHandler) StartAuction(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Call service
	auction, err := h.auctionService.StartAuction(id)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrNoItemsInAuction):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "No items found in auction",
			})
		case errors.Is(err, service.ErrItemsMissingStartingPrice):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Items missing starting price",
			})
		case errors.Is(err, service.ErrAuctionNotPending):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction is not in pending status",
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
	c.JSON(http.StatusOK, auction)
}

// EndAuction handles POST /api/admin/auctions/:id/end
func (h *AuctionHandler) EndAuction(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Call service
	auction, err := h.auctionService.EndAuction(id)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotActive):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction is not in active status",
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
	c.JSON(http.StatusOK, auction)
}

// CancelAuction handles POST /api/admin/auctions/:id/cancel
func (h *AuctionHandler) CancelAuction(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Call service
	auction, err := h.auctionService.CancelAuction(id)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotActive):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction is not in active status",
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
	c.JSON(http.StatusOK, auction)
}

// CreateAuction handles POST /api/admin/auctions
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	// Parse request body
	var req domain.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	response, err := h.auctionService.CreateAuction(&req)
	if err != nil {
		// Log internal errors but don't expose details to client
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	// Return successful response with 201 Created
	c.JSON(http.StatusCreated, response)
}

// GetBidderAuctionList handles GET /api/auctions (public endpoint, no authentication required)
func (h *AuctionHandler) GetBidderAuctionList(c *gin.Context) {
	// Parse query parameters
	var req domain.BidderAuctionListRequest

	// Parse offset (default: 0)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	req.Offset = offset

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

	// Parse status filter (default: active)
	statusStr := c.DefaultQuery("status", "active")
	if statusStr != "" {
		req.Status = domain.AuctionStatus(statusStr)
	}

	// Parse sort mode (default: started_at_desc)
	req.Sort = c.DefaultQuery("sort", "started_at_desc")

	// Call service
	response, err := h.auctionService.GetBidderAuctionList(&req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrInvalidSortMode),
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

// GetAuctionDetail handles GET /api/auctions/:id
func (h *AuctionHandler) GetAuctionDetail(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Call service
	auction, err := h.auctionService.GetAuctionDetail(id)
	if err != nil {
		// Handle different error types
		if errors.Is(err, service.ErrAuctionNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// StartItem handles POST /api/items/:id/start
func (h *AuctionHandler) StartItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Call service
	response, err := h.auctionService.StartItem(itemID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemAlreadyStarted):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item already started",
			})
		case errors.Is(err, service.ErrStartingPriceNotSet):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Starting price not set",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// OpenPrice handles POST /api/items/:id/open-price
func (h *AuctionHandler) OpenPrice(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Get admin ID from context (set by auth middleware)
	adminIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	adminID, ok := adminIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid admin ID",
		})
		return
	}

	// Parse request body
	var req domain.OpenPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	response, err := h.auctionService.OpenPrice(itemID, req.NewPrice, adminID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotStarted):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item not started",
			})
		case errors.Is(err, service.ErrItemAlreadyEnded):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item already ended",
			})
		case errors.Is(err, service.ErrPriceTooLow):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "New price must be higher than current price",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// EndItem handles POST /api/items/:id/end
func (h *AuctionHandler) EndItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Call service
	response, err := h.auctionService.EndItem(itemID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotStarted):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item not started",
			})
		case errors.Is(err, service.ErrItemAlreadyEnded):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item already ended",
			})
		case errors.Is(err, service.ErrNoBidsFound):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "No bids found for this item",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBidHistory handles GET /api/items/:id/bids
func (h *AuctionHandler) GetBidHistory(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Call service
	response, err := h.auctionService.GetBidHistory(itemID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPriceHistory handles GET /api/items/:id/price-history
func (h *AuctionHandler) GetPriceHistory(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Call service
	response, err := h.auctionService.GetPriceHistory(itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetParticipants handles GET /api/auctions/:id/participants
func (h *AuctionHandler) GetParticipants(c *gin.Context) {
	// Get auction ID from URL parameter
	auctionID := c.Param("id")

	// Call service
	response, err := h.auctionService.GetParticipants(auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CancelAuctionWithReason handles POST /api/auctions/:id/cancel with reason
func (h *AuctionHandler) CancelAuctionWithReason(c *gin.Context) {
	// Get auction ID from URL parameter
	auctionID := c.Param("id")

	// Parse request body
	var req domain.CancelAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	response, err := h.auctionService.CancelAuctionWithReason(auctionID, req.Reason)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotActive):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction is not in active status",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAuctionForEdit handles GET /api/admin/auctions/:id (for edit)
func (h *AuctionHandler) GetAuctionForEdit(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Call service
	response, err := h.auctionService.GetAuctionForEdit(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}
	if response == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Auction not found",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAuction handles PUT /api/admin/auctions/:id
func (h *AuctionHandler) UpdateAuction(c *gin.Context) {
	// Get auction ID from URL parameter
	id := c.Param("id")

	// Parse request body
	var req domain.UpdateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	auction, err := h.auctionService.UpdateAuction(id, &req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotEditable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction cannot be edited",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, auction)
}

// UpdateItem handles PUT /api/admin/auctions/:id/items/:itemId
func (h *AuctionHandler) UpdateItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("itemId")

	// Parse request body
	var req domain.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	item, err := h.auctionService.UpdateItem(itemID, &req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotEditable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item cannot be edited",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem handles DELETE /api/admin/auctions/:id/items/:itemId
func (h *AuctionHandler) DeleteItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("itemId")

	// Call service
	err := h.auctionService.DeleteItem(itemID)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotDeletable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item cannot be deleted (already started)",
			})
		case errors.Is(err, service.ErrItemHasBids):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item has bids and cannot be deleted",
			})
		case errors.Is(err, service.ErrAuctionNotEditable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction cannot be edited",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddItem handles POST /api/admin/auctions/:id/items
func (h *AuctionHandler) AddItem(c *gin.Context) {
	// Get auction ID from URL parameter
	auctionID := c.Param("id")

	// Parse request body
	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	item, err := h.auctionService.AddItem(auctionID, &req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotEditable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction cannot be edited",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, item)
}

// ReorderItems handles PUT /api/admin/auctions/:id/items/reorder
func (h *AuctionHandler) ReorderItems(c *gin.Context) {
	// Get auction ID from URL parameter
	auctionID := c.Param("id")

	// Parse request body
	var req domain.ReorderItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	err := h.auctionService.ReorderItems(auctionID, &req)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrAuctionNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Auction not found",
			})
		case errors.Is(err, service.ErrAuctionNotEditable):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Auction cannot be edited",
			})
		case errors.Is(err, service.ErrInvalidItemIDs):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid item IDs",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items reordered successfully"})
}
