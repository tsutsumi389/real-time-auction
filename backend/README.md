# Backend

Go言語によるバックエンドサーバー（REST API + WebSocket）

## 構造

```
backend/
├── cmd/
│   ├── api/              # REST APIサーバーエントリーポイント
│   └── ws/               # WebSocketサーバーエントリーポイント
├── internal/             # プライベートアプリケーションコード
│   ├── domain/           # ドメインモデル
│   ├── repository/       # データアクセス層
│   ├── service/          # ビジネスロジック
│   ├── handler/          # HTTPハンドラー
│   ├── ws/               # WebSocketハンドラー
│   └── middleware/       # ミドルウェア
├── pkg/                  # 公開パッケージ
│   ├── config/           # 設定管理
│   ├── logger/           # ロギング
│   └── validator/        # カスタムバリデーター
├── migrations/           # データベースマイグレーション
├── docs/                 # API仕様書
├── .air.toml             # Air設定（APIサーバー）
├── .air.ws.toml          # Air設定（WebSocketサーバー）
├── Dockerfile            # 本番用Dockerfile
├── Dockerfile.dev        # 開発用Dockerfile
├── go.mod                # Go modules
└── go.sum                # Go modules checksum
```

## 開発

### ローカル開発（Dockerなし）

```bash
# 依存関係のインストール
go mod download

# APIサーバー起動
go run cmd/api/main.go

# WebSocketサーバー起動
go run cmd/ws/main.go
```

### Docker開発環境

プロジェクトルートから:

```bash
# 起動
make up

# APIログ確認
make logs-api

# WSログ確認
make logs-ws
```

## 実装予定

- [ ] データベース接続
- [ ] JWT認証
- [ ] ユーザー管理API
- [ ] オークション管理API
- [ ] WebSocket実装
- [ ] Redis Pub/Sub
- [ ] 入札ロジック
