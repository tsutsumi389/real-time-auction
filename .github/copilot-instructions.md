# Real-Time Auction System - AI Coding Agent Instructions

## Project Overview

競走馬セリをモデルとした主催者主導型のリアルタイムオークションシステム。仮想ポイント制で、WebSocketによる双方向通信を採用。

**Key Architectural Principles:**
- **Separated Services**: REST API (`cmd/api`) と WebSocket (`cmd/ws`) を独立したGoサーバーで実装
- **Auctioneer-Driven**: 主催者が開始価格と次の入札価格を決定・開示（タイマーレス）
- **3-Role System**: `system_admin` (全権限)、`auctioneer` (オークション管理)、`bidder` (入札のみ)
- **Point-Based**: 仮想ポイントで入札（実際の金銭取引なし）

## Development Workflow

### Quick Start
```bash
# 初回セットアップ & 起動
make up

# アクセス: http://localhost (Frontend), http://localhost/api (REST), ws://localhost/ws (WebSocket)
```

### Essential Make Commands
- `make up` - 全サービス起動 (初回は自動的に `.env` 作成)
- `make logs` - 全ログ表示 / `make logs service=api` - 特定サービスのみ
- `make down` - 停止 / `make clean` - ボリューム含めて完全削除
- `make shell-api` / `make shell-ws` - コンテナ内シェル
- `make shell-postgres` - PostgreSQL接続 / `make shell-redis` - Redis CLI

**Important**: コマンドはプロジェクトルートから実行。`Makefile` に全てのワークフローが定義済み。

## Backend (Go)

### Architecture Pattern
```
cmd/
  api/main.go        # REST APIサーバー (Gin) - ポート 8080
  ws/main.go         # WebSocketサーバー (Gin + Gorilla WebSocket) - ポート 8081
internal/
  domain/            # ドメインモデル (struct定義)
  repository/        # データアクセス層 (GORM, go-redis)
    postgres/
    redis/
  service/           # ビジネスロジック (入札検証、ポイント管理)
  handler/           # HTTPハンドラー (Gin handlers)
  ws/
    hub.go           # WebSocket接続管理 (goroutineベース)
    client.go        # クライアント管理
    handler.go       # イベント処理 (auction:join, auction:bid等)
  middleware/        # JWT認証、CORS、ロギング
pkg/                 # 共通パッケージ (config, logger, validator)
```

### Technology Stack
- **Gin**: HTTPルーター (高速、ミドルウェア充実)
- **Gorilla WebSocket**: RFC 6455準拠、goroutineと組み合わせて数万接続に対応
- **GORM**: ORM (PostgreSQL)
- **go-redis**: Redisクライアント (Pub/Sub、キャッシュ、セッション管理)
- **golang-jwt/jwt**: JWT認証
- **go-playground/validator**: 構造体バリデーション

### Hot Reload
- **Air** によるホットリロード設定済み
  - `.air.toml` - REST APIサーバー用
  - `.air.ws.toml` - WebSocketサーバー用
- `make up` で自動的にAirが起動
- `.go` ファイルの変更を検知して自動再ビルド

### Code Conventions
- **Error Handling**: 明示的な `if err != nil` チェック必須
- **Context**: goroutineには必ず `context.Context` を渡す
- **Logging**: `pkg/logger` を使用 (標準出力ではなく構造化ログ)
- **環境変数**: `os.Getenv()` でフォールバック値を提供

### Data Flow Example (入札処理)
1. Client → WebSocket: `{"type":"auction:bid", "auction_id":1, "price":150}`
2. `ws/handler.go`: JWT検証、ロール確認 (`bidder`)
3. `service/bid_service.go`: ポイント残高確認、開示価格との一致確認
4. PostgreSQL: `INSERT INTO bids` + `UPDATE user_points SET reserved_points += price` (トランザクション)
5. Redis: `SET auction:{id}:last_bidder {user_id}`, `SET auction:{id}:has_bid true`
6. Redis Pub/Sub: `PUBLISH auction:bid {auction_id, user_id, price}`
7. All WebSocket Servers: ブロードキャスト (`hub.go` の goroutine経由)

## Frontend (Vue.js 3)

### Architecture
```
src/
  views/            # ページコンポーネント (HomeView.vue)
  components/       # 再利用可能コンポーネント (未実装)
  router/           # Vue Router設定
  stores/           # Pinia状態管理 (未実装)
  services/         # API通信 (axios) (未実装)
  assets/           # CSS、画像
```

### Technology Stack
- **Vue 3 Composition API**: `<script setup>` スタイル
- **Vite**: ビルドツール、HMR対応
- **Pinia**: 状態管理 (未実装)
- **Axios**: HTTP通信 (未実装)
- **標準WebSocket API**: `new WebSocket(import.meta.env.VITE_WS_URL)`

### Environment Variables
- `VITE_API_BASE_URL` - REST APIベースURL (デフォルト: `http://localhost/api`)
- `VITE_WS_URL` - WebSocket URL (デフォルト: `ws://localhost/ws`)
- `.env.example` の `VITE_*` 変数を参照

### Code Conventions
- **Composition API** のみ使用 (Options API禁止)
- **TypeScript未導入**: 現在はJavaScriptのみ (将来的にTS化予定)

## Nginx (Reverse Proxy)

### Routing Rules
```
/api/*  → api:8080  (REST API、CORS設定、60秒タイムアウト)
/ws     → ws:8081   (WebSocket、7日間タイムアウト、buffering無効)
/*      → frontend:5173 (Vite dev server、HMR対応)
```

**CORS**: 開発環境では全オリジン許可 (`*`)、本番では要制限

## Database & Cache

### PostgreSQL (GORM)
- **Host**: `postgres:5432` (Docker内)
- **DB**: `auction_db` / **User**: `auction_user` / **Password**: `.env` 参照
- **接続**: `make shell-postgres` でpsql接続
- **マイグレーション**: `migrations/` (未実装)

**Planned Schema:**
```
users (id, email, password_hash, role, created_at)
user_points (user_id, total_points, available_points, reserved_points)
auctions (id, status, started_at, ended_at, winner_id)
items (id, name, description, image_url, auction_id)
bids (id, auction_id, user_id, price, bid_at)
price_history (id, auction_id, price, opened_by, opened_at)
```

### Redis (go-redis)
- **Host**: `redis:6379` (Docker内)
- **接続**: `make shell-redis` でredis-cli接続
- **用途**:
  - セッション管理: `session:{token}`
  - オークション状態: `auction:{id}:current_price`, `auction:{id}:status`, `auction:{id}:has_bid`, `auction:{id}:last_bidder`
  - Pub/Sub: `auction:started`, `auction:bid`, `auction:price_open`, `auction:ended`
  - ロック: 入札時の排他制御

## Current Development Status

**Implemented:**
- ✅ Docker開発環境 (docker-compose.yml)
- ✅ Makefile による開発ワークフロー
- ✅ 基本的なGoサーバー構造 (`cmd/api`, `cmd/ws`)
- ✅ Vue 3 + Vite基盤
- ✅ Nginx リバースプロキシ設定

**Not Yet Implemented (最優先タスク):**
- ❌ PostgreSQL スキーマ定義 & マイグレーション
- ❌ JWT認証ミドルウェア
- ❌ WebSocket実装 (Gorilla WebSocket、hub.go、client.go)
- ❌ GORM リポジトリ層 (`internal/repository/postgres`)
- ❌ Redis リポジトリ層 (`internal/repository/redis`)
- ❌ ビジネスロジック (`internal/service`)
- ❌ REST APIエンドポイント実装
- ❌ フロントエンド画面実装

## Critical Implementation Notes

### WebSocket Server Design
- **Hub Pattern**: `ws/hub.go` でgoroutineベースの接続管理
  - `clients map[*Client]bool` - アクティブ接続
  - `rooms map[int][]*Client` - オークションIDごとのルーム
  - `broadcast chan []byte` - ブロードキャストチャネル
- **Client**: `ws/client.go` で各接続を管理
  - `readPump()` goroutine - クライアントからのメッセージ受信
  - `writePump()` goroutine - クライアントへのメッセージ送信
- **Ping/Pong**: 30秒間隔でping、60秒タイムアウト

### Auctioneer-Driven Price Management
主催者が価格を**開示**する独自システム:
1. 主催者: 開始価格設定 → `POST /api/auctions/:id/start` → Redis: `SET auction:{id}:current_price`
2. 入札者: 開示価格で入札 → WebSocket: `auction:bid` イベント → リクエストに `price` 含める（開示価格と一致確認）
3. 主催者: 次の価格開示 → `POST /api/auctions/:id/open-price` → 前の価格で入札があったか確認 (`has_bid`)
4. 入札なし → 前の価格で入札したユーザーを落札者として確定 → `auction:ended` イベント

**重要**: 入札者は開示された価格でしか入札できない（自由な価格入札は不可）。

### Role-Based Access Control
- **system_admin**: ユーザー管理、ポイント付与、全体状況確認
- **auctioneer**: オークション作成・開始・終了、商品登録、価格開示
- **bidder**: 入札のみ、自分のポイント・入札履歴閲覧

JWT Claimsに `role` を含める。各エンドポイント/WebSocketイベントでロール確認必須。

## Testing & Debugging

- **Health Checks**:
  - `curl http://localhost/api/health` - REST API
  - `curl http://localhost/ws/health` - WebSocket
- **Logs**: `make logs` または `make logs service=api`
- **Database**: `make shell-postgres` → `\dt` (テーブル一覧), `\d+ users` (スキーマ詳細)
- **Redis**: `make shell-redis` → `KEYS *`, `GET auction:1:current_price`

## Common Pitfalls

1. **環境変数**: `.env` がないと `make up` で自動作成されるが、`JWT_SECRET` は本番では必ず変更
2. **Docker内からのDB接続**: ホスト名は `postgres:5432` (localhost不可)
3. **WebSocketアップグレード**: Nginxの `Upgrade` ヘッダー設定必須 (nginx.conf参照)
4. **CORS**: 開発環境では `*` 許可だが、本番では `CORS_ORIGINS` を制限
5. **goroutine leak**: WebSocket接続終了時に `readPump()`, `writePump()` を必ずクリーンアップ

## Future Phases

- **Phase 1** (現在): Webアプリケーション (Vue.js 3)
- **Phase 2** (将来): iOSネイティブアプリ (Swift + SwiftUI、APNsプッシュ通知)

Webアプリはレスポンシブデザインでモバイルブラウザ対応だが、iOSアプリは入札者専用の最適化されたUXを提供予定。
