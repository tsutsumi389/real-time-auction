package service

import "errors"

// Common service errors
var (
	// Email already exists error (shared across admin and bidder services)
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// Auction service errors
var (
	ErrAuctionNotFound           = errors.New("auction not found")
	ErrNoItemsInAuction          = errors.New("no items found in auction")
	ErrItemsMissingStartingPrice = errors.New("items missing starting price")
	ErrInvalidAuctionStatus      = errors.New("invalid auction status")
	ErrAuctionNotPending         = errors.New("auction is not in pending status")
	ErrAuctionNotActive          = errors.New("auction is not in active status")
	ErrInvalidPage               = errors.New("invalid page number")
	ErrInvalidLimit              = errors.New("invalid limit")
	ErrInvalidSortMode           = errors.New("invalid sort mode")
	ErrInvalidStatus             = errors.New("invalid status")
	ErrAuctionNotEditable        = errors.New("auction cannot be edited")
)

// Item service errors
var (
	ErrItemNotFound           = errors.New("item not found")
	ErrItemAlreadyStarted     = errors.New("item already started")
	ErrItemNotStarted         = errors.New("item not started")
	ErrItemAlreadyEnded       = errors.New("item already ended")
	ErrStartingPriceNotSet    = errors.New("starting price not set")
	ErrPriceTooLow            = errors.New("new price must be higher than current price")
	ErrNoBidsFound            = errors.New("no bids found for this item")
	ErrItemNotEditable        = errors.New("item cannot be edited")
	ErrItemNotDeletable       = errors.New("item cannot be deleted")
	ErrItemHasBids            = errors.New("item has bids and cannot be deleted")
	ErrInvalidItemIDs         = errors.New("invalid item IDs")
	ErrItemAssignedToAuction  = errors.New("item is assigned to an auction")
	ErrItemAlreadyAssigned    = errors.New("item is already assigned to an auction")
	ErrItemNotAssigned        = errors.New("item is not assigned to any auction")
	ErrItemNotInAuction       = errors.New("item is not in this auction")
	ErrAuctionAlreadyStarted  = errors.New("auction has already started")
)
