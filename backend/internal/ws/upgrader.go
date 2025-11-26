package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 本番環境では適切なオリジンチェックを実装
		return true
	},
}

// ServeWs はWebSocket接続をアップグレードし、クライアントを登録する
func ServeWs(hub *Hub, c *gin.Context) {
	// TODO: JWTトークンからユーザー情報を取得
	// 現在は仮の実装
	userID := c.Query("user_id")
	userRole := c.Query("role")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id"})
		return
	}

	if userRole == "" {
		userRole = "bidder"
	}

	// WebSocket接続にアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// クライアントを作成
	client := NewClient(hub, conn, userID, userRole)

	// Hubに登録
	hub.register <- client

	// readPumpとwritePumpをgoroutineで開始
	go client.writePump()
	go client.readPump()
}
