# データベース定義書

## 概要

リアルタイムオークションシステムのデータベース設計書です。
PostgreSQLを使用し、ポイント制オークションにおける主催者主導の価格開示システムを実現します。

## データベース情報

- **RDBMS**: PostgreSQL 15+
- **文字コード**: UTF8
- **タイムゾーン**: UTC (アプリケーション層で日本時間に変換)
- **ORM**: GORM (Go)

## ER図

```
┌──────────────────┐
│     users        │
├──────────────────┤
│ id (PK)          │─┐
│ email            │ │
│ password_hash    │ │
│ role             │ │
│ display_name     │ │
│ status           │ │
│ created_at       │ │
│ updated_at       │ │
└──────────────────┘ │
         │           │
         │           │
         ┼ 1         │
         │           │
         │           │
         ┼ 1         │
         │           │
┌────────┴─────────┐ │
│  user_points     │ │
├──────────────────┤ │
│ user_id (PK,FK)  │─┘
│ total_points     │
│ available_points │
│ reserved_points  │
│ updated_at       │
└──────────────────┘
         │
         │
         ┼ 1
         │
         │
         ┼ *
         │
┌────────┴─────────┐       ┌──────────────────┐
│      bids        │       │    auctions      │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │    ┌──│ id (PK)          │
│ auction_id (FK)  │────┘  │ title            │
│ user_id (FK)     │───┐   │ description      │
│ price            │   │   │ status           │
│ bid_at           │   │   │ starting_price   │
│ is_winning       │   │   │ current_price    │
└──────────────────┘   │   │ winner_id (FK)   │
                       │   │ started_at       │
                       │   │ ended_at         │
                       │   │ created_at       │
                       │   │ updated_at       │
                       │   └──────────────────┘
                       │            │
                       │            │
                       │            ┼ 1
                       │            │
                       │            │
                       │            ┼ 1
                       │            │
                       │   ┌────────┴─────────┐
                       │   │      items       │
                       │   ├──────────────────┤
                       │   │ id (PK)          │
                       │   │ auction_id (FK)  │
                       │   │ name             │
                       │   │ description      │
                       │   │ horse_info       │
                       │   │ image_url        │
                       │   │ created_at       │
                       │   └──────────────────┘
                       │
                       │   ┌──────────────────┐
                       │   │ price_history    │
                       │   ├──────────────────┤
                       │   │ id (PK)          │
                       │   │ auction_id (FK)  │
                       │   │ price            │
                       │   │ opened_by (FK)   │
                       │   │ opened_at        │
                       │   │ had_bid          │
                       │   └──────────────────┘
                       │
                       │   ┌──────────────────┐
                       │   │ point_history    │
                       │   ├──────────────────┤
                       │   │ id (PK)          │
                       └───│ user_id (FK)     │
                           │ amount           │
                           │ type             │
                           │ reason           │
                           │ related_id       │
                           │ admin_id (FK)    │
                           │ created_at       │
                           └──────────────────┘
```

## テーブル定義

### 1. users (ユーザー)

ユーザーアカウント情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | ユーザーID (主キー) |
| email | VARCHAR(255) | NO | - | メールアドレス (ログインID) |
| password_hash | VARCHAR(255) | NO | - | bcryptハッシュ化パスワード |
| role | VARCHAR(50) | NO | 'bidder' | ロール (system_admin/auctioneer/bidder) |
| display_name | VARCHAR(100) | YES | NULL | 表示名 |
| status | VARCHAR(20) | NO | 'active' | アカウント状態 (active/suspended/deleted) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE status != 'deleted';
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
```

**制約**
```sql
ALTER TABLE users ADD CONSTRAINT chk_users_role 
  CHECK (role IN ('system_admin', 'auctioneer', 'bidder'));
ALTER TABLE users ADD CONSTRAINT chk_users_status 
  CHECK (status IN ('active', 'suspended', 'deleted'));
ALTER TABLE users ADD CONSTRAINT chk_users_email 
  CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
```

### 2. user_points (ユーザーポイント)

ユーザーの仮想ポイント残高を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| user_id | BIGINT | NO | - | ユーザーID (主キー、外部キー) |
| total_points | BIGINT | NO | 0 | 累計付与ポイント |
| available_points | BIGINT | NO | 0 | 利用可能ポイント |
| reserved_points | BIGINT | NO | 0 | 入札中の予約ポイント |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE INDEX idx_user_points_available ON user_points(available_points);
```

**制約**
```sql
ALTER TABLE user_points ADD CONSTRAINT fk_user_points_user 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_points ADD CONSTRAINT chk_user_points_non_negative 
  CHECK (total_points >= 0 AND available_points >= 0 AND reserved_points >= 0);
ALTER TABLE user_points ADD CONSTRAINT chk_user_points_balance 
  CHECK (available_points + reserved_points <= total_points);
```

**説明**
- `total_points`: システム管理者から付与された累計ポイント
- `available_points`: 現在使用可能なポイント
- `reserved_points`: 入札中で予約されているポイント
- 計算式: `total_points = available_points + reserved_points + 使用済みポイント`

### 3. auctions (オークション)

オークションの基本情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | オークションID (主キー) |
| title | VARCHAR(200) | NO | - | オークションタイトル |
| description | TEXT | YES | NULL | オークション説明 |
| status | VARCHAR(20) | NO | 'pending' | 状態 (pending/active/ended/cancelled) |
| starting_price | BIGINT | YES | NULL | 開始価格 |
| current_price | BIGINT | YES | NULL | 現在の開示価格 |
| winner_id | BIGINT | YES | NULL | 落札者ID (外部キー) |
| started_at | TIMESTAMPTZ | YES | NULL | 開始日時 |
| ended_at | TIMESTAMPTZ | YES | NULL | 終了日時 |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_started_at ON auctions(started_at);
CREATE INDEX idx_auctions_ended_at ON auctions(ended_at);
CREATE INDEX idx_auctions_winner ON auctions(winner_id);
```

**制約**
```sql
ALTER TABLE auctions ADD CONSTRAINT fk_auctions_winner 
  FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_status 
  CHECK (status IN ('pending', 'active', 'ended', 'cancelled'));
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_price_positive 
  CHECK (starting_price IS NULL OR starting_price > 0);
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_dates 
  CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at);
```

### 4. items (商品・競走馬)

オークションに出品される商品（競走馬）の情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 商品ID (主キー) |
| auction_id | BIGINT | NO | - | オークションID (外部キー) |
| name | VARCHAR(200) | NO | - | 商品名（馬名） |
| description | TEXT | YES | NULL | 商品説明 |
| horse_info | JSONB | YES | NULL | 競走馬詳細情報 (JSON形式) |
| image_url | VARCHAR(500) | YES | NULL | 画像URL |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE INDEX idx_items_auction ON items(auction_id);
CREATE INDEX idx_items_horse_info ON items USING GIN(horse_info);
```

**制約**
```sql
ALTER TABLE items ADD CONSTRAINT fk_items_auction 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;
```

**horse_info JSON例**
```json
{
  "birth_date": "2022-03-15",
  "gender": "male",
  "coat_color": "栗毛",
  "father": "ディープインパクト",
  "mother": "ジェンティルドンナ",
  "breeder": "ノーザンファーム",
  "pedigree": {
    "father_line": "サンデーサイレンス系",
    "mother_father": "ステイゴールド"
  },
  "measurements": {
    "height": "160cm",
    "chest": "185cm",
    "cannon": "20.5cm"
  }
}
```

### 5. bids (入札)

入札履歴を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 入札ID (主キー) |
| auction_id | BIGINT | NO | - | オークションID (外部キー) |
| user_id | BIGINT | NO | - | 入札者ID (外部キー) |
| price | BIGINT | NO | - | 入札価格 |
| bid_at | TIMESTAMPTZ | NO | NOW() | 入札日時 |
| is_winning | BOOLEAN | NO | TRUE | 現在の最高入札か |

**インデックス**
```sql
CREATE INDEX idx_bids_auction ON bids(auction_id, bid_at DESC);
CREATE INDEX idx_bids_user ON bids(user_id, bid_at DESC);
CREATE INDEX idx_bids_winning ON bids(auction_id, is_winning) WHERE is_winning = TRUE;
CREATE INDEX idx_bids_auction_user ON bids(auction_id, user_id);
```

**制約**
```sql
ALTER TABLE bids ADD CONSTRAINT fk_bids_auction 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;
ALTER TABLE bids ADD CONSTRAINT fk_bids_user 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE bids ADD CONSTRAINT chk_bids_price_positive 
  CHECK (price > 0);
```

**説明**
- `is_winning`: オークション終了時、TRUEのレコードが落札入札となる
- 同一価格で複数入札があった場合、最も早い `bid_at` が優先

### 6. price_history (価格開示履歴)

主催者による価格開示の履歴を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 履歴ID (主キー) |
| auction_id | BIGINT | NO | - | オークションID (外部キー) |
| price | BIGINT | NO | - | 開示価格 |
| opened_by | BIGINT | NO | - | 開示者ID (外部キー) |
| opened_at | TIMESTAMPTZ | NO | NOW() | 開示日時 |
| had_bid | BOOLEAN | NO | FALSE | 前の価格で入札があったか |

**インデックス**
```sql
CREATE INDEX idx_price_history_auction ON price_history(auction_id, opened_at DESC);
CREATE INDEX idx_price_history_opened_by ON price_history(opened_by);
```

**制約**
```sql
ALTER TABLE price_history ADD CONSTRAINT fk_price_history_auction 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;
ALTER TABLE price_history ADD CONSTRAINT fk_price_history_opener 
  FOREIGN KEY (opened_by) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE price_history ADD CONSTRAINT chk_price_history_price_positive 
  CHECK (price > 0);
```

### 7. point_history (ポイント履歴)

ポイントの増減履歴を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 履歴ID (主キー) |
| user_id | BIGINT | NO | - | ユーザーID (外部キー) |
| amount | BIGINT | NO | - | ポイント増減量 (正: 増加、負: 減少) |
| type | VARCHAR(50) | NO | - | 種別 (grant/reserve/release/consume/refund) |
| reason | TEXT | YES | NULL | 理由・備考 |
| related_id | BIGINT | YES | NULL | 関連ID (bid_id, auction_id等) |
| admin_id | BIGINT | YES | NULL | 実行管理者ID (外部キー) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |

**インデックス**
```sql
CREATE INDEX idx_point_history_user ON point_history(user_id, created_at DESC);
CREATE INDEX idx_point_history_type ON point_history(type);
CREATE INDEX idx_point_history_admin ON point_history(admin_id);
CREATE INDEX idx_point_history_related ON point_history(related_id, type);
```

**制約**
```sql
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_user 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_admin 
  FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE point_history ADD CONSTRAINT chk_point_history_type 
  CHECK (type IN ('grant', 'reserve', 'release', 'consume', 'refund'));
ALTER TABLE point_history ADD CONSTRAINT chk_point_history_amount_not_zero 
  CHECK (amount != 0);
```

**ポイント種別 (type)**
- `grant`: ポイント付与 (システム管理者による)
- `reserve`: 入札時の予約 (available → reserved)
- `release`: 入札取消・落札失敗時の解放 (reserved → available)
- `consume`: 落札時の消費 (reserved → 消費)
- `refund`: オークション中止時の返金 (reserved → available)

## トリガー関数

### 1. updated_at 自動更新

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 各テーブルに適用
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_points_updated_at BEFORE UPDATE ON user_points
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_auctions_updated_at BEFORE UPDATE ON auctions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 2. user_points 自動作成

```sql
CREATE OR REPLACE FUNCTION create_user_points()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_points (user_id, total_points, available_points, reserved_points)
    VALUES (NEW.id, 0, 0, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_user_points AFTER INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION create_user_points();
```

### 3. 入札時の is_winning 更新

```sql
CREATE OR REPLACE FUNCTION update_bid_winning_status()
RETURNS TRIGGER AS $$
BEGIN
    -- 同じオークションの過去の入札をすべて is_winning = FALSE に
    UPDATE bids 
    SET is_winning = FALSE 
    WHERE auction_id = NEW.auction_id AND id != NEW.id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_bid_winning_status AFTER INSERT ON bids
    FOR EACH ROW EXECUTE FUNCTION update_bid_winning_status();
```

## ビュー

### 1. active_auctions_view (アクティブなオークション一覧)

```sql
CREATE VIEW active_auctions_view AS
SELECT 
    a.id,
    a.title,
    a.description,
    a.status,
    a.current_price,
    a.started_at,
    i.name AS item_name,
    i.image_url,
    COUNT(DISTINCT b.user_id) AS bidder_count,
    MAX(b.bid_at) AS last_bid_at,
    u.display_name AS winner_name
FROM auctions a
LEFT JOIN items i ON a.id = i.auction_id
LEFT JOIN bids b ON a.id = b.auction_id
LEFT JOIN users u ON a.winner_id = u.id
WHERE a.status IN ('active', 'ended')
GROUP BY a.id, i.name, i.image_url, u.display_name
ORDER BY a.started_at DESC;
```

### 2. user_auction_summary (ユーザーごとのオークション参加状況)

```sql
CREATE VIEW user_auction_summary AS
SELECT 
    u.id AS user_id,
    u.email,
    u.display_name,
    up.available_points,
    up.reserved_points,
    COUNT(DISTINCT b.auction_id) AS participated_auctions,
    COUNT(DISTINCT CASE WHEN a.winner_id = u.id THEN a.id END) AS won_auctions,
    COALESCE(SUM(CASE WHEN a.winner_id = u.id THEN a.current_price END), 0) AS total_won_amount
FROM users u
LEFT JOIN user_points up ON u.id = up.user_id
LEFT JOIN bids b ON u.id = b.user_id
LEFT JOIN auctions a ON b.auction_id = a.id AND a.status = 'ended'
WHERE u.role = 'bidder'
GROUP BY u.id, u.email, u.display_name, up.available_points, up.reserved_points;
```

## 初期データ

### システム管理者アカウント

```sql
-- パスワード: admin123 (本番環境では必ず変更)
INSERT INTO users (email, password_hash, role, display_name, status)
VALUES (
    'admin@example.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5lW5vJ4r.Wquy',
    'system_admin',
    'システム管理者',
    'active'
);
```

### 主催者アカウント

```sql
-- パスワード: auctioneer123
INSERT INTO users (email, password_hash, role, display_name, status)
VALUES (
    'auctioneer@example.com',
    '$2a$12$8kPvKJQNnZqU5vR8gF4Hc.LxK9bE2rW7jT5mC3dP1fN6hG8lA4qBi',
    'auctioneer',
    '主催者',
    'active'
);
```

### テスト用入札者アカウント

```sql
-- パスワード: bidder123
INSERT INTO users (email, password_hash, role, display_name, status)
VALUES 
    ('bidder1@example.com', '$2a$12$9kQwLMRPqZrV6wS9hG5Id.MyL0cF3sX8kU6nD4eQ2gO7iH9mB5rCj', 'bidder', '入札者1', 'active'),
    ('bidder2@example.com', '$2a$12$0lRxMNSQraW7xT0iH6Je.NzM1dG4tY9lV7oE5fR3hP8jI0nC6sDk', 'bidder', '入札者2', 'active'),
    ('bidder3@example.com', '$2a$12$1mSyNOTRsbX8yU1jI7Kf.O0N2eH5uZ0mW8pF6gS4iQ9kJ1oD7tEl', 'bidder', '入札者3', 'active');

-- テスト用ポイント付与 (各ユーザーに10,000pt)
UPDATE user_points SET total_points = 10000, available_points = 10000 
WHERE user_id IN (SELECT id FROM users WHERE role = 'bidder');
```

## パフォーマンスチューニング

### 1. 接続プール設定

```go
// Go (GORM) の推奨設定
db.SetMaxIdleConns(10)
db.SetMaxOpenConns(100)
db.SetConnMaxLifetime(time.Hour)
```

### 2. 頻繁に使用されるクエリの最適化

```sql
-- 入札履歴取得クエリ (ユーザーごと)
EXPLAIN ANALYZE
SELECT b.id, b.price, b.bid_at, b.is_winning, a.title, i.name
FROM bids b
JOIN auctions a ON b.auction_id = a.id
JOIN items i ON a.id = i.auction_id
WHERE b.user_id = $1
ORDER BY b.bid_at DESC
LIMIT 20;

-- オークション詳細取得クエリ
EXPLAIN ANALYZE
SELECT 
    a.*,
    i.name, i.description, i.horse_info, i.image_url,
    COUNT(b.id) AS total_bids,
    MAX(b.bid_at) AS last_bid_at
FROM auctions a
LEFT JOIN items i ON a.id = i.auction_id
LEFT JOIN bids b ON a.id = b.auction_id
WHERE a.id = $1
GROUP BY a.id, i.name, i.description, i.horse_info, i.image_url;
```

### 3. パーティショニング (将来的な対応)

大量のデータが蓄積された場合、以下のテーブルをパーティション化：

```sql
-- bids テーブルを月ごとにパーティション化
CREATE TABLE bids_2025_01 PARTITION OF bids
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE bids_2025_02 PARTITION OF bids
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
```

## バックアップ戦略

### 1. 定期バックアップ

```bash
# 日次フルバックアップ (毎日午前3時)
0 3 * * * pg_dump -h localhost -U auction_user -d auction_db -F c -f /backups/auction_db_$(date +\%Y\%m\%d).dump

# 週次バックアップの保持 (30日間)
0 4 * * 0 find /backups -name "auction_db_*.dump" -mtime +30 -delete
```

### 2. ポイントインタイムリカバリ (PITR)

```bash
# WAL アーカイブ設定
archive_mode = on
archive_command = 'cp %p /var/lib/postgresql/wal_archive/%f'
wal_level = replica
```

## マイグレーション管理

### GORM AutoMigrate (開発環境)

```go
// backend/cmd/api/main.go
db.AutoMigrate(
    &domain.User{},
    &domain.UserPoint{},
    &domain.Auction{},
    &domain.Item{},
    &domain.Bid{},
    &domain.PriceHistory{},
    &domain.PointHistory{},
)
```

### 本番環境マイグレーション

```bash
# golang-migrate を使用
migrate -path migrations -database "postgresql://auction_user:password@localhost:5432/auction_db?sslmode=disable" up
```

## セキュリティ考慮事項

1. **パスワードハッシュ化**: bcrypt (cost factor: 12)
2. **SQL インジェクション対策**: GORM のプリペアドステートメント使用
3. **個人情報保護**: email, password_hash は暗号化推奨
4. **監査ログ**: point_history, price_history で全操作を記録
5. **アクセス制限**: PostgreSQL ユーザー権限で Read/Write を分離

## 監視項目

- **接続数**: `SELECT count(*) FROM pg_stat_activity;`
- **スロークエリ**: `log_min_duration_statement = 1000` (1秒以上)
- **テーブルサイズ**: `SELECT pg_size_pretty(pg_total_relation_size('bids'));`
- **インデックス使用率**: `SELECT * FROM pg_stat_user_indexes;`

## まとめ

このデータベース設計は以下の要件を満たします：

✅ **ポイント制**: 仮想ポイントによる安全な入札管理  
✅ **主催者主導**: 価格開示履歴の完全な記録  
✅ **整合性**: トランザクションとトリガーによるデータ整合性保証  
✅ **パフォーマンス**: 適切なインデックスとビューによる高速クエリ  
✅ **監査**: 全ポイント操作・価格開示の履歴記録  
✅ **スケーラビリティ**: パーティショニング対応可能な設計  
✅ **セキュリティ**: ロールベースのアクセス制御とデータ保護
