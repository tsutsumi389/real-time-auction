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
- **Go 1.21+**
- **Gin** - Webフレームワーク
- **Gorilla WebSocket** - WebSocket通信
- **GORM** - ORM
- **PostgreSQL 15** - データベース
- **Redis 7** - キャッシュ・Pub/Sub

### フロントエンド
- **Vue.js 3** (Composition API)
- **Vite** - ビルドツール
- **Pinia** - 状態管理
- **Axios** - HTTP通信

### インフラ
- **Docker** & **Docker Compose**
- **Nginx** - リバースプロキシ

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
# 起動
make up

# 停止
make down

# ログ確認
make logs

# 特定サービスのログ
make logs service=api

# ステータス確認
make ps

# 再起動
make restart

# データベース・ボリュームを含めて完全削除
make clean

# ヘルプ
make help
```

## プロジェクト構造

```
real-time-auction/
├── backend/              # Go バックエンド
│   ├── cmd/              # エントリーポイント
│   │   ├── api/          # REST APIサーバー
│   │   └── ws/           # WebSocketサーバー
│   ├── internal/         # アプリケーションロジック
│   ├── pkg/              # 共通パッケージ
│   └── migrations/       # DBマイグレーション
├── frontend/             # Vue.js フロントエンド
│   ├── src/              # ソースコード
│   ├── public/           # 静的ファイル
│   └── package.json      # 依存関係
├── nginx/                # Nginx設定
│   └── nginx.conf        # リバースプロキシ設定
├── docs/                 # ドキュメント
│   └── architecture.md   # アーキテクチャ設計書
├── scripts/              # ユーティリティスクリプト
├── docker-compose.yml    # Docker構成
├── Makefile              # 開発コマンド
└── README.md             # このファイル
```

## ドキュメント

- [アーキテクチャ設計書](docs/architecture.md) - システム全体の設計と技術選定
- [API仕様書](backend/docs/api.md) - REST API エンドポイント仕様（準備中）
- [WebSocketイベント仕様](backend/docs/websocket.md) - WebSocketイベント定義（準備中）

## 開発ロードマップ

### Phase 1: Webアプリケーション（現在）

- [x] プロジェクト基盤構築
- [x] Docker開発環境セットアップ
- [ ] データベーススキーマ設計
- [ ] 認証・認可システム
- [ ] REST API実装
- [ ] WebSocket実装
- [ ] Vue.jsフロントエンド実装
- [ ] システム管理機能
- [ ] 主催者管理機能
- [ ] 入札者機能

### Phase 2: iOSネイティブアプリ（将来）

- [ ] Swift/SwiftUI アプリケーション
- [ ] APNsプッシュ通知
- [ ] オフライン対応
- [ ] Face ID/Touch ID認証

## ライセンス

MIT License

## 作者

tsutsumi389
