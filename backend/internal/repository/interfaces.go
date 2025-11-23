package repository

import "github.com/tsutsumi389/real-time-auction/internal/domain"

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
}

// BidderRepositoryInterface defines the interface for bidder repository operations
type BidderRepositoryInterface interface {
	FindByID(id string) (*domain.Bidder, error)
	FindByEmail(email string) (*domain.Bidder, error)
	FindBiddersWithFilters(req *domain.BidderListRequest) ([]domain.BidderWithPoints, error)
	CountBiddersWithFilters(req *domain.BidderListRequest) (int64, error)
	GetBidderPoints(bidderID string) (*domain.BidderPoints, error)
	GrantPoints(bidderID string, points int64, adminID int64) (*domain.BidderWithPoints, *domain.PointHistory, error)
	GetPointHistory(bidderID string, page int, limit int) ([]domain.PointHistoryWithAuction, error)
	CountPointHistory(bidderID string) (int64, error)
	UpdateBidderStatus(id string, status domain.BidderStatus) error
	CreateBidderWithPoints(bidder *domain.Bidder, initialPoints int64, adminID int64) (*domain.BidderResponse, error)
}

// AuctionRepositoryInterface defines the interface for auction repository operations
type AuctionRepositoryInterface interface {
	FindByID(id string) (*domain.Auction, error)
	FindAuctionsWithFilters(req *domain.AuctionListRequest) ([]domain.AuctionWithItemCount, error)
	CountAuctionsWithFilters(req *domain.AuctionListRequest) (int64, error)
	UpdateAuctionStatus(id string, status domain.AuctionStatus) error
	CountItemsByAuctionID(auctionID string) (int64, error)
	FindItemsByAuctionID(auctionID string) ([]domain.Item, error)
}
