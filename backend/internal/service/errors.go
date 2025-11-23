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
)
