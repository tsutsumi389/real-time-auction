# Database Migrations

このディレクトリには、PostgreSQLデータベースのマイグレーションファイルが含まれています。

## マイグレーションファイル一覧

### 001_create_initial_tables
- テーブル作成: `bidders`, `admins`, `bidder_points`, `auctions`, `items`, `item_media`, `bids`, `price_history`, `point_history`
- インデックス、制約、外部キー制約の設定
- UUID拡張機能の有効化

### 002_create_triggers_and_functions
- トリガー関数:
  - `update_updated_at_column()`: `updated_at`カラムの自動更新
  - `create_bidder_points()`: 入札者作成時にポイントレコードを自動作成
  - `update_bid_winning_status()`: 入札時に`is_winning`フラグを更新
  - `record_point_history()`: ポイント変動時に履歴を自動記録

### 003_create_views
- ビュー作成:
  - `active_auctions_view`: アクティブなオークション一覧
  - `bidder_auction_summary`: 入札者ごとのオークション参加状況
  - `point_history_detailed`: ポイント履歴詳細
  - `bidder_point_balance_view`: 入札者ポイント残高
  - `point_transactions_by_auction`: オークション別ポイント取引

### 004_insert_initial_data
- 初期データ挿入:
  - システム管理者アカウント (`admin@example.com` / パスワード: `admin123`)
  - 主催者アカウント (`auctioneer@example.com` / パスワード: `auctioneer123`)
  - テスト用入札者アカウント3件 (各10,000ポイント付与)

⚠️ **本番環境では必ずパスワードを変更してください**

## マイグレーション実行方法

### golang-migrate を使用する場合

```bash
# マイグレーション実行（環境変数使用を推奨）
export DB_PASSWORD="your_secure_password"
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  up

# 特定のバージョンまで実行
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  goto 2

# 1つ戻す
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  down 1

# 全て戻す
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  down

# バージョン確認
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  version
```

### Docker環境での実行

```bash
# コンテナ内で実行
docker-compose exec api sh -c "migrate -path /app/migrations -database postgresql://\$DB_USER:\$DB_PASSWORD@\$DB_HOST:\$DB_PORT/\$DB_NAME?sslmode=disable up"

# または、PostgreSQLコンテナから直接実行
docker-compose exec postgres psql -U auction_user -d auction_db -f /docker-entrypoint-initdb.d/001_create_initial_tables.up.sql
```

### GORMのAutoMigrate (開発環境のみ推奨)

`cmd/api/main.go` や `cmd/ws/main.go` で以下のように実行:

```go
package main

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "your-module/internal/domain"
)

func main() {
    dsn := fmt.Sprintf("host=localhost user=auction_user password=%s dbname=auction_db port=5432 sslmode=disable", os.Getenv("DB_PASSWORD"))
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }

    // AutoMigrate (開発環境のみ)
    db.AutoMigrate(
        &domain.Bidder{},
        &domain.Admin{},
        &domain.BidderPoint{},
        &domain.Auction{},
        &domain.Item{},
        &domain.ItemMedia{},
        &domain.Bid{},
        &domain.PriceHistory{},
        &domain.PointHistory{},
    )
}
```

⚠️ **本番環境では`golang-migrate`等のマイグレーションツールを使用してください**

## ファイル命名規則

- `<version>_<description>.up.sql`: マイグレーション適用用
- `<version>_<description>.down.sql`: ロールバック用

バージョン番号は3桁の連番 (001, 002, 003...)

## データベース構造

詳細なデータベース設計については、`docs/database_definition.md` を参照してください。

### 主要テーブル

- **bidders**: 入札者アカウント (UUID主キー)
- **admins**: 管理者・主催者アカウント
- **bidder_points**: 入札者の仮想ポイント残高
- **auctions**: オークション情報
- **items**: 商品情報 (JSONBメタデータ対応)
- **item_media**: 商品の画像・動画
- **bids**: 入札履歴
- **price_history**: 価格開示履歴 (主催者主導型)
- **point_history**: ポイント変動履歴 (監査証跡)

## トラブルシューティング

### マイグレーションが失敗する場合

```bash
# マイグレーション状態確認
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  version

# 強制的にバージョンを設定
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  force <version>

# schema_migrations テーブルを確認
docker-compose exec postgres psql -U auction_user -d auction_db -c "SELECT * FROM schema_migrations;"
```

### データベースをリセットしたい場合

```bash
# 全テーブル削除
docker-compose exec postgres psql -U auction_user -d auction_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# マイグレーション再実行
migrate -path backend/migrations \
  -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" \
  up
```

## セキュリティ注意事項

1. **本番環境では必ず初期パスワードを変更**
   - システム管理者: `admin@example.com`
   - 主催者: `auctioneer@example.com`

2. **データベース接続情報を環境変数で管理**
   ```bash
   export DB_USER=auction_user
   export DB_PASSWORD=your_secure_password
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_NAME=auction_db
   ```

3. **本番環境ではSSL接続を有効化**
   ```
   sslmode=require
   ```

## 関連ドキュメント

- [データベース定義書](../../docs/database_definition.md)
- [アーキテクチャ概要](../../docs/architecture.md)
- [プロジェクトREADME](../../README.md)
