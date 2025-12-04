# Real-Time Auction System

競走馬セリをモデルとしたリアルタイムオークションシステム

## 概要

仮想ポイント制を採用した主催者主導型のリアルタイムオークションシステムです。WebSocketによる双方向通信で、入札状況や価格変動を即座に反映します。

### 主な特徴

- **ポイント制**: 仮想ポイントで入札（実際の金銭取引なし）
- **主催者主導**: 主催者が開始価格と次の入札価格を決定・開示
- **タイマーレス**: 制限時間なし、主催者の判断で終了
- **リアルタイム**: WebSocketによる瞬時の状態同期
- **3ロール体制**: システム管理者・主催者・入札者の明確な権限分離

## 技術スタック

### バックエンド
- **Go 1.21+** - 主要言語
- **Gin** - RESTful APIフレームワーク
- **Gorilla WebSocket** - WebSocket通信
- **GORM** - PostgreSQL ORM
- **go-redis** - Redisクライアント
- **golang-jwt/jwt** - JWT認証
- **Air** - ホットリロード（開発環境）

### データベース・キャッシュ
- **PostgreSQL 15** - メインデータベース
- **Redis 7** - キャッシュ・Pub/Sub

### フロントエンド
- **Vue.js 3** (Composition API)
- **Vite** - ビルドツール・開発サーバー
- **Pinia** - 状態管理
- **Shadcn Vue** - UIコンポーネント
- **Tailwind CSS** - CSSフレームワーク
- **Axios** - HTTP通信

### インフラ
- **Docker** & **Docker Compose** - コンテナ環境
- **Nginx** - リバースプロキシ
- **golang-migrate** - DBマイグレーション

## クイックスタート

### 前提条件

- Docker Desktop (macOS/Windows) または Docker Engine + Docker Compose (Linux)
- Make (オプションだが推奨)

### セットアップ

1. **リポジトリのクローン**

```bash
git clone https://github.com/tsutsumi389/real-time-auction.git
cd real-time-auction
```

2. **環境変数の設定**

```bash
cp .env.example .env
# 必要に応じて .env を編集
```

3. **Docker環境の起動**

```bash
make up
# または
docker-compose up -d
```

4. **アクセス**

- **フロントエンド**: http://localhost
- **REST API**: http://localhost/api
- **WebSocket**: ws://localhost/ws

### 開発コマンド

```bash
# サービス起動・停止
make up              # 全サービス起動
make down            # 全サービス停止
make restart         # 全サービス再起動
make clean           # データベース・ボリュームを含めて完全削除

# ログ確認
make logs            # 全サービスのログ
make logs service=api  # 特定サービスのログ (api, ws, frontend, postgres, redis)

# ステータス確認
make ps              # コンテナ一覧

# サービスシェルアクセス
make shell-api       # REST APIサーバーコンテナ
make shell-ws        # WebSocketサーバーコンテナ
make shell-postgres  # PostgreSQL (psql)
make shell-redis     # Redis CLI

# データベースマイグレーション
make db-migrate      # マイグレーション適用
make db-migrate-down # 1つ前にロールバック
make db-status       # マイグレーション状態確認
make db-create-migration name=description  # 新規マイグレーション作成

# ヘルプ
make help
```

### ビルドコマンド

コンテナ内でのビルドが必要な場合:

```bash
# REST APIサーバーのビルド
docker compose exec api go build -o /app/bin/api ./cmd/api

# WebSocketサーバーのビルド
docker compose exec ws go build -o /app/bin/ws ./cmd/ws

# フロントエンドの本番ビルド
docker compose exec frontend npm run build
```

## プロジェクト構造

```
real-time-auction/
├── backend/                  # Go バックエンド
│   ├── cmd/                  # エントリーポイント
│   │   ├── api/              # REST APIサーバー
│   │   └── ws/               # WebSocketサーバー
│   ├── internal/             # アプリケーションロジック
│   │   ├── domain/           # ドメインモデル
│   │   ├── usecase/          # ビジネスロジック
│   │   ├── adapter/          # インフラ実装
│   │   └── handler/          # HTTPハンドラ
│   ├── pkg/                  # 共通パッケージ
│   └── migrations/           # DBマイグレーション
├── frontend/                 # Vue.js 3 フロントエンド
│   ├── src/
│   │   ├── components/       # Vueコンポーネント
│   │   ├── stores/           # Pinia状態管理
│   │   ├── views/            # 画面コンポーネント
│   │   ├── router/           # ルーティング
│   │   └── api/              # API通信
│   ├── public/               # 静的ファイル
│   └── package.json
├── nginx/                    # Nginx設定
│   └── nginx.conf            # リバースプロキシ設定
├── docs/                     # ドキュメント
│   ├── architecture.md       # アーキテクチャ設計書
│   ├── database_definition.md # DB定義書
│   ├── screen_list.md        # 画面一覧
│   ├── rule/                 # 開発ガイドライン
│   └── plan/                 # 実装計画
├── scripts/                  # ユーティリティスクリプト
├── docker-compose.yml        # Docker構成
├── Makefile                  # 開発コマンド
├── CLAUDE.md                 # Claude Code用ガイド
└── README.md                 # このファイル
```

## ドキュメント

### 設計書
- [アーキテクチャ設計書](docs/architecture.md) - システム全体の設計と技術選定
- [データベース定義書](docs/database_definition.md) - テーブル設計、ER図、最適化戦略
- [画面一覧](docs/screen_list.md) - 各画面の仕様とユーザーフロー

### 開発ガイドライン
- [CLAUDE.md](CLAUDE.md) - Claude Code用プロジェクト概要と開発コマンド
- [バックエンドガイドライン](docs/rule/backend.md) - Go開発のベストプラクティス
- [フロントエンドガイドライン](docs/rule/frontend.md) - Vue.js開発のベストプラクティス
- [実装計画ガイドライン](docs/rule/planning.md) - 実装計画の書き方

### API仕様（準備中）
- REST API エンドポイント仕様
- WebSocketイベント定義

## ライセンス

MIT License

## 作者

tsutsumi389
