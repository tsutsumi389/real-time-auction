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

	// Redis接続
	redisClient, err := repository.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)
	bidderRepo := repository.NewBidderRepository(db)
	auctionRepo := repository.NewAuctionRepository(db)
	pointRepo := repository.NewPointRepository(db)
	bidRepo := repository.NewBidRepository(db)
	itemRepo := repository.NewItemRepository(db)

	// サービス初期化
	jwtService := service.NewJWTService(jwtSecret)
	authService := service.NewAuthService(adminRepo, bidderRepo, jwtService)
	adminService := service.NewAdminService(adminRepo)
	bidderService := service.NewBidderService(bidderRepo)
	auctionService := service.NewAuctionService(db, auctionRepo, bidRepo, pointRepo, redisClient)
	pointService := service.NewPointService(pointRepo)
	bidService := service.NewBidService(db, redisClient, bidRepo, pointRepo, auctionRepo)
	itemService := service.NewItemService(itemRepo)

	// ハンドラ初期化
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)
	bidderHandler := handler.NewBidderHandler(bidderService)
	auctionHandler := handler.NewAuctionHandler(auctionService)
	bidHandler := handler.NewBidHandler(pointService, bidService)
	itemHandler := handler.NewItemHandler(itemService)

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
			auth.POST("/bidder/login", authHandler.BidderLogin)
		}

		// 公開エンドポイント（認証不要）
		// 入札者用オークション一覧取得（すべてのユーザーがアクセス可能）
		api.GET("/auctions", auctionHandler.GetBidderAuctionList)
		// オークション詳細取得（すべてのユーザーがアクセス可能）
		api.GET("/auctions/:id", auctionHandler.GetAuctionDetail)

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

			// 入札者専用エンドポイント
			bidder := protected.Group("/bidder")
			bidder.Use(middleware.RequireBidder())
			{
				bidder.GET("/points", bidHandler.GetPoints)
				bidder.POST("/items/:id/bid", bidHandler.PlaceBid)
				bidder.GET("/items/:id/bids", bidHandler.GetBidHistory)
			}

			// システム管理者専用エンドポイント
			systemAdmin := protected.Group("")
			systemAdmin.Use(middleware.RequireSystemAdmin())
			{
				// 管理者登録
				systemAdmin.POST("/admin/admins", adminHandler.RegisterAdmin)
				// 管理者一覧取得
				systemAdmin.GET("/admin/admins", adminHandler.GetAdminList)
				// 管理者詳細取得
				systemAdmin.GET("/admin/admins/:id", adminHandler.GetAdmin)
				// 管理者更新
				systemAdmin.PUT("/admin/admins/:id", adminHandler.UpdateAdmin)
				// 管理者状態変更
				systemAdmin.PATCH("/admin/admins/:id/status", adminHandler.UpdateAdminStatus)

				// 入札者登録
				systemAdmin.POST("/admin/bidders", bidderHandler.RegisterBidder)
				// 入札者一覧取得
				systemAdmin.GET("/admin/bidders", bidderHandler.GetBidderList)
				// 入札者詳細取得
				systemAdmin.GET("/admin/bidders/:id", bidderHandler.GetBidderByID)
				// 入札者更新
				systemAdmin.PUT("/admin/bidders/:id", bidderHandler.UpdateBidder)
				// 入札者へのポイント付与
				systemAdmin.POST("/admin/bidders/:id/points", bidderHandler.GrantPoints)
				// 入札者のポイント履歴取得
				systemAdmin.GET("/admin/bidders/:id/points/history", bidderHandler.GetPointHistory)
				// 入札者状態変更
				systemAdmin.PATCH("/admin/bidders/:id/status", bidderHandler.UpdateBidderStatus)

				// オークション中止（system_adminのみ）
				systemAdmin.POST("/admin/auctions/:id/cancel", auctionHandler.CancelAuctionWithReason)
			}

			// システム管理者と主催者がアクセス可能なエンドポイント
			adminOrAuctioneer := protected.Group("")
			adminOrAuctioneer.Use(middleware.RequireAdminOrAuctioneer())
			{
				// オークション一覧取得
				adminOrAuctioneer.GET("/admin/auctions", auctionHandler.GetAuctionList)
				// オークション作成
				adminOrAuctioneer.POST("/admin/auctions", auctionHandler.CreateAuction)
				// オークション詳細取得（編集用）
				adminOrAuctioneer.GET("/admin/auctions/:id", auctionHandler.GetAuctionForEdit)
				// オークション更新
				adminOrAuctioneer.PUT("/admin/auctions/:id", auctionHandler.UpdateAuction)
				// オークション開始
				adminOrAuctioneer.POST("/admin/auctions/:id/start", auctionHandler.StartAuction)
				// オークション終了
				adminOrAuctioneer.POST("/admin/auctions/:id/end", auctionHandler.EndAuction)
				// オークション参加者一覧取得
				adminOrAuctioneer.GET("/admin/auctions/:id/participants", auctionHandler.GetParticipants)

				// オークション商品紐づけ
				adminOrAuctioneer.POST("/admin/auctions/:id/items/assign", itemHandler.AssignItems)
				// オークション商品解除
				adminOrAuctioneer.DELETE("/admin/auctions/:id/items/:itemId/unassign", itemHandler.UnassignItem)

				// 商品追加
				adminOrAuctioneer.POST("/admin/auctions/:id/items", auctionHandler.AddItem)
				// 商品更新
				adminOrAuctioneer.PUT("/admin/auctions/:id/items/:itemId", auctionHandler.UpdateItem)
				// 商品削除
				adminOrAuctioneer.DELETE("/admin/auctions/:id/items/:itemId", auctionHandler.DeleteItem)
				// 商品順序変更
				adminOrAuctioneer.PUT("/admin/auctions/:id/items/reorder", auctionHandler.ReorderItems)

				// 商品管理API（オークションから独立）
				adminOrAuctioneer.GET("/admin/items", itemHandler.GetItemList)
				adminOrAuctioneer.POST("/admin/items", itemHandler.CreateItem)
				adminOrAuctioneer.GET("/admin/items/:id", itemHandler.GetItemDetail)
				adminOrAuctioneer.PUT("/admin/items/:id", itemHandler.UpdateItem)
				adminOrAuctioneer.DELETE("/admin/items/:id", itemHandler.DeleteItem)

				// 商品開始
				adminOrAuctioneer.POST("/admin/items/:id/start", auctionHandler.StartItem)
				// 価格開示
				adminOrAuctioneer.POST("/admin/items/:id/open-price", auctionHandler.OpenPrice)
				// 商品終了
				adminOrAuctioneer.POST("/admin/items/:id/end", auctionHandler.EndItem)
				// 入札履歴取得
				adminOrAuctioneer.GET("/admin/items/:id/bids", auctionHandler.GetBidHistory)
				// 価格開示履歴取得
				adminOrAuctioneer.GET("/admin/items/:id/price-history", auctionHandler.GetPriceHistory)
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
