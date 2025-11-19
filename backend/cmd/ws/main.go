package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 環境変数取得
	port := getEnv("PORT", "8081")
	env := getEnv("ENV", "development")

	// Ginモード設定
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Ginルーター初期化
	router := gin.Default()

	// ヘルスチェックエンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "WebSocket Server",
			"version": "0.1.0",
		})
	})

	// WebSocketエンドポイント（準備中）
	router.GET("/ws", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "WebSocket endpoint (not yet implemented)",
		})
	})

	// サーバー起動
	log.Printf("Starting WebSocket server on port %s (env: %s)", port, env)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
