package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tsutsumi389/real-time-auction/internal/ws"
)

func main() {
	// 環境変数取得
	port := getEnv("PORT", "8081")
	env := getEnv("ENV", "development")
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// Ginモード設定
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Redis接続
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Redis接続確認
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")

	// Hubを初期化
	hub := ws.NewHub(redisClient)
	go hub.Run()

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

	// WebSocketエンドポイント
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
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
