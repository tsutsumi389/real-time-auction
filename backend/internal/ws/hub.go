package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

// Hub はWebSocket接続を管理する
type Hub struct {
	// クライアント管理
	clients    map[*Client]bool       // 登録されているクライアント
	rooms      map[string][]*Client   // オークションID -> クライアントリスト
	roomsMutex sync.RWMutex           // ルームマップのロック

	// チャネル
	register     chan *Client        // クライアント登録
	unregister   chan *Client        // クライアント登録解除
	broadcast    chan *BroadcastMsg  // ブロードキャストメッセージ
	handleEvent  chan *ClientEvent   // クライアントイベント

	// Redis
	redisClient  *redis.Client
	ctx          context.Context

	// Repository
	auctionRepo  *repository.AuctionRepository

	// イベントハンドラー
	eventHandler *EventHandler
}

// BroadcastMsg はブロードキャストメッセージを表す
type BroadcastMsg struct {
	auctionID string  // 空文字列の場合は全クライアントに送信
	event     *Event
}

// ClientEvent はクライアントからのイベントを表す
type ClientEvent struct {
	client *Client
	event  *Event
}

// NewHub は新しいHubを作成する
func NewHub(redisClient *redis.Client, auctionRepo *repository.AuctionRepository) *Hub {
	hub := &Hub{
		clients:      make(map[*Client]bool),
		rooms:        make(map[string][]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan *BroadcastMsg, 256),
		handleEvent:  make(chan *ClientEvent, 256),
		redisClient:  redisClient,
		ctx:          context.Background(),
		auctionRepo:  auctionRepo,
	}

	// イベントハンドラーを初期化
	hub.eventHandler = NewEventHandler(hub)

	return hub
}

// Run はHubのメインループを開始する
func (h *Hub) Run() {
	// Redis Pub/Subリスナーを開始
	go h.listenRedis()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case msg := <-h.broadcast:
			h.broadcastMessage(msg)

		case clientEvent := <-h.handleEvent:
			h.eventHandler.Handle(clientEvent.client, clientEvent.event)
		}
	}
}

// registerClient はクライアントを登録する
func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
	log.Printf("Client registered: userID=%s, role=%s", client.userID, client.userRole)
}

// unregisterClient はクライアントの登録を解除する
func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// ルームから削除
		h.roomsMutex.Lock()
		for auctionID := range client.auctionIDs {
			h.removeClientFromRoom(auctionID, client)
		}
		h.roomsMutex.Unlock()

		log.Printf("Client unregistered: userID=%s, role=%s", client.userID, client.userRole)
	}
}

// broadcastMessage はメッセージをブロードキャストする
func (h *Hub) broadcastMessage(msg *BroadcastMsg) {
	message, err := json.Marshal(msg.event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	if msg.auctionID == "" {
		// 全クライアントにブロードキャスト
		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	} else {
		// 特定のオークションルームにブロードキャスト
		h.roomsMutex.RLock()
		clients := h.rooms[msg.auctionID]
		h.roomsMutex.RUnlock()

		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// AddClientToRoom はクライアントをオークションルームに追加する
func (h *Hub) AddClientToRoom(auctionID string, client *Client) {
	h.roomsMutex.Lock()
	defer h.roomsMutex.Unlock()

	if _, ok := h.rooms[auctionID]; !ok {
		h.rooms[auctionID] = make([]*Client, 0)
	}

	// 既に存在する場合は追加しない
	for _, c := range h.rooms[auctionID] {
		if c == client {
			return
		}
	}

	h.rooms[auctionID] = append(h.rooms[auctionID], client)
	client.subscribe(auctionID)

	log.Printf("Client added to room: userID=%s, auctionID=%s", client.userID, auctionID)

	// bidderの場合のみ参加イベントを送信
	if client.userRole == "bidder" && client.bidderID != nil {
		h.broadcastParticipantJoined(auctionID, client)
	}
}

// broadcastParticipantJoined は参加者参加イベントをブロードキャストする
func (h *Hub) broadcastParticipantJoined(auctionID string, client *Client) {
	// データベースから入札者情報を取得
	bidderUUID, err := uuid.Parse(*client.bidderID)
	if err != nil {
		log.Printf("Failed to parse bidder ID: %v", err)
		return
	}

	participantInfo, err := h.auctionRepo.GetBidderInfo(bidderUUID, auctionID)
	if err != nil {
		log.Printf("Failed to get bidder info: %v", err)
		return
	}

	// イベントデータを作成
	participantData := ParticipantData{
		BidderID:    participantInfo.BidderID.String(),
		DisplayName: participantInfo.DisplayName,
		IsOnline:    true,
		BidCount:    participantInfo.BidCount,
		LastBidAt:   participantInfo.LastBidAt,
	}

	event := NewEvent(EventParticipantJoined, auctionID, ParticipantJoinedData{
		AuctionID:   auctionID,
		Participant: participantData,
	})

	// オークションルームにブロードキャスト
	h.BroadcastToAuction(auctionID, event)
}

// RemoveClientFromRoom はクライアントをオークションルームから削除する
func (h *Hub) RemoveClientFromRoom(auctionID string, client *Client) {
	h.roomsMutex.Lock()
	defer h.roomsMutex.Unlock()

	h.removeClientFromRoom(auctionID, client)
}

// removeClientFromRoom はクライアントをオークションルームから削除する（ロックなし）
func (h *Hub) removeClientFromRoom(auctionID string, client *Client) {
	clients, ok := h.rooms[auctionID]
	if !ok {
		return
	}

	for i, c := range clients {
		if c == client {
			h.rooms[auctionID] = append(clients[:i], clients[i+1:]...)
			client.unsubscribe(auctionID)
			log.Printf("Client removed from room: userID=%s, auctionID=%s", client.userID, auctionID)

			// bidderの場合のみ退出イベントを送信
			if c.userRole == "bidder" && c.bidderID != nil {
				h.broadcastParticipantLeft(auctionID, c)
			}
			break
		}
	}

	// ルームが空になった場合は削除
	if len(h.rooms[auctionID]) == 0 {
		delete(h.rooms, auctionID)
	}
}

// broadcastParticipantLeft は参加者退出イベントをブロードキャストする
func (h *Hub) broadcastParticipantLeft(auctionID string, client *Client) {
	event := NewEvent(EventParticipantLeft, auctionID, ParticipantLeftData{
		AuctionID: auctionID,
		BidderID:  *client.bidderID,
	})

	// オークションルームにブロードキャスト
	h.BroadcastToAuction(auctionID, event)
}

// BroadcastToAuction はオークションルームにイベントをブロードキャストする
func (h *Hub) BroadcastToAuction(auctionID string, event *Event) {
	h.broadcast <- &BroadcastMsg{
		auctionID: auctionID,
		event:     event,
	}
}

// BroadcastToAll は全クライアントにイベントをブロードキャストする
func (h *Hub) BroadcastToAll(event *Event) {
	h.broadcast <- &BroadcastMsg{
		auctionID: "",
		event:     event,
	}
}

// PublishEvent はRedis Pub/Subにイベントを発行する
func (h *Hub) PublishEvent(channel string, event *Event) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return h.redisClient.Publish(h.ctx, channel, message).Err()
}

// listenRedis はRedis Pub/Subからイベントを受信する
func (h *Hub) listenRedis() {
	pubsub := h.redisClient.Subscribe(h.ctx,
		"auction:started",
		"auction:price_open",
		"auction:bid",
		"auction:ended",
		"auction:cancelled",
		"auction:item_started",
		"auction:item_ended",
	)
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("Redis Pub/Sub listener started")

	for msg := range ch {
		// Redisからのメッセージを汎用的なマップとして解析
		var rawEvent map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &rawEvent); err != nil {
			log.Printf("Failed to unmarshal Redis message: %v", err)
			continue
		}

		// typeフィールドを取得
		eventType, _ := rawEvent["type"].(string)

		// WebSocketクライアントが期待する形式に変換: { type, data }
		// dataにはtype以外の全フィールドを含める
		delete(rawEvent, "type")

		wsMessage := map[string]interface{}{
			"type": eventType,
			"data": rawEvent,
		}

		// メッセージをJSONに変換
		messageBytes, err := json.Marshal(wsMessage)
		if err != nil {
			log.Printf("Failed to marshal WebSocket message: %v", err)
			continue
		}

		// 全クライアントにブロードキャスト
		for client := range h.clients {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}

		log.Printf("Broadcasted event from Redis: type=%s, channel=%s", eventType, msg.Channel)
	}
}

// GetRoomSize はオークションルームのクライアント数を返す
func (h *Hub) GetRoomSize(auctionID string) int {
	h.roomsMutex.RLock()
	defer h.roomsMutex.RUnlock()

	return len(h.rooms[auctionID])
}

// GetActiveParticipants はオークションルームのアクティブ参加者一覧を返す
func (h *Hub) GetActiveParticipants(auctionID string) ([]ParticipantData, error) {
	h.roomsMutex.RLock()
	clients := h.rooms[auctionID]
	h.roomsMutex.RUnlock()

	// bidder_idで重複排除するためのマップ
	uniqueBidders := make(map[string]bool)
	var bidderIDs []uuid.UUID

	// bidderのみフィルタリングし、重複排除
	for _, client := range clients {
		if client.userRole == "bidder" && client.bidderID != nil {
			if !uniqueBidders[*client.bidderID] {
				uniqueBidders[*client.bidderID] = true
				bidderUUID, err := uuid.Parse(*client.bidderID)
				if err != nil {
					log.Printf("Failed to parse bidder ID: %v", err)
					continue
				}
				bidderIDs = append(bidderIDs, bidderUUID)
			}
		}
	}

	// データベースから各入札者の情報を取得
	var participants []ParticipantData
	for _, bidderID := range bidderIDs {
		participantInfo, err := h.auctionRepo.GetBidderInfo(bidderID, auctionID)
		if err != nil {
			log.Printf("Failed to get bidder info: %v", err)
			continue
		}

		participants = append(participants, ParticipantData{
			BidderID:    participantInfo.BidderID.String(),
			DisplayName: participantInfo.DisplayName,
			IsOnline:    true,
			BidCount:    participantInfo.BidCount,
			LastBidAt:   participantInfo.LastBidAt,
		})
	}

	return participants, nil
}
