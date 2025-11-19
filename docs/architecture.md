# システムアーキテクチャ

## 概要

競走馬セリをモデルとしたリアルタイムオークションシステムのアーキテクチャ設計書です。
仮想ポイント制を採用し、主催者主導で価格を開示していく方式のオークションシステムです。
WebSocketによるリアルタイム通信で、入札状況や価格変動を即座に反映します。

**オークション方式の特徴**
- **ポイント制**: 仮想ポイントで入札（実際の金銭取引なし）
- **主催者主導**: 主催者が開始価格と次の入札価格を決定・開示
- **タイマーレス**: 制限時間なし、主催者の判断で終了
- **競り下げも可能**: 入札がなければ前の価格で入札したユーザーが落札

**開発フェーズ**
- **第一フェーズ**: Webアプリケーション (Vue.js) の開発
- **第二フェーズ**: ネイティブモバイルアプリ (Swift/iOS) の開発

## システム全体構成図

```
┌─────────────────────────────────────────────────────────────────┐
│                        クライアント層                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │   Web Client (Phase 1)                                 │    │
│  │   (Vue.js 3 + Vite)                                    │    │
│  │                                                        │    │
│  │  - オークション閲覧 (全ユーザー)                          │    │
│  │  - 管理機能 (管理者・主催者)                             │    │
│  │  - 入札機能 (入札者)                                     │    │
│  │  - レスポンシブデザイン (PC/タブレット/スマホ対応)         │    │
│  └────────────────────────────────────────────────────────┘    │
│         │                                                       │
│         │ HTTPS/WSS                                            │
│         │                                                       │
│  ┌─────┴──────────────────────────────────────────────────┐    │
│  │   iOS App (Phase 2 - 予定)                             │    │
│  │   (Swift + SwiftUI)                                    │    │
│  │                                                        │    │
│  │  - 入札専用アプリ                                        │    │
│  │  - プッシュ通知 (APNs)                                   │    │
│  │  - リアルタイム更新                                      │    │
│  │  - オフライン対応                                        │    │
│  └────────────────────────────────────────────────────────┘    │
│                                                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          │ HTTPS/WSS
                          │
┌─────────────────────┴───────────────────────────────────────────┐
│                      API Gateway層                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              Nginx / AWS ALB / Cloud Load Balancer       │  │
│  │                                                          │  │
│  │  - SSL終端                                                │  │
│  │  - ルーティング (REST API / WebSocket)                    │  │
│  │  - レート制限                                              │  │
│  │  - 負荷分散                                                │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
└─────────┬───────────────────────────────────┬───────────────────┘
          │                                   │
          │ HTTP/REST                         │ WebSocket
          │                                   │
┌─────────┴─────────────────┐    ┌────────────┴─────────────────┐
│    アプリケーション層        │    │  リアルタイム通信層           │
├───────────────────────────┤    ├──────────────────────────────┤
│                           │    │                              │
│  ┌─────────────────────┐  │    │  ┌────────────────────────┐ │
│  │   REST API Server   │  │    │  │  WebSocket Server      │ │
│  │   (Node.js/Express) │  │    │  │  (Node.js/Socket.io)   │ │
│  │                     │  │    │  │                        │ │
│  │  エンドポイント:      │  │    │  │  イベント:              │ │
│  │  - 認証/認可         │  │    │  │  - auction:join        │ │
│  │  - ユーザー管理       │  │    │  │  - auction:bid         │ │
│  │  - オークションCRUD   │  │    │  │  - auction:update      │ │
│  │  - 商品(馬)管理      │  │    │  │  - auction:timer       │ │
│  │  - 入札履歴取得      │  │    │  │  - auction:extended    │ │
│  │  - レポート生成      │  │    │  │  - auction:ended       │ │
│  └─────────────────────┘  │    │  └────────────────────────┘ │
│            │              │    │             │                │
│            │              │    │             │                │
│  ┌─────────┴────────────┐ │    │  ┌──────────┴──────────────┐ │
│  │  Business Logic      │ │    │  │   Connection Manager    │ │
│  │                      │ │    │  │                         │ │
│  │  - 入札検証ロジック    │ │    │  │  - WebSocket接続管理    │ │
│  │  - タイマー延長ロジック │ │    │  │  - ルーム管理           │ │
│  │  - 落札判定          │ │    │  │  - ブロードキャスト      │ │
│  │  - 通知トリガー       │ │    │  │  - 再接続ハンドリング    │ │
│  └──────────────────────┘ │    │  └─────────────────────────┘ │
│                           │    │                              │
└───────────┬───────────────┘    └────────────┬─────────────────┘
            │                                 │
            │                                 │
            └────────────┬────────────────────┘
                         │
┌────────────────────────┴─────────────────────────────────────┐
│                        データ層                                │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────┐      ┌──────────────────────┐    │
│  │  Primary Database    │      │   Cache / Session    │    │
│  │  (PostgreSQL)        │      │   (Redis)            │    │
│  │                      │      │                      │    │
│  │  テーブル:            │      │  用途:                │    │
│  │  - users             │      │  - セッション管理      │    │
│  │  - user_points       │      │  - リアルタイム状態    │    │
│  │  - auctions          │      │  - 入札キュー         │    │
│  │  - items (horses)    │      │  - 現在価格状態       │    │
│  │  - bids              │      │  - アクティブ接続情報  │    │
│  │  - price_history     │      │  - Pub/Sub           │    │
│  │  - notifications     │      │                      │    │
│  └──────────────────────┘      └──────────────────────┘    │
│           │                              │                  │
│           │                              │                  │
│  ┌────────┴──────────────────────────────┴────────────┐    │
│  │              Message Queue (Optional)              │    │
│  │              (Redis / RabbitMQ / AWS SQS)          │    │
│  │                                                    │    │
│  │  - 非同期処理 (メール送信、通知、レポート生成)        │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│                      外部サービス層                            │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  - APNs (Apple Push Notification service) - Phase 2         │
│  - SendGrid / AWS SES (メール通知)                           │
│  - Stripe (決済処理 - オプション)                             │
│  - AWS S3 / Cloud Storage (画像・ドキュメント保存)            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## アーキテクチャの特徴

### 1. マイクロサービス指向の分離設計

**REST APIサーバーとWebSocketサーバーを分離**することで、以下のメリットを実現：

- **独立したスケーリング**: リアルタイム通信の負荷に応じて個別にスケールアウト可能
- **障害の局所化**: 一方のサービス障害が他方に影響しない
- **技術選択の柔軟性**: 各サービスに最適な技術を選択可能

### 2. リアルタイム通信の実装

**WebSocket (Gorilla WebSocket)** を使用したリアルタイム双方向通信：

- **低レイテンシ**: 入札から通知まで100ms以内の応答（Goの高速処理）
- **永続的接続**: goroutineによる効率的な接続管理
- **軽量**: 数万の同時接続を1サーバーで処理可能
- **ルームベース管理**: チャネルを使った安全な並行処理
- **標準準拠**: RFC 6455完全準拠の実装
- **クライアント対応**: Vue.js (標準WebSocket API) / iOS Swift (URLSessionWebSocketTask)

### 3. Redis活用による高速性とスケーラビリティ

**Redis** を多目的に活用：

- **セッション管理**: 分散環境でのユーザーセッション共有
- **リアルタイム状態管理**: オークションの現在価格、入札状況をキャッシュ
- **Pub/Sub**: 複数WebSocketサーバー間でのイベント配信（go-redis使用）
- **入札キュー**: Goの並行処理とRedisロックによる競合制御と順序保証
- **高速アクセス**: Goのgoroutineと組み合わせて非同期処理を実現

### 4. 価格管理アーキテクチャ

**主催者主導の価格開示システム**：

```
┌─────────────────────────────────────────────────────┐
│           価格管理システム                             │
├─────────────────────────────────────────────────────┤
│                                                     │
│  1. オークション開始時                                │
│     ↓                                               │
│     主催者: 開始価格を設定 (例: 100pt)                │
│     Redis: SET auction:{id}:current_price 100       │
│     Redis: SET auction:{id}:status "active"         │
│     WebSocket: broadcast "auction:started"          │
│                                                     │
│  2. 入札イベント受信                                  │
│     ↓                                               │
│     ポイント残高確認                                  │
│     PostgreSQL: INSERT bid record                   │
│     Redis: SET auction:{id}:has_bid true            │
│     Redis: SET auction:{id}:last_bidder {user_id}   │
│     WebSocket: broadcast "auction:bid"              │
│                                                     │
│  3. 主催者が次の価格を開示                            │
│     ↓                                               │
│     主催者: 次の価格を設定 (例: 150pt)                │
│     Redis: SET auction:{id}:current_price 150       │
│     PostgreSQL: INSERT price_history                │
│     WebSocket: broadcast "auction:price_open"       │
│                                                     │
│  4. 入札がない場合                                    │
│     ↓                                               │
│     主催者: 「入札なし」を確認                        │
│     前の価格で入札したユーザーを落札者として確定        │
│     PostgreSQL: UPDATE auction SET winner           │
│     ポイント確定処理                                  │
│     WebSocket: broadcast "auction:ended"            │
│                                                     │
│  5. 主催者が手動で終了                                │
│     ↓                                               │
│     最後の入札者を落札者として確定                     │
│     ポイント確定処理                                  │
│     WebSocket: broadcast "auction:ended"            │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 5. データ整合性とトランザクション

**PostgreSQLとRedisの適材適所な使い分け**：

- **PostgreSQL**: 永続化が必要なマスターデータ、入札履歴、監査ログ
- **Redis**: リアルタイム性が必要な一時データ、キャッシュ
- **トランザクション**: 入札処理時はPostgreSQLトランザクションで整合性保証

### 6. 水平スケーラビリティ

**ステートレス設計**によるスケールアウト対応：

```
┌──────────────────────────────────────────────────┐
│         Load Balancer (Sticky Session)           │
└───────┬──────────────┬──────────────┬───────────┘
        │              │              │
   ┌────┴───┐    ┌────┴───┐    ┌────┴───┐
   │ WS #1  │    │ WS #2  │    │ WS #3  │
   └────┬───┘    └────┬───┘    └────┬───┘
        │              │              │
        └──────────────┴──────────────┘
                       │
              ┌────────┴────────┐
              │  Redis Pub/Sub  │  ← イベント配信
              └─────────────────┘
```

## データフロー

### 入札処理フロー

```
入札者 (Client)
  │
  │ 1. bid event (WebSocket) + current_price
  ↓
WebSocket Server
  │
  │ 2. validate connection & auth (role: bidder)
  ↓
Business Logic
  │
  │ 3. validate bid
  │    - 現在の開示価格と一致するか
  │    - ポイント残高が十分か
  │    - オークションがactiveか
  ↓
PostgreSQL
  │ 4. BEGIN TRANSACTION
  │ 5. INSERT INTO bids (auction_id, user_id, price, points)
  │ 6. UPDATE user_points SET reserved_points += price
  │ 7. COMMIT
  ↓
Redis
  │ 8. SET auction:{id}:has_bid true
  │ 9. SET auction:{id}:last_bidder {user_id}
  │ 10. SET auction:{id}:last_bid_price {price}
  ↓
Redis Pub/Sub
  │ 11. PUBLISH auction:bid {auction_id, user_id, price}
  ↓
All WebSocket Server Instances
  │ 12. broadcast to auction room
  ↓
All Connected Clients
  │ 13. UI update (新しい入札を表示)
```

### オークション開始フロー

```
主催者 (Auctioneer)
  │
  │ 1. POST /api/auctions/:id/start
  │    { starting_price: 100 }
  ↓
REST API Server
  │
  │ 2. validate permissions (role: auctioneer)
  │ 3. validate starting_price
  ↓
PostgreSQL
  │ 4. UPDATE auctions SET
  │    status = 'active',
  │    started_at = NOW()
  │ 5. INSERT INTO price_history
  │    (auction_id, price, opened_by, opened_at)
  ↓
Redis
  │ 6. SET auction:{id}:status "active"
  │ 7. SET auction:{id}:current_price 100
  │ 8. SET auction:{id}:has_bid false
  ↓
Redis Pub/Sub
  │ 9. PUBLISH auction:started {auction_id, price: 100}
  ↓
WebSocket Servers
  │ 10. broadcast to all connected clients
  ↓
All Clients
  │ 11. show auction UI with starting price
  │ 12. enable bid button for bidders
```

### 価格開示フロー

```
主催者 (Auctioneer)
  │
  │ 1. POST /api/auctions/:id/open-price
  │    { next_price: 150 }
  ↓
REST API Server
  │
  │ 2. validate permissions (role: auctioneer)
  │ 3. validate next_price > current_price
  ↓
PostgreSQL
  │ 4. INSERT INTO price_history
  │    (auction_id, price, opened_by, opened_at)
  ↓
Redis
  │ 5. GET auction:{id}:has_bid
  │ 6. SET auction:{id}:current_price 150
  │ 7. SET auction:{id}:has_bid false
  │ 8. DEL auction:{id}:last_bidder (reset)
  ↓
Redis Pub/Sub
  │ 9. PUBLISH auction:price_open
  │    {auction_id, new_price: 150, had_previous_bid}
  ↓
WebSocket Servers
  │ 10. broadcast to auction room
  ↓
All Clients
  │ 11. UI update (new price displayed)
  │ 12. reset bid status display
```

## セキュリティ設計

### 認証・認可

1. **JWT (JSON Web Token)** による認証
   - アクセストークン (短期: 15分)
   - リフレッシュトークン (長期: 7日)

2. **WebSocket接続の認証**
   - ハンドシェイク時にJWTトークンを検証
   - 接続後も定期的にトークンの有効性をチェック

3. **ロールベースアクセス制御 (RBAC)**
   - `system_admin`: システム管理者（全権限、ユーザー管理、ポイント付与、全体状況確認）
   - `auctioneer`: オークション主催者（オークション作成・開始・終了、商品登録、価格開示）
   - `bidder`: 入札者（入札のみ、自分のポイント・入札履歴閲覧）

### データ保護

- **通信の暗号化**: TLS 1.3 (HTTPS/WSS)
- **パスワードハッシュ化**: bcrypt (cost factor: 12)
- **APIレート制限**: IPアドレスごとに 100 req/min
- **入札レート制限**: ユーザーごとに 10 bids/min

### 不正防止

- **二重入札防止**: Redisロックによる排他制御
- **ポイント残高検証**: 入札時にリアルタイムで残高確認
- **価格整合性チェック**: 開示価格との一致を検証
- **ロール権限チェック**: 主催者のみが価格開示・終了可能
- **監査ログ**: 全入札イベント、価格開示、ポイント付与を記録
- **ポイント不正利用防止**: システム管理者のみがポイント付与可能

## パフォーマンス要件

| 項目 | 目標値 |
|------|--------|
| WebSocket接続確立 | < 100ms |
| 入札イベント処理 | < 50ms |
| 入札通知配信 | < 100ms |
| 価格開示通知配信 | < 100ms |
| REST API応答時間 | < 200ms (P95) |
| ポイント残高照会 | < 50ms |
| 同時接続数 | 10,000+ connections |
| 同時開催オークション数 | 100+ auctions |
| データベースクエリ | < 50ms (P95) |
| Redisクエリ | < 5ms (P99) |

## 監視・運用

### メトリクス収集

- **APM**: New Relic / Datadog
- **ログ集約**: ELK Stack / CloudWatch Logs
- **メトリクス**: Prometheus + Grafana

### 重要な監視項目

1. **WebSocket接続数**: アクティブ接続数の推移
2. **入札レート**: 秒間入札数
3. **応答時間**: エンドツーエンドのレイテンシ
4. **エラー率**: 4xx, 5xxエラーの発生率
5. **Redis/DB接続プール**: 使用率と待機時間
6. **CPU/メモリ使用率**: リソース使用状況

### アラート設定

- WebSocket切断率が10%を超えた場合
- 入札処理の失敗率が5%を超えた場合
- API応答時間が500msを超えた場合
- データベース接続プールが90%を超えた場合

## デプロイ構成

### 推奨環境: AWS

```
┌────────────────────────────────────────────────────────┐
│                     AWS Cloud                          │
├────────────────────────────────────────────────────────┤
│                                                        │
│  ┌──────────────────────────────────────────────┐    │
│  │     CloudFront (CDN)                         │    │
│  │     - Static assets caching                  │    │
│  │     - DDoS protection                        │    │
│  └──────────────────┬───────────────────────────┘    │
│                     │                                 │
│  ┌──────────────────┴───────────────────────────┐    │
│  │  Application Load Balancer                   │    │
│  │  - SSL termination                           │    │
│  │  - Path-based routing                        │    │
│  │  - Health checks                             │    │
│  └──────────────────┬───────────────────────────┘    │
│                     │                                 │
│         ┌───────────┴───────────┐                    │
│         │                       │                    │
│  ┌──────┴──────┐         ┌─────┴──────┐            │
│  │   ECS/Fargate         │   ECS/Fargate            │
│  │   (REST API)  │         │ (WebSocket) │            │
│  │   Auto Scaling│         │ Auto Scaling│            │
│  └──────┬──────┘         └─────┬──────┘            │
│         │                       │                    │
│         └───────────┬───────────┘                    │
│                     │                                 │
│  ┌──────────────────┴───────────────────────────┐    │
│  │  ElastiCache for Redis (Cluster Mode)       │    │
│  │  - Multi-AZ                                  │    │
│  │  - Automatic failover                        │    │
│  └──────────────────────────────────────────────┘    │
│                                                        │
│  ┌────────────────────────────────────────────┐      │
│  │  RDS for PostgreSQL                        │      │
│  │  - Multi-AZ deployment                     │      │
│  │  - Read replicas                           │      │
│  │  - Automated backups                       │      │
│  └────────────────────────────────────────────┘      │
│                                                        │
│  ┌────────────────────────────────────────────┐      │
│  │  S3                                        │      │
│  │  - Static assets                           │      │
│  │  - Auction images                          │      │
│  │  - Logs & backups                          │      │
│  └────────────────────────────────────────────┘      │
│                                                        │
└────────────────────────────────────────────────────────┘
```

### 代替環境: GCP / 自前VPS

- **GCP**: Cloud Run + Cloud SQL + Memorystore
- **自前VPS**: Docker Swarm / Kubernetes + NGINX + PostgreSQL + Redis

## 拡張性・将来対応

1. **マルチテナント対応**: 複数の主催者が独立してオークションを開催
2. **ビデオストリーミング**: 競走馬の映像をライブ配信
3. **AIレコメンデーション**: 過去の入札履歴から興味のある馬を推薦
4. **ブロックチェーン**: 入札履歴の改ざん防止と透明性確保
5. **国際化**: 多言語・多通貨対応

## まとめ

本アーキテクチャは以下の要件を満たします：

✅ **リアルタイム性**: WebSocketによる双方向通信で瞬時の入札・価格開示を反映  
✅ **ポイント制**: 仮想ポイントによる安全なオークション運営  
✅ **主催者主導**: 主催者が価格をコントロールする柔軟な運営  
✅ **3ロール体制**: システム管理者・主催者・入札者の明確な権限分離  
✅ **スケーラビリティ**: ステートレス設計で水平スケール可能  
✅ **可用性**: Multi-AZ構成と自動フェイルオーバー  
✅ **セキュリティ**: JWT認証、TLS暗号化、ロール別アクセス制御  
✅ **マルチプラットフォーム**:  
   - **Phase 1**: Web (Vue.js 3 + Vite) - レスポンシブ対応  
   - **Phase 2**: iOS App (Swift + SwiftUI) - ネイティブアプリ

競走馬セリの要件に最適化された、主催者主導型のエンタープライズグレードのアーキテクチャです。

## バックエンド技術スタック

### Go言語採用の理由

**パフォーマンス**
- ネイティブコンパイルによる高速実行
- 低レイテンシ（入札処理 < 50ms）
- 効率的なメモリ管理

**並行処理**
- goroutineによる軽量スレッド（数万の同時接続に対応）
- チャネルを使った安全な並行処理
- WebSocket接続管理に最適

**スケーラビリティ**
- 単一バイナリで簡単デプロイ
- 水平スケールが容易
- Kubernetes/Docker対応

**信頼性**
- 静的型付けによる安全性
- 明示的なエラーハンドリング
- 大規模本番環境での実績

### 主要ライブラリ・フレームワーク

**Webフレームワーク: Gin**
- 高速なHTTPルーター（httprouterベース）
- 充実したミドルウェア
- JSONバインディング・バリデーション
- 豊富なコミュニティとドキュメント

**WebSocket: Gorilla WebSocket**
- RFC 6455完全準拠
- 安定した実装と豊富な実績
- 柔軟なメッセージハンドリング

**データベース: GORM**
- 型安全なORM
- マイグレーション機能
- リレーション管理

**Redis: go-redis**
- 高性能なRedisクライアント
- Pub/Sub対応
- コネクションプーリング

**認証: golang-jwt/jwt**
- JWT生成・検証
- 標準的な実装

**バリデーション: go-playground/validator**
- 構造体ベースのバリデーション
- カスタムルール対応

### バックエンドプロジェクト構造

```
backend/
├── cmd/
│   ├── api/              # REST APIサーバー
│   │   └── main.go
│   └── ws/               # WebSocketサーバー
│       └── main.go
├── internal/
│   ├── domain/           # ドメインモデル
│   │   ├── auction.go
│   │   ├── user.go
│   │   ├── bid.go
│   │   └── point.go
│   ├── repository/       # データアクセス層
│   │   ├── postgres/
│   │   └── redis/
│   ├── service/          # ビジネスロジック
│   │   ├── auction_service.go
│   │   ├── bid_service.go
│   │   └── point_service.go
│   ├── handler/          # HTTPハンドラー
│   │   ├── auction_handler.go
│   │   ├── user_handler.go
│   │   └── auth_handler.go
│   ├── ws/               # WebSocketハンドラー
│   │   ├── hub.go        # 接続管理
│   │   ├── client.go     # クライアント管理
│   │   └── handler.go    # イベント処理
│   └── middleware/       # ミドルウェア
│       ├── auth.go
│       ├── cors.go
│       └── logger.go
├── pkg/
│   ├── config/           # 設定管理
│   ├── logger/           # ロギング
│   └── validator/        # カスタムバリデーター
├── migrations/           # DBマイグレーション
├── docs/                 # API仕様書
├── go.mod
└── go.sum
```

## フェーズ別開発計画

### Phase 1: Webアプリケーション (優先)

**目的**: 全ユーザー向けの基本機能を提供

**技術スタック**
- **フロントエンド**: Vue.js 3 (Composition API) + Vite
- **UIフレームワーク**: Vuetify 3 / Element Plus / Ant Design Vue
- **状態管理**: Pinia
- **WebSocket**: 標準WebSocket API (RFC 6455準拠)
- **HTTP Client**: Axios
- **認証**: JWT (localStorage/sessionStorage)

**対応機能**
- オークション一覧・詳細閲覧
- リアルタイム入札
- システム管理画面 (ユーザー管理・ポイント付与・全体状況確認)
- 主催者管理画面 (オークション作成・開始・終了・価格開示・商品登録)
- ポイント残高表示
- 入札履歴閲覧
- ユーザー登録・認証
- レスポンシブデザイン (モバイルブラウザ対応)

**配信方法**
- CloudFront + S3 (静的ホスティング)
- PWA対応でモバイルでもアプリライクな体験

### Phase 2: iOSネイティブアプリ (将来)

**目的**: 入札者向けの専用アプリで最高のUXを提供

**技術スタック**
- **言語**: Swift 5+
- **UIフレームワーク**: SwiftUI
- **アーキテクチャ**: MVVM + Combine
- **WebSocket**: URLSessionWebSocketTask (iOS 13+) / Starscream
- **HTTP Client**: URLSession / Alamofire
- **認証**: Keychain による安全なトークン保存
- **プッシュ通知**: APNs (Apple Push Notification service)

**追加機能**
- プッシュ通知 (オークション開始・入札通知)
- オフライン時の入札キュー
- Face ID / Touch ID 認証
- ウィジェット対応 (進行中オークションの表示)
- Apple Watch対応 (Phase 2.1)

**配信方法**
- App Store配信
- TestFlight によるベータテスト
