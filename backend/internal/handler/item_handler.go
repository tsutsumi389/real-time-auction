package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// ItemHandler handles item management HTTP requests
type ItemHandler struct {
	itemService *service.ItemService
}

// NewItemHandler creates a new ItemHandler instance
func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// GetItemList handles GET /api/admin/items
func (h *ItemHandler) GetItemList(c *gin.Context) {
	// Parse query parameters
	status := strings.TrimSpace(c.DefaultQuery("status", "all"))
	search := strings.TrimSpace(c.Query("search"))

	// Parse page (default: 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse limit (default: 20, max: 100)
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Call service
	response, err := h.itemService.GetItemList(status, search, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid status filter. Valid values: all, assigned, unassigned",
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

// GetItemDetail handles GET /api/admin/items/:id
func (h *ItemHandler) GetItemDetail(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(itemID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Call service
	item, err := h.itemService.GetItemDetail(itemID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
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

// CreateItem handles POST /api/admin/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	// Parse request body
	var req domain.StandaloneItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem handles PUT /api/admin/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(itemID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Parse request body
	var req domain.StandaloneItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	item, err := h.itemService.UpdateItem(itemID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotEditable):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Item cannot be edited (already started)",
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

// DeleteItem handles DELETE /api/admin/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	// Get item ID from URL parameter
	itemID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(itemID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Call service
	err := h.itemService.DeleteItem(itemID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemAssignedToAuction):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot delete item that is assigned to an auction",
			})
		case errors.Is(err, service.ErrItemHasBids):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot delete item that has bids",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignItems handles POST /api/admin/auctions/:id/items/assign
func (h *ItemHandler) AssignItems(c *gin.Context) {
	// Get auction ID from URL parameter
	auctionID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(auctionID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid auction ID format",
		})
		return
	}

	// Parse request body
	var req domain.AssignItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	err := h.itemService.AssignItemsToAuction(auctionID, req.ItemIDs)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "One or more items not found",
			})
		case errors.Is(err, service.ErrItemAlreadyAssigned):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "One or more items are already assigned to an auction",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items assigned successfully"})
}

// UnassignItem handles DELETE /api/admin/auctions/:id/items/:itemId/unassign
func (h *ItemHandler) UnassignItem(c *gin.Context) {
	// Get auction ID and item ID from URL parameters
	auctionID := c.Param("id")
	itemID := c.Param("itemId")

	// Validate UUID formats
	if _, err := uuid.Parse(auctionID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid auction ID format",
		})
		return
	}
	if _, err := uuid.Parse(itemID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid item ID format",
		})
		return
	}

	// Call service
	err := h.itemService.UnassignItemFromAuction(auctionID, itemID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Item not found",
			})
		case errors.Is(err, service.ErrItemNotAssigned):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item is not assigned to any auction",
			})
		case errors.Is(err, service.ErrItemNotInAuction):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Item is not in this auction",
			})
		case errors.Is(err, service.ErrAuctionAlreadyStarted):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Cannot unassign items after auction has started",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
