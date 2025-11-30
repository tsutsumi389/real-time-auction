package ws

import (
	"encoding/json"
	"log"
)

// EventHandler はクライアントからのイベントを処理する
type EventHandler struct {
	hub *Hub
}

// NewEventHandler は新しいEventHandlerを作成する
func NewEventHandler(hub *Hub) *EventHandler {
	return &EventHandler{
		hub: hub,
	}
}

// Handle はイベントを処理する
func (h *EventHandler) Handle(client *Client, event *Event) {
	switch event.Type {
	case EventSubscribe:
		h.handleSubscribe(client, event)
	case EventUnsubscribe:
		h.handleUnsubscribe(client, event)
	case EventPing:
		h.handlePing(client, event)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		client.sendError("UNKNOWN_EVENT", "Unknown event type")
	}
}

// handleSubscribe はオークションルームへの参加リクエストを処理する
func (h *EventHandler) handleSubscribe(client *Client, event *Event) {
	log.Printf("[Subscribe] Received subscribe event from userID=%s, role=%s", client.userID, client.userRole)

	var data SubscribeData
	if err := h.parseEventData(event, &data); err != nil {
		log.Printf("[Subscribe] Failed to parse event data: %v", err)
		client.sendError("INVALID_DATA", "Invalid subscribe data")
		return
	}

	log.Printf("[Subscribe] Parsed auction_id=%s", data.AuctionID)

	if data.AuctionID == "" {
		log.Printf("[Subscribe] Empty auction ID")
		client.sendError("INVALID_AUCTION_ID", "Invalid auction ID")
		return
	}

	// TODO: オークションが存在し、アクティブかチェック
	// TODO: 権限チェック（入札者は自分が参加可能なオークションのみ）

	// クライアントをルームに追加（これにより participant:joined イベントが送信される）
	h.hub.AddClientToRoom(data.AuctionID, client)

	// 現在のアクティブ参加者一覧を取得
	participants, err := h.hub.GetActiveParticipants(data.AuctionID)
	if err != nil {
		log.Printf("Failed to get active participants: %v", err)
		participants = []ParticipantData{} // エラー時は空配列
	}

	// 初期参加者リストを送信
	participantsListEvent := NewEvent(EventParticipantsList, data.AuctionID, ParticipantsListData{
		AuctionID:    data.AuctionID,
		Participants: participants,
	})
	client.sendEvent(participantsListEvent)

	// 確認メッセージを送信
	response := NewEvent("subscribed", data.AuctionID, map[string]interface{}{
		"auction_id": data.AuctionID,
		"message":    "Successfully subscribed to auction",
	})
	client.sendEvent(response)
}

// handleUnsubscribe はオークションルームからの退出リクエストを処理する
func (h *EventHandler) handleUnsubscribe(client *Client, event *Event) {
	var data SubscribeData
	if err := h.parseEventData(event, &data); err != nil {
		client.sendError("INVALID_DATA", "Invalid unsubscribe data")
		return
	}

	if data.AuctionID == "" {
		client.sendError("INVALID_AUCTION_ID", "Invalid auction ID")
		return
	}

	// クライアントをルームから削除
	h.hub.RemoveClientFromRoom(data.AuctionID, client)

	// 確認メッセージを送信
	response := NewEvent("unsubscribed", data.AuctionID, map[string]interface{}{
		"auction_id": data.AuctionID,
		"message":    "Successfully unsubscribed from auction",
	})
	client.sendEvent(response)
}

// handlePing はPingリクエストを処理する
func (h *EventHandler) handlePing(client *Client, event *Event) {
	response := NewEvent(EventPong, "", map[string]interface{}{
		"message": "pong",
	})
	client.sendEvent(response)
}

// parseEventData はイベントデータをパースする
func (h *EventHandler) parseEventData(event *Event, v interface{}) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
