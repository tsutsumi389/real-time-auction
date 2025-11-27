package ws

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
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
	// JWTトークンをクエリパラメータから取得
	tokenString := c.Query("token")
	if tokenString == "" {
		// Authorization ヘッダーからも確認
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}

	// JWTトークンを検証
	jwtService := service.NewJWTService("")
	claims, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// ユーザー情報を取得
	var userID string
	var userRole string

	if claims.UserType == domain.UserTypeAdmin {
		// 管理者の場合、int64のIDを文字列に変換
		if id, ok := claims.GetUserIDAsInt64(); ok {
			userID = strconv.FormatInt(id, 10)
		} else {
			log.Printf("Failed to get admin user ID")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		// ロールを文字列に変換
		switch claims.Role {
		case domain.RoleSystemAdmin:
			userRole = "system_admin"
		case domain.RoleAuctioneer:
			userRole = "auctioneer"
		default:
			userRole = "admin"
		}
	} else {
		// 入札者の場合、文字列のUUIDをそのまま使用
		if id, ok := claims.GetUserIDAsString(); ok {
			userID = id
		} else {
			log.Printf("Failed to get bidder user ID")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
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
