package service

import "github.com/tsutsumi389/real-time-auction/internal/domain"

// AuthServiceInterface defines the interface for authentication service operations
type AuthServiceInterface interface {
	LoginAdmin(email, password string) (*domain.LoginResponse, error)
	LoginBidder(email, password string) (*domain.LoginResponse, error)
}

// JWTServiceInterface defines the interface for JWT service operations
type JWTServiceInterface interface {
	GenerateTokenForAdmin(admin *domain.Admin) (string, error)
	GenerateTokenForBidder(bidder *domain.Bidder) (string, error)
	ValidateToken(tokenString string) (*domain.JWTClaims, error)
}

// AdminServiceInterface defines the interface for admin service operations
type AdminServiceInterface interface {
	RegisterAdmin(req *domain.AdminCreateRequest) (*domain.Admin, error)
	GetAdminByID(id int64) (*domain.Admin, error)
	GetAdminList(req *domain.AdminListRequest) (*domain.AdminListResponse, error)
	UpdateAdminStatus(id int64, status domain.AdminStatus) (*domain.Admin, error)
	UpdateAdmin(id int64, req *domain.AdminUpdateRequest, currentUserID int64) (*domain.Admin, error)
}

// BidderServiceInterface defines the interface for bidder service operations
type BidderServiceInterface interface {
	RegisterBidder(req *domain.BidderCreateRequest, adminID int64) (*domain.BidderResponse, error)
	GetBidderByID(id string) (*domain.Bidder, error)
	GetBidderDetail(id string) (*domain.BidderDetailResponse, error)
	GetBidderList(req *domain.BidderListRequest) (*domain.BidderListResponse, error)
	GrantPoints(bidderID string, points int64, adminID int64) (*domain.GrantPointsResponse, error)
	GetPointHistory(bidderID string, page int, limit int) (*domain.PointHistoryListResponse, error)
	UpdateBidderStatus(id string, status domain.BidderStatus) (*domain.Bidder, error)
	UpdateBidder(id string, req *domain.BidderUpdateRequest) (*domain.BidderDetailResponse, error)
}

// AuctionServiceInterface defines the interface for auction service operations
type AuctionServiceInterface interface {
	// Auction-level operations
	GetAuctionList(req *domain.AuctionListRequest) (*domain.AuctionListResponse, error)
	GetBidderAuctionList(req *domain.BidderAuctionListRequest) (*domain.BidderAuctionListResponse, error)
	GetAuctionDetail(id string) (*domain.GetAuctionDetailResponse, error)
	StartAuction(id string) (*domain.AuctionWithItemCount, error)
	EndAuction(id string) (*domain.AuctionWithItemCount, error)
	CancelAuction(id string) (*domain.AuctionWithItemCount, error)
	CancelAuctionWithReason(id string, reason string) (*domain.CancelAuctionResponse, error)
	CreateAuction(req *domain.CreateAuctionRequest) (*domain.CreateAuctionResponse, error)

	// Item-level operations
	StartItem(itemID string) (*domain.StartItemResponse, error)
	OpenPrice(itemID string, newPrice int64, adminID int64) (*domain.OpenPriceResponse, error)
	EndItem(itemID string) (*domain.EndItemResponse, error)

	// Query operations
	GetBidHistory(itemID string, limit int, offset int) (*domain.BidHistoryResponse, error)
	GetPriceHistory(itemID string) (*domain.PriceHistoryResponse, error)
	GetParticipants(auctionID string) (*domain.ParticipantsResponse, error)

	// Edit operations
	GetAuctionForEdit(id string) (*domain.AuctionEditResponse, error)
	UpdateAuction(id string, req *domain.UpdateAuctionRequest) (*domain.Auction, error)
	UpdateItem(itemID string, req *domain.UpdateItemRequest) (*domain.Item, error)
	DeleteItem(itemID string) error
	AddItem(auctionID string, req *domain.AddItemRequest) (*domain.Item, error)
	ReorderItems(auctionID string, req *domain.ReorderItemsRequest) error
}
