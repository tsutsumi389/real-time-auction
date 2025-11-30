package repository

import (
	"github.com/google/uuid"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
)

// AdminRepositoryInterface defines the interface for admin repository operations
type AdminRepositoryInterface interface {
	FindByEmail(email string) (*domain.Admin, error)
	FindByID(id int64) (*domain.Admin, error)
	Create(admin *domain.Admin) error
	Update(admin *domain.Admin) error
	Delete(id int64) error
	FindAdminsWithFilters(req *domain.AdminListRequest) ([]domain.Admin, error)
	CountAdminsWithFilters(req *domain.AdminListRequest) (int64, error)
	UpdateAdminStatus(id int64, status domain.AdminStatus) error
	CountActiveSystemAdmins() (int64, error)
	FindByEmailExcludeID(email string, excludeID int64) (*domain.Admin, error)
}

// BidderRepositoryInterface defines the interface for bidder repository operations
type BidderRepositoryInterface interface {
	FindByID(id string) (*domain.Bidder, error)
	FindByEmail(email string) (*domain.Bidder, error)
	FindByIDWithPoints(id string) (*domain.BidderDetailResponse, error)
	FindBiddersWithFilters(req *domain.BidderListRequest) ([]domain.BidderWithPoints, error)
	CountBiddersWithFilters(req *domain.BidderListRequest) (int64, error)
	GetBidderPoints(bidderID string) (*domain.BidderPoints, error)
	GrantPoints(bidderID string, points int64, adminID int64) (*domain.BidderWithPoints, *domain.PointHistory, error)
	GetPointHistory(bidderID string, page int, limit int) ([]domain.PointHistoryWithAuction, error)
	CountPointHistory(bidderID string) (int64, error)
	UpdateBidderStatus(id string, status domain.BidderStatus) error
	CreateBidderWithPoints(bidder *domain.Bidder, initialPoints int64, adminID int64) (*domain.BidderResponse, error)
	UpdateBidder(id string, req *domain.BidderUpdateRequest, passwordHash *string) error
}

// AuctionRepositoryInterface defines the interface for auction repository operations
type AuctionRepositoryInterface interface {
	FindByID(id string) (*domain.Auction, error)
	FindAuctionWithItems(id string) (*domain.GetAuctionDetailResponse, error)
	FindAuctionsWithFilters(req *domain.AuctionListRequest) ([]domain.AuctionWithItemCount, error)
	CountAuctionsWithFilters(req *domain.AuctionListRequest) (int64, error)
	FindPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) ([]domain.BidderAuctionSummary, error)
	CountPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) (int64, error)
	UpdateAuctionStatus(id string, status domain.AuctionStatus) error
	CountItemsByAuctionID(auctionID string) (int64, error)
	FindItemsByAuctionID(auctionID string) ([]domain.Item, error)
	CreateAuction(auction *domain.Auction) error
	CreateItems(items []domain.Item) error
	CreateAuctionWithItems(auction *domain.Auction, items []domain.Item) error

	// Item operations
	FindItemByID(itemID string) (*domain.Item, error)
	StartItem(itemID string) (*domain.Item, error)
	UpdateItemCurrentPrice(itemID string, price int64) error
	EndItem(itemID string, winnerID uuid.UUID, finalPrice int64) (*domain.Item, error)

	// Bid operations
	FindBidsByItemID(itemID string, limit int, offset int) ([]domain.BidWithBidderInfo, error)
	CountBidsByItemID(itemID string) (int64, error)
	FindWinningBidByItemID(itemID string) (*domain.Bid, error)
	CreateBid(bid *domain.Bid) error
	UpdateBidWinningStatus(itemID string, winningBidID int64) error

	// Price history operations
	FindPriceHistoryByItemID(itemID string) ([]domain.PriceHistoryWithAdmin, error)
	CreatePriceHistory(history *domain.PriceHistory) error

	// Participant operations
	FindParticipantsByAuctionID(auctionID string) ([]domain.ParticipantInfo, error)

	// Cancel auction operations
	CancelAuctionWithRefunds(auctionID string, reason string) (*domain.CancelAuctionResponse, error)

	// Edit operations
	UpdateAuction(id string, req *domain.UpdateAuctionRequest) (*domain.Auction, error)
	GetAuctionForEdit(id string) (*domain.AuctionEditResponse, error)
	UpdateItem(itemID string, req *domain.UpdateItemRequest) (*domain.Item, error)
	DeleteItem(itemID string) error
	AddItem(auctionID string, req *domain.AddItemRequest) (*domain.Item, error)
	ReorderItems(auctionID string, itemIDs []uuid.UUID) error
}
