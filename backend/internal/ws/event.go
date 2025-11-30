package ws

import "time"

// EventType はWebSocketイベントのタイプを表す
type EventType string

const (
	// オークションイベント
	EventAuctionStarted   EventType = "auction:started"
	EventAuctionPriceOpen EventType = "auction:price_open"
	EventAuctionBid       EventType = "auction:bid"
	EventAuctionEnded     EventType = "auction:ended"
	EventAuctionCancelled EventType = "auction:cancelled"

	// 参加者イベント
	EventParticipantJoined EventType = "participant:joined"
	EventParticipantLeft   EventType = "participant:left"
	EventParticipantsList  EventType = "participants:list"

	// システムイベント
	EventError       EventType = "error"
	EventPing        EventType = "ping"
	EventPong        EventType = "pong"
	EventSubscribe   EventType = "subscribe"
	EventUnsubscribe EventType = "unsubscribe"
)

// Event はWebSocketイベントの基本構造
type Event struct {
	Type      EventType   `json:"type"`
	AuctionID string      `json:"auction_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// AuctionStartedData はオークション開始イベントのデータ
type AuctionStartedData struct {
	AuctionID     string `json:"auction_id"`
	Title         string `json:"title"`
	ItemName      string `json:"item_name"`
	StartingPrice int    `json:"starting_price"`
	CurrentPrice  int    `json:"current_price"`
}

// PriceOpenData は価格開示イベントのデータ
type PriceOpenData struct {
	AuctionID     string    `json:"auction_id"`
	NewPrice      int       `json:"new_price"`
	PreviousPrice int       `json:"previous_price,omitempty"`
	OpenedAt      time.Time `json:"opened_at"`
}

// BidData は入札イベントのデータ
type BidData struct {
	BidID      int64     `json:"bid_id"`
	AuctionID  string    `json:"auction_id"`
	BidderName string    `json:"bidder_name"` // display_name or anonymized
	Price      int       `json:"price"`
	BidAt      time.Time `json:"bid_at"`
	IsWinning  bool      `json:"is_winning"`
}

// AuctionEndedData はオークション終了イベントのデータ
type AuctionEndedData struct {
	AuctionID  string     `json:"auction_id"`
	WinnerID   *string    `json:"winner_id,omitempty"` // UUID string
	WinnerName *string    `json:"winner_name,omitempty"`
	FinalPrice *int       `json:"final_price,omitempty"`
	EndedAt    time.Time  `json:"ended_at"`
	Reason     string     `json:"reason"` // "sold", "no_bids", "cancelled"
}

// AuctionCancelledData はオークション中止イベントのデータ
type AuctionCancelledData struct {
	AuctionID   string    `json:"auction_id"`
	CancelledAt time.Time `json:"cancelled_at"`
	Reason      string    `json:"reason"`
}

// ErrorData はエラーイベントのデータ
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SubscribeData はサブスクライブリクエストのデータ
type SubscribeData struct {
	AuctionID string `json:"auction_id"`
}

// ParticipantData は参加者情報のデータ
type ParticipantData struct {
	BidderID    string     `json:"bidder_id"`
	DisplayName string     `json:"display_name"`
	IsOnline    bool       `json:"is_online"`
	BidCount    int64      `json:"bid_count"`
	LastBidAt   *time.Time `json:"last_bid_at,omitempty"`
}

// ParticipantJoinedData は参加者参加イベントのデータ
type ParticipantJoinedData struct {
	AuctionID   string          `json:"auction_id"`
	Participant ParticipantData `json:"participant"`
}

// ParticipantLeftData は参加者退出イベントのデータ
type ParticipantLeftData struct {
	AuctionID string `json:"auction_id"`
	BidderID  string `json:"bidder_id"`
}

// ParticipantsListData は参加者一覧イベントのデータ
type ParticipantsListData struct {
	AuctionID    string            `json:"auction_id"`
	Participants []ParticipantData `json:"participants"`
}

// NewEvent は新しいイベントを作成する
func NewEvent(eventType EventType, auctionID string, data interface{}) *Event {
	return &Event{
		Type:      eventType,
		AuctionID: auctionID,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorEvent はエラーイベントを作成する
func NewErrorEvent(code, message string) *Event {
	return &Event{
		Type: EventError,
		Data: ErrorData{
			Code:    code,
			Message: message,
		},
		Timestamp: time.Now(),
	}
}
