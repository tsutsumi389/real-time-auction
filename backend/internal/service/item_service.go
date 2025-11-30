package service

import (
	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

// ItemService handles business logic for item management operations
type ItemService struct {
	itemRepo *repository.ItemRepository
}

// NewItemService creates a new ItemService instance
func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

// GetItemList retrieves a list of items with filters, search, and pagination
func (s *ItemService) GetItemList(status string, search string, page int, limit int) (*domain.ItemListResponse, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Validate status
	validStatuses := map[string]bool{
		"all":        true,
		"assigned":   true,
		"unassigned": true,
		"":           true, // empty means "all"
	}
	if !validStatuses[status] {
		return nil, ErrInvalidStatus
	}

	// Get items from repository
	items, err := s.itemRepo.FindItemsWithFilters(status, search, page, limit)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.itemRepo.CountItemsWithFilters(status, search)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &domain.ItemListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

// GetItemDetail retrieves detailed information for a single item
func (s *ItemService) GetItemDetail(itemID string) (*domain.ItemDetailResponse, error) {
	item, err := s.itemRepo.FindItemDetailByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	return item, nil
}

// CreateItem creates a new item without auction assignment
func (s *ItemService) CreateItem(req *domain.StandaloneItemRequest) (*domain.Item, error) {
	// Create item entity
	item := &domain.Item{
		Name:          req.Name,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		AuctionID:     nil, // Not assigned to any auction
		LotNumber:     0,   // No lot number when unassigned
	}

	// Create item in repository
	if err := s.itemRepo.CreateItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem updates an item's basic information
func (s *ItemService) UpdateItem(itemID string, req *domain.StandaloneItemRequest) (*domain.Item, error) {
	// Check if item exists
	existingItem, err := s.itemRepo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if existingItem == nil {
		return nil, ErrItemNotFound
	}

	// Check if item is editable (not started)
	if existingItem.StartedAt != nil {
		return nil, ErrItemNotEditable
	}

	// Update item
	item, err := s.itemRepo.UpdateItem(itemID, req)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem deletes an item
func (s *ItemService) DeleteItem(itemID string) error {
	// Check if item exists
	item, err := s.itemRepo.FindItemByID(itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	// Check if item is assigned to an auction
	if item.AuctionID != nil {
		return ErrItemAssignedToAuction
	}

	// Check if item has bids
	bidCount, err := s.itemRepo.GetBidCountByItemID(itemID)
	if err != nil {
		return err
	}
	if bidCount > 0 {
		return ErrItemHasBids
	}

	// Delete item
	return s.itemRepo.DeleteItem(itemID)
}

// AssignItemsToAuction assigns multiple items to an auction
func (s *ItemService) AssignItemsToAuction(auctionID string, itemIDs []uuid.UUID) error {
	// Verify all items are unassigned
	for _, itemID := range itemIDs {
		item, err := s.itemRepo.FindItemByID(itemID.String())
		if err != nil {
			return err
		}
		if item == nil {
			return ErrItemNotFound
		}
		if item.AuctionID != nil {
			return ErrItemAlreadyAssigned
		}
	}

	// Assign items with auto-incrementing lot numbers
	return s.itemRepo.AssignItemsToAuction(auctionID, itemIDs)
}

// UnassignItemFromAuction removes an item from an auction
func (s *ItemService) UnassignItemFromAuction(auctionID string, itemID string) error {
	// Check if any item in the auction has started
	hasStarted, err := s.itemRepo.HasAnyItemStarted(auctionID)
	if err != nil {
		return err
	}
	if hasStarted {
		return ErrAuctionAlreadyStarted
	}

	// Check if item exists and belongs to the auction
	item, err := s.itemRepo.FindItemByID(itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}
	if item.AuctionID == nil {
		return ErrItemNotAssigned
	}

	// Parse auction ID for comparison
	aucID, err := uuid.Parse(auctionID)
	if err != nil {
		return err
	}
	if *item.AuctionID != aucID {
		return ErrItemNotInAuction
	}

	// Unassign item
	return s.itemRepo.UnassignItemFromAuction(auctionID, itemID)
}
