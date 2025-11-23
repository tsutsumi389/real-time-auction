package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/handler"
	"github.com/tsutsumi389/real-time-auction/internal/middleware"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

func main() {
	// 環境変数取得
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	// Ginモード設定
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// データベース接続
	db, err := repository.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)
	bidderRepo := repository.NewBidderRepository(db)
	auctionRepo := repository.NewAuctionRepository(db)

	// サービス初期化
	jwtService := service.NewJWTService(jwtSecret)
	authService := service.NewAuthService(adminRepo, jwtService)
	adminService := service.NewAdminService(adminRepo)
	bidderService := service.NewBidderService(bidderRepo)
	auctionService := service.NewAuctionService(auctionRepo)

	// ハンドラ初期化
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)
	bidderHandler := handler.NewBidderHandler(bidderService)
	auctionHandler := handler.NewAuctionHandler(auctionService)

	// Ginルーター初期化
	router := gin.Default()

	// CORS ミドルウェア設定
	router.Use(middleware.CORSMiddleware())

	// ヘルスチェックエンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "REST API Server",
			"version": "0.1.0",
		})
	})

	// APIルートグループ
	api := router.Group("/api")
	{
		// 認証エンドポイント
		auth := api.Group("/auth")
		{
			auth.POST("/admin/login", authHandler.AdminLogin)
		}

		// 保護されたエンドポイント（認証が必要）
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			// 現在のユーザー情報取得
			protected.GET("/admin/me", adminHandler.GetCurrentAdmin)

			// システム管理者専用エンドポイント
			systemAdmin := protected.Group("")
			systemAdmin.Use(middleware.RequireSystemAdmin())
			{
				// 管理者登録
				systemAdmin.POST("/admin/admins", adminHandler.RegisterAdmin)
				// 管理者一覧取得
				systemAdmin.GET("/admin/admins", adminHandler.GetAdminList)
				// 管理者状態変更
				systemAdmin.PATCH("/admin/admins/:id/status", adminHandler.UpdateAdminStatus)

				// 入札者登録
				systemAdmin.POST("/admin/bidders", bidderHandler.RegisterBidder)
				// 入札者一覧取得
				systemAdmin.GET("/admin/bidders", bidderHandler.GetBidderList)
				// 入札者へのポイント付与
				systemAdmin.POST("/admin/bidders/:id/points", bidderHandler.GrantPoints)
				// 入札者のポイント履歴取得
				systemAdmin.GET("/admin/bidders/:id/points/history", bidderHandler.GetPointHistory)
				// 入札者状態変更
				systemAdmin.PATCH("/admin/bidders/:id/status", bidderHandler.UpdateBidderStatus)

				// オークション中止（system_adminのみ）
				systemAdmin.POST("/admin/auctions/:id/cancel", auctionHandler.CancelAuction)
			}

			// システム管理者と主催者がアクセス可能なエンドポイント
			adminOrAuctioneer := protected.Group("")
			adminOrAuctioneer.Use(middleware.RequireAdminOrAuctioneer())
			{
				// オークション一覧取得
				adminOrAuctioneer.GET("/admin/auctions", auctionHandler.GetAuctionList)
				// オークション開始
				adminOrAuctioneer.POST("/admin/auctions/:id/start", auctionHandler.StartAuction)
				// オークション終了
				adminOrAuctioneer.POST("/admin/auctions/:id/end", auctionHandler.EndAuction)
			}
		}
	}

	// サーバー起動
	log.Printf("Starting REST API server on port %s (env: %s)", port, env)
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
