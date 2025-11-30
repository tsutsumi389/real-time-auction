package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"github.com/tsutsumi389/real-time-auction/internal/ws"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 環境変数取得
	port := getEnv("PORT", "8081")
	env := getEnv("ENV", "development")
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "auction_user")
	dbPassword := getEnv("DB_PASSWORD", "auction_pass")
	dbName := getEnv("DB_NAME", "auction_db")

	// Ginモード設定
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// データベース接続
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to PostgreSQL")

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

	// Repository初期化
	auctionRepo := repository.NewAuctionRepository(db)

	// Hubを初期化
	hub := ws.NewHub(redisClient, auctionRepo)
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
