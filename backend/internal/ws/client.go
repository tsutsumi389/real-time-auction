package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Ping/Pong設定
	writeWait      = 10 * time.Second    // WebSocketへの書き込みタイムアウト
	pongWait       = 60 * time.Second    // Pong待機タイムアウト
	pingPeriod     = 30 * time.Second    // Ping送信間隔 (pongWaitより短く設定)
	maxMessageSize = 512 * 1024          // 最大メッセージサイズ (512KB)
)

// Client はWebSocket接続を表す
type Client struct {
	hub         *Hub            // Hubへの参照
	conn        *websocket.Conn // WebSocket接続
	send        chan []byte     // 送信チャネル
	userID      string          // ユーザーID (bidder UUID or admin ID)
	userRole    string          // ユーザーロール (bidder, auctioneer, system_admin)
	bidderID    *string         // 入札者ID (bidderの場合のみ、UUID文字列)
	displayName string          // 表示名
	auctionIDs  map[int64]bool  // 購読中のオークションID
}

// NewClient は新しいクライアントを作成する
func NewClient(hub *Hub, conn *websocket.Conn, userID, userRole string, bidderID *string, displayName string) *Client {
	return &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		userID:      userID,
		userRole:    userRole,
		bidderID:    bidderID,
		displayName: displayName,
		auctionIDs:  make(map[int64]bool),
	}
}

// readPump はWebSocketからのメッセージを読み取る
// クライアントごとに1つのgoroutineで実行される
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// イベントをパース
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			c.sendError("INVALID_EVENT", "Invalid event format")
			continue
		}

		// イベントをハンドラーに渡す
		c.hub.handleEvent <- &ClientEvent{
			client: c,
			event:  &event,
		}
	}
}

// writePump はクライアントへメッセージを送信する
// クライアントごとに1つのgoroutineで実行される
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hubがチャネルを閉じた
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// キューにある追加メッセージも送信
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sendEvent はクライアントにイベントを送信する
func (c *Client) sendEvent(event *Event) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	select {
	case c.send <- message:
	default:
		// 送信バッファがいっぱいの場合、クライアントを切断
		close(c.send)
	}

	return nil
}

// sendError はクライアントにエラーイベントを送信する
func (c *Client) sendError(code, message string) {
	event := NewErrorEvent(code, message)
	c.sendEvent(event)
}

// subscribe はオークションルームに参加する
func (c *Client) subscribe(auctionID int64) {
	c.auctionIDs[auctionID] = true
}

// unsubscribe はオークションルームから退出する
func (c *Client) unsubscribe(auctionID int64) {
	delete(c.auctionIDs, auctionID)
}

// isSubscribed はオークションルームに参加しているかチェック
func (c *Client) isSubscribed(auctionID int64) bool {
	return c.auctionIDs[auctionID]
}
