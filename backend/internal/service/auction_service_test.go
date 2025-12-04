package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
)

// MockAuctionRepository is a mock implementation of AuctionRepository
type MockAuctionRepository struct {
	mock.Mock
}

func (m *MockAuctionRepository) CreateAuctionWithItems(auction *domain.Auction, items []domain.Item) error {
	args := m.Called(auction, items)
	return args.Error(0)
}

func (m *MockAuctionRepository) FindByID(id string) (*domain.Auction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Auction), args.Error(1)
}

func (m *MockAuctionRepository) FindAuctionWithItems(id string) (*domain.GetAuctionDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GetAuctionDetailResponse), args.Error(1)
}

func (m *MockAuctionRepository) FindAuctionsWithFilters(req *domain.AuctionListRequest) ([]domain.AuctionWithItemCount, error) {
	args := m.Called(req)
	return args.Get(0).([]domain.AuctionWithItemCount), args.Error(1)
}

func (m *MockAuctionRepository) CountAuctionsWithFilters(req *domain.AuctionListRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuctionRepository) UpdateAuctionStatus(id string, status domain.AuctionStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockAuctionRepository) CountItemsByAuctionID(auctionID string) (int64, error) {
	args := m.Called(auctionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuctionRepository) FindItemsByAuctionID(auctionID string) ([]domain.Item, error) {
	args := m.Called(auctionID)
	return args.Get(0).([]domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) CreateAuction(auction *domain.Auction) error {
	args := m.Called(auction)
	return args.Error(0)
}

func (m *MockAuctionRepository) CreateItems(items []domain.Item) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *MockAuctionRepository) FindPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) ([]domain.BidderAuctionSummary, error) {
	args := m.Called(req)
	return args.Get(0).([]domain.BidderAuctionSummary), args.Error(1)
}

func (m *MockAuctionRepository) CountPublicAuctionsWithFilters(req *domain.BidderAuctionListRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuctionRepository) FindItemByID(itemID string) (*domain.Item, error) {
	args := m.Called(itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) StartItem(itemID string) (*domain.Item, error) {
	args := m.Called(itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) UpdateItemCurrentPrice(itemID string, price int64) error {
	args := m.Called(itemID, price)
	return args.Error(0)
}

func (m *MockAuctionRepository) EndItem(itemID string, winnerID uuid.UUID, finalPrice int64) (*domain.Item, error) {
	args := m.Called(itemID, winnerID, finalPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) FindBidsByItemID(itemID string, limit int, offset int) ([]domain.BidWithBidderInfo, error) {
	args := m.Called(itemID, limit, offset)
	return args.Get(0).([]domain.BidWithBidderInfo), args.Error(1)
}

func (m *MockAuctionRepository) CountBidsByItemID(itemID string) (int64, error) {
	args := m.Called(itemID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuctionRepository) FindWinningBidByItemID(itemID string) (*domain.Bid, error) {
	args := m.Called(itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bid), args.Error(1)
}

func (m *MockAuctionRepository) CreateBid(bid *domain.Bid) error {
	args := m.Called(bid)
	return args.Error(0)
}

func (m *MockAuctionRepository) UpdateBidWinningStatus(itemID string, winningBidID int64) error {
	args := m.Called(itemID, winningBidID)
	return args.Error(0)
}

func (m *MockAuctionRepository) FindPriceHistoryByItemID(itemID string) ([]domain.PriceHistoryWithAdmin, error) {
	args := m.Called(itemID)
	return args.Get(0).([]domain.PriceHistoryWithAdmin), args.Error(1)
}

func (m *MockAuctionRepository) CreatePriceHistory(history *domain.PriceHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockAuctionRepository) FindParticipantsByAuctionID(auctionID string) ([]domain.ParticipantInfo, error) {
	args := m.Called(auctionID)
	return args.Get(0).([]domain.ParticipantInfo), args.Error(1)
}

func (m *MockAuctionRepository) GetBidderInfo(bidderID uuid.UUID, auctionIDStr string) (*domain.ParticipantInfo, error) {
	args := m.Called(bidderID, auctionIDStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ParticipantInfo), args.Error(1)
}

func (m *MockAuctionRepository) CancelAuctionWithRefunds(auctionID string, reason string) (*domain.CancelAuctionResponse, error) {
	args := m.Called(auctionID, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CancelAuctionResponse), args.Error(1)
}

func (m *MockAuctionRepository) UpdateAuction(id string, req *domain.UpdateAuctionRequest) (*domain.Auction, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Auction), args.Error(1)
}

func (m *MockAuctionRepository) GetAuctionForEdit(id string) (*domain.AuctionEditResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuctionEditResponse), args.Error(1)
}

func (m *MockAuctionRepository) UpdateItem(itemID string, req *domain.UpdateItemRequest) (*domain.Item, error) {
	args := m.Called(itemID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) DeleteItem(itemID string) error {
	args := m.Called(itemID)
	return args.Error(0)
}

func (m *MockAuctionRepository) AddItem(auctionID string, req *domain.AddItemRequest) (*domain.Item, error) {
	args := m.Called(auctionID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockAuctionRepository) ReorderItems(auctionID string, itemIDs []uuid.UUID) error {
	args := m.Called(auctionID, itemIDs)
	return args.Error(0)
}

func (m *MockAuctionRepository) GetMaxLotNumber(auctionID string) (int, error) {
	args := m.Called(auctionID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockAuctionRepository) FindItemsForEdit(auctionID string) ([]domain.ItemEditInfo, error) {
	args := m.Called(auctionID)
	return args.Get(0).([]domain.ItemEditInfo), args.Error(1)
}

// TestCreateAuction_WithZeroItems tests creating an auction with no items
func TestCreateAuction_WithZeroItems(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(nil, mockRepo, nil, nil, nil)

	startedAt := time.Now().Add(24 * time.Hour)
	req := &domain.CreateAuctionRequest{
		Title:       "Test Auction",
		Description: "Test Description",
		StartedAt:   &startedAt,
		Items:       []domain.CreateItemRequest{}, // Empty items array
	}

	// Mock the repository call
	mockRepo.On("CreateAuctionWithItems", mock.AnythingOfType("*domain.Auction"), mock.AnythingOfType("[]domain.Item")).
		Run(func(args mock.Arguments) {
			auction := args.Get(0).(*domain.Auction)
			// Simulate DB setting the ID and timestamps
			auction.ID = uuid.New()
			auction.CreatedAt = time.Now()
			auction.UpdatedAt = time.Now()
		}).
		Return(nil)

	// Act
	result, err := service.CreateAuction(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Auction", result.Title)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, domain.AuctionStatusPending, result.Status)
	assert.Equal(t, 0, result.ItemCount) // Should be 0 items
	mockRepo.AssertExpectations(t)
}

// TestCreateAuction_WithOneItem tests creating an auction with one item
func TestCreateAuction_WithOneItem(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(nil, mockRepo, nil, nil, nil)

	startedAt := time.Now().Add(24 * time.Hour)
	startingPrice := int64(1000)
	req := &domain.CreateAuctionRequest{
		Title:       "Test Auction",
		Description: "Test Description",
		StartedAt:   &startedAt,
		Items: []domain.CreateItemRequest{
			{
				Name:          "Item 1",
				Description:   "Item 1 Description",
				LotNumber:     1,
				StartingPrice: &startingPrice,
			},
		},
	}

	// Mock the repository call
	mockRepo.On("CreateAuctionWithItems", mock.AnythingOfType("*domain.Auction"), mock.AnythingOfType("[]domain.Item")).
		Run(func(args mock.Arguments) {
			auction := args.Get(0).(*domain.Auction)
			auction.ID = uuid.New()
			auction.CreatedAt = time.Now()
			auction.UpdatedAt = time.Now()
		}).
		Return(nil)

	// Act
	result, err := service.CreateAuction(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Auction", result.Title)
	assert.Equal(t, 1, result.ItemCount) // Should be 1 item
	mockRepo.AssertExpectations(t)
}

// TestCreateAuction_WithMultipleItems tests creating an auction with multiple items
func TestCreateAuction_WithMultipleItems(t *testing.T) {
	// Arrange
	mockRepo := new(MockAuctionRepository)
	service := NewAuctionService(nil, mockRepo, nil, nil, nil)

	startedAt := time.Now().Add(24 * time.Hour)
	startingPrice1 := int64(1000)
	startingPrice2 := int64(2000)
	startingPrice3 := int64(3000)
	req := &domain.CreateAuctionRequest{
		Title:       "Test Auction",
		Description: "Test Description",
		StartedAt:   &startedAt,
		Items: []domain.CreateItemRequest{
			{
				Name:          "Item 1",
				Description:   "Item 1 Description",
				LotNumber:     1,
				StartingPrice: &startingPrice1,
			},
			{
				Name:          "Item 2",
				Description:   "Item 2 Description",
				LotNumber:     2,
				StartingPrice: &startingPrice2,
			},
			{
				Name:          "Item 3",
				Description:   "Item 3 Description",
				LotNumber:     3,
				StartingPrice: &startingPrice3,
			},
		},
	}

	// Mock the repository call
	mockRepo.On("CreateAuctionWithItems", mock.AnythingOfType("*domain.Auction"), mock.AnythingOfType("[]domain.Item")).
		Run(func(args mock.Arguments) {
			auction := args.Get(0).(*domain.Auction)
			auction.ID = uuid.New()
			auction.CreatedAt = time.Now()
			auction.UpdatedAt = time.Now()
		}).
		Return(nil)

	// Act
	result, err := service.CreateAuction(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Auction", result.Title)
	assert.Equal(t, 3, result.ItemCount) // Should be 3 items
	mockRepo.AssertExpectations(t)
}
