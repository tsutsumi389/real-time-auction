package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"gorm.io/gorm"
)

// ItemRepository handles database operations for Item entities
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new ItemRepository instance
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// FindItemsWithFilters retrieves items with filters, search, and pagination
func (r *ItemRepository) FindItemsWithFilters(status string, search string, page int, limit int) ([]domain.ItemListItem, error) {
	var results []domain.ItemListItem

	query := r.db.Model(&domain.Item{}).
		Select(`items.id, items.name, items.description, items.starting_price,
			items.auction_id, items.lot_number, items.created_at,
			auctions.title as auction_title,
			(SELECT COUNT(*) FROM bids WHERE bids.item_id = items.id) as bid_count`).
		Joins("LEFT JOIN auctions ON auctions.id = items.auction_id")

	// Apply status filter
	switch status {
	case "assigned":
		query = query.Where("items.auction_id IS NOT NULL")
	case "unassigned":
		query = query.Where("items.auction_id IS NULL")
	// "all" or empty - no filter
	}

	// Apply keyword search
	if search != "" {
		query = query.Where("items.name ILIKE ?", "%"+search+"%")
	}

	// Order by created_at descending
	query = query.Order("items.created_at DESC")

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Execute query
	result := query.Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate can_delete for each item
	for i := range results {
		// Item can be deleted if not assigned to auction and has no bids
		results[i].CanDelete = results[i].AuctionID == nil && results[i].BidCount == 0
	}

	return results, nil
}

// CountItemsWithFilters counts items matching the filters
func (r *ItemRepository) CountItemsWithFilters(status string, search string) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Item{})

	// Apply status filter
	switch status {
	case "assigned":
		query = query.Where("auction_id IS NOT NULL")
	case "unassigned":
		query = query.Where("auction_id IS NULL")
	// "all" or empty - no filter
	}

	// Apply keyword search
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Execute count query
	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// FindItemDetailByID retrieves detailed item information
func (r *ItemRepository) FindItemDetailByID(itemID string) (*domain.ItemDetailResponse, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	var result domain.ItemDetailResponse

	query := r.db.Model(&domain.Item{}).
		Select(`items.id, items.name, items.description, items.starting_price,
			items.current_price, items.auction_id, items.lot_number,
			items.winner_id, items.started_at, items.ended_at,
			items.created_at, items.updated_at,
			auctions.title as auction_title,
			(SELECT COUNT(*) FROM bids WHERE bids.item_id = items.id) as bid_count`).
		Joins("LEFT JOIN auctions ON auctions.id = items.auction_id").
		Where("items.id = ?", id)

	if err := query.Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Check if record was found (uuid.Nil means not found)
	if result.ID == uuid.Nil {
		return nil, nil
	}

	// Calculate can_edit: can edit if item has not started
	result.CanEdit = result.StartedAt == nil

	// Calculate can_delete: can delete if not assigned to auction and has no bids
	result.CanDelete = result.AuctionID == nil && result.BidCount == 0

	return &result, nil
}

// CreateItem creates a new item without auction assignment
func (r *ItemRepository) CreateItem(item *domain.Item) error {
	return r.db.Create(item).Error
}

// UpdateItem updates an item's basic information
func (r *ItemRepository) UpdateItem(itemID string, req *domain.StandaloneItemRequest) (*domain.Item, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	// Get the existing item first
	var item domain.Item
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Update fields
	item.Name = req.Name
	item.Description = req.Description
	item.StartingPrice = req.StartingPrice

	// Save the updated item
	if err := r.db.Save(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem deletes an item by ID
func (r *ItemRepository) DeleteItem(itemID string) error {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return err
	}

	result := r.db.Delete(&domain.Item{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetBidCountByItemID returns the number of bids for an item
func (r *ItemRepository) GetBidCountByItemID(itemID string) (int64, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return 0, err
	}

	var count int64
	result := r.db.Model(&domain.Bid{}).Where("item_id = ?", id).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// FindItemByID finds an item by ID
func (r *ItemRepository) FindItemByID(itemID string) (*domain.Item, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, err
	}

	var item domain.Item
	result := r.db.First(&item, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &item, nil
}

// FindUnassignedItems retrieves items that are not assigned to any auction
func (r *ItemRepository) FindUnassignedItems(search string, limit int) ([]domain.Item, error) {
	var items []domain.Item

	query := r.db.Model(&domain.Item{}).
		Where("auction_id IS NULL")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	query = query.Order("created_at DESC").Limit(limit)

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// AssignItemsToAuction assigns multiple items to an auction with auto-incrementing lot numbers
func (r *ItemRepository) AssignItemsToAuction(auctionID string, itemIDs []uuid.UUID) error {
	aucID, err := uuid.Parse(auctionID)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current max lot_number for the auction
		var maxLotNumber int
		tx.Model(&domain.Item{}).
			Where("auction_id = ?", aucID).
			Select("COALESCE(MAX(lot_number), 0)").
			Scan(&maxLotNumber)

		// Assign each item with incrementing lot_number
		for i, itemID := range itemIDs {
			result := tx.Model(&domain.Item{}).
				Where("id = ? AND auction_id IS NULL", itemID).
				Updates(map[string]interface{}{
					"auction_id": aucID,
					"lot_number": maxLotNumber + i + 1,
				})

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("item not found or already assigned")
			}
		}

		return nil
	})
}

// UnassignItemFromAuction removes an item from an auction
func (r *ItemRepository) UnassignItemFromAuction(auctionID string, itemID string) error {
	aucID, err := uuid.Parse(auctionID)
	if err != nil {
		return err
	}

	itID, err := uuid.Parse(itemID)
	if err != nil {
		return err
	}

	result := r.db.Model(&domain.Item{}).
		Where("id = ? AND auction_id = ?", itID, aucID).
		Updates(map[string]interface{}{
			"auction_id": nil,
			"lot_number": 0,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// HasAnyItemStarted checks if any item in an auction has started
func (r *ItemRepository) HasAnyItemStarted(auctionID string) (bool, error) {
	aucID, err := uuid.Parse(auctionID)
	if err != nil {
		return false, err
	}

	var count int64
	result := r.db.Model(&domain.Item{}).
		Where("auction_id = ? AND started_at IS NOT NULL", aucID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
