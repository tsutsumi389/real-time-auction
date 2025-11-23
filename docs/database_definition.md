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
┌──────────────────┐         ┌──────────────────┐
│     bidders      │         │     admins       │
├──────────────────┤         ├──────────────────┤
│ id (UUID,PK)     │─┐       │ id (PK)          │
│ email            │ │       │ email            │
│ password_hash    │ │       │ password_hash    │
│ display_name     │ │       │ role             │
│ status           │ │       │ display_name     │
│ created_at       │ │       │ status           │
│ updated_at       │ │       │ created_at       │
└──────────────────┘ │       │ updated_at       │
         │           │       └──────────────────┘
         │           │                │
         ┼ 1         │                │
         │           │                │
         │           │                │
         ┼ 1         │                │
         │           │                │
┌────────┴─────────┐ │                │
│  bidder_points   │ │                │
├──────────────────┤ │                │
│ bidder_id (PK,FK)│─┘                │
│ total_points     │                  │
│ available_points │                  │
│ reserved_points  │                  │
│ updated_at       │                  │
└──────────────────┘                  │
         │                            │
         │                            │
         ┼ 1                          │
         │                            │
         │                            │
         ┼ *                          │
         │                            │
┌────────┴─────────┐       ┌──────────────────┐
│      bids        │       │    auctions      │
├──────────────────┤       ├──────────────────┤
│ id (PK)          │    ┌──│ id (PK)          │
│ auction_id (FK)  │────┘  │ title            │
│ bidder_id (FK)   │───┐   │ description      │
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
                       │   │ metadata         │
                       │   │ created_at       │
                       │   │ updated_at       │
                       │   └────────┬─────────┘
                       │            │
                       │            ┼ 1
                       │            │
                       │            │
                       │            ┼ *
                       │            │
                       │   ┌────────┴─────────┐
                       │   │  item_media      │
                       │   ├──────────────────┤
                       │   │ id (PK)          │
                       │   │ item_id (FK)     │
                       │   │ media_type       │
                       │   │ url              │
                       │   │ thumbnail_url    │
                       │   │ display_order    │
                       │   │ created_at       │
                       │   └──────────────────┘
                       │
                       │   ┌──────────────────┐
                       │   │ price_history    │
                       │   ├──────────────────┤
                       │   │ id (PK)          │
                       │   │ auction_id (FK)  │
                       │   │ price            │
                       │   │ opened_by (FK)   │─┐
                       │   │ opened_at        │ │
                       │   │ had_bid          │ │
                       │   └──────────────────┘ │
                       │                        │
                       │   ┌────────────────────┴┐
                       │   │ point_history       │
                       │   ├─────────────────────┤
                       │   │ id (PK)             │
                       └───│ bidder_id (FK)      │
                           │ amount              │
                           │ type                │
                           │ reason              │
                           │ related_id          │
                           │ admin_id (FK)       │──┘
                           │ created_at          │
                           └─────────────────────┘
```

## テーブル定義

### 1. bidders (入札者)

入札者アカウント情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | UUID | NO | gen_random_uuid() | 入札者ID (主キー) |
| email | VARCHAR(255) | NO | - | メールアドレス (ログインID) |
| password_hash | VARCHAR(255) | NO | - | bcryptハッシュ化パスワード |
| display_name | VARCHAR(100) | YES | NULL | 表示名 |
| status | VARCHAR(20) | NO | 'active' | アカウント状態 (active/suspended/deleted) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE UNIQUE INDEX idx_bidders_email ON bidders(email) WHERE status != 'deleted';
CREATE INDEX idx_bidders_status ON bidders(status);
```

**制約**
```sql
ALTER TABLE bidders ADD CONSTRAINT chk_bidders_status 
  CHECK (status IN ('active', 'suspended', 'deleted'));
ALTER TABLE bidders ADD CONSTRAINT chk_bidders_email 
  CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
```

### 2. admins (管理者・主催者)

管理者およびオークション主催者のアカウント情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 管理者ID (主キー) |
| email | VARCHAR(255) | NO | - | メールアドレス (ログインID) |
| password_hash | VARCHAR(255) | NO | - | bcryptハッシュ化パスワード |
| role | VARCHAR(50) | NO | 'auctioneer' | ロール (system_admin/auctioneer) |
| display_name | VARCHAR(100) | YES | NULL | 表示名 |
| status | VARCHAR(20) | NO | 'active' | アカウント状態 (active/suspended/deleted) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE UNIQUE INDEX idx_admins_email ON admins(email) WHERE status != 'deleted';
CREATE INDEX idx_admins_role ON admins(role);
CREATE INDEX idx_admins_status ON admins(status);
```

**制約**
```sql
ALTER TABLE admins ADD CONSTRAINT chk_admins_role 
  CHECK (role IN ('system_admin', 'auctioneer'));
ALTER TABLE admins ADD CONSTRAINT chk_admins_status 
  CHECK (status IN ('active', 'suspended', 'deleted'));
ALTER TABLE admins ADD CONSTRAINT chk_admins_email 
  CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
```

### 3. bidder_points (入札者ポイント)

入札者の仮想ポイント残高を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| bidder_id | UUID | NO | - | 入札者ID (主キー、外部キー) |
| total_points | BIGINT | NO | 0 | 累計付与ポイント |
| available_points | BIGINT | NO | 0 | 利用可能ポイント |
| reserved_points | BIGINT | NO | 0 | 入札中の予約ポイント |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE INDEX idx_bidder_points_available ON bidder_points(available_points);
```

**制約**
```sql
ALTER TABLE bidder_points ADD CONSTRAINT fk_bidder_points_bidder 
  FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE;
ALTER TABLE bidder_points ADD CONSTRAINT chk_bidder_points_non_negative 
  CHECK (total_points >= 0 AND available_points >= 0 AND reserved_points >= 0);
ALTER TABLE bidder_points ADD CONSTRAINT chk_bidder_points_balance 
  CHECK (available_points + reserved_points <= total_points);
```

**説明**
- `total_points`: システム管理者から付与された累計ポイント
- `available_points`: 現在使用可能なポイント
- `reserved_points`: 入札中で予約されているポイント
- 計算式: `total_points = available_points + reserved_points + 使用済みポイント`

### 4. auctions (オークション)

オークションの基本情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | UUID | NO | gen_random_uuid() | オークションID (主キー) |
| title | VARCHAR(200) | NO | - | オークションタイトル |
| description | TEXT | YES | NULL | オークション説明 |
| status | VARCHAR(20) | NO | 'pending' | 状態 (pending/active/ended/cancelled) |
| starting_price | BIGINT | YES | NULL | 開始価格 |
| current_price | BIGINT | YES | NULL | 現在の開示価格 |
| winner_id | UUID | YES | NULL | 落札者ID (外部キー) |
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
  FOREIGN KEY (winner_id) REFERENCES bidders(id) ON DELETE SET NULL;
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_status 
  CHECK (status IN ('pending', 'active', 'ended', 'cancelled'));
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_price_positive 
  CHECK (starting_price IS NULL OR starting_price > 0);
ALTER TABLE auctions ADD CONSTRAINT chk_auctions_dates 
  CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at);
```

### 5. items (商品)

オークションに出品される商品の基本情報を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | UUID | NO | gen_random_uuid() | 商品ID (主キー) |
| auction_id | UUID | NO | - | オークションID (外部キー) |
| name | VARCHAR(200) | NO | - | 商品名 |
| description | TEXT | YES | NULL | 商品説明 |
| metadata | JSONB | YES | NULL | 商品詳細情報 (JSON形式、自由フォーマット) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NO | NOW() | 更新日時 |

**インデックス**
```sql
CREATE INDEX idx_items_auction ON items(auction_id);
CREATE INDEX idx_items_metadata ON items USING GIN(metadata);
```

**制約**
```sql
ALTER TABLE items ADD CONSTRAINT fk_items_auction 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;
```

**metadata JSON例（競走馬の場合）**
```json
{
  "category": "horse",
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

**metadata JSON例（その他商品の場合）**
```json
{
  "category": "art",
  "artist": "山田太郎",
  "year": "2023",
  "dimensions": "100cm x 80cm",
  "material": "油彩・キャンバス",
  "condition": "良好"
}
```

### 6. item_media (商品メディア)

商品に紐づく画像・動画を管理するテーブル。1つの商品に複数のメディアを登録可能。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | メディアID (主キー) |
| item_id | UUID | NO | - | 商品ID (外部キー) |
| media_type | VARCHAR(20) | NO | - | メディア種別 (image/video) |
| url | VARCHAR(500) | NO | - | メディアURL (S3等) |
| thumbnail_url | VARCHAR(500) | YES | NULL | サムネイルURL (動画の場合) |
| display_order | INT | NO | 0 | 表示順序 (小さい順に表示) |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |

**インデックス**
```sql
CREATE INDEX idx_item_media_item ON item_media(item_id, display_order);
CREATE INDEX idx_item_media_type ON item_media(media_type);
```

**制約**
```sql
ALTER TABLE item_media ADD CONSTRAINT fk_item_media_item 
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE;
ALTER TABLE item_media ADD CONSTRAINT chk_item_media_type 
  CHECK (media_type IN ('image', 'video'));
ALTER TABLE item_media ADD CONSTRAINT chk_item_media_order_non_negative 
  CHECK (display_order >= 0);
```

**使用例**
```sql
-- 商品に3つの画像と1つの動画を登録
INSERT INTO item_media (item_id, media_type, url, display_order) VALUES
  (1, 'image', 'https://cdn.example.com/items/1/main.jpg', 0),
  (1, 'image', 'https://cdn.example.com/items/1/side.jpg', 1),
  (1, 'video', 'https://cdn.example.com/items/1/movie.mp4', 2),
  (1, 'image', 'https://cdn.example.com/items/1/back.jpg', 3);
```

### 7. bids (入札)

入札履歴を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 入札ID (主キー) |
| auction_id | UUID | NO | - | オークションID (外部キー) |
| bidder_id | UUID | NO | - | 入札者ID (外部キー) |
| price | BIGINT | NO | - | 入札価格 |
| bid_at | TIMESTAMPTZ | NO | NOW() | 入札日時 |
| is_winning | BOOLEAN | NO | TRUE | 現在の最高入札か |

**インデックス**
```sql
CREATE INDEX idx_bids_auction ON bids(auction_id, bid_at DESC);
CREATE INDEX idx_bids_bidder ON bids(bidder_id, bid_at DESC);
CREATE INDEX idx_bids_winning ON bids(auction_id, is_winning) WHERE is_winning = TRUE;
CREATE INDEX idx_bids_auction_bidder ON bids(auction_id, bidder_id);
```

**制約**
```sql
ALTER TABLE bids ADD CONSTRAINT fk_bids_auction 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;
ALTER TABLE bids ADD CONSTRAINT fk_bids_bidder 
  FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE;
ALTER TABLE bids ADD CONSTRAINT chk_bids_price_positive 
  CHECK (price > 0);
```

**説明**
- `is_winning`: オークション終了時、TRUEのレコードが落札入札となる
- 同一価格で複数入札があった場合、最も早い `bid_at` が優先

### 8. price_history (価格開示履歴)

主催者による価格開示の履歴を管理するテーブル。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 履歴ID (主キー) |
| auction_id | UUID | NO | - | オークションID (外部キー) |
| price | BIGINT | NO | - | 開示価格 |
| opened_by | BIGINT | NO | - | 開示者ID (外部キー: admins) |
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
  FOREIGN KEY (opened_by) REFERENCES admins(id) ON DELETE CASCADE;
ALTER TABLE price_history ADD CONSTRAINT chk_price_history_price_positive 
  CHECK (price > 0);
```

### 9. point_history (ポイント履歴)

ポイントの増減履歴を管理するテーブル。全てのポイント操作を記録し、監査証跡として機能。

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGSERIAL | NO | - | 履歴ID (主キー) |
| bidder_id | UUID | NO | - | 入札者ID (外部キー) |
| amount | BIGINT | NO | - | ポイント増減量 (正: 増加、負: 減少) |
| type | VARCHAR(50) | NO | - | 種別 (grant/reserve/release/consume/refund) |
| reason | TEXT | YES | NULL | 理由・備考 |
| related_auction_id | UUID | YES | NULL | 関連オークションID (外部キー) |
| related_bid_id | BIGINT | YES | NULL | 関連入札ID (外部キー) |
| admin_id | BIGINT | YES | NULL | 実行管理者ID (外部キー: admins) |
| balance_before | BIGINT | NO | - | 操作前の利用可能ポイント |
| balance_after | BIGINT | NO | - | 操作後の利用可能ポイント |
| reserved_before | BIGINT | NO | - | 操作前の予約ポイント |
| reserved_after | BIGINT | NO | - | 操作後の予約ポイント |
| total_before | BIGINT | NO | - | 操作前の累計ポイント |
| total_after | BIGINT | NO | - | 操作後の累計ポイント |
| created_at | TIMESTAMPTZ | NO | NOW() | 作成日時 |

**インデックス**
```sql
CREATE INDEX idx_point_history_bidder ON point_history(bidder_id, created_at DESC);
CREATE INDEX idx_point_history_type ON point_history(type);
CREATE INDEX idx_point_history_admin ON point_history(admin_id);
CREATE INDEX idx_point_history_auction ON point_history(related_auction_id);
CREATE INDEX idx_point_history_bid ON point_history(related_bid_id);
CREATE INDEX idx_point_history_created_at ON point_history(created_at DESC);
```

**制約**
```sql
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_bidder 
  FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE;
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_admin 
  FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE SET NULL;
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_auction 
  FOREIGN KEY (related_auction_id) REFERENCES auctions(id) ON DELETE SET NULL;
ALTER TABLE point_history ADD CONSTRAINT fk_point_history_bid 
  FOREIGN KEY (related_bid_id) REFERENCES bids(id) ON DELETE SET NULL;
ALTER TABLE point_history ADD CONSTRAINT chk_point_history_type 
  CHECK (type IN ('grant', 'reserve', 'release', 'consume', 'refund'));
ALTER TABLE point_history ADD CONSTRAINT chk_point_history_amount_not_zero 
  CHECK (amount != 0);
ALTER TABLE point_history ADD CONSTRAINT chk_point_history_balance_non_negative 
  CHECK (balance_before >= 0 AND balance_after >= 0 AND 
         reserved_before >= 0 AND reserved_after >= 0 AND
         total_before >= 0 AND total_after >= 0);
```

**ポイント種別 (type) と動作**

| 種別 | 説明 | amount | available | reserved | total | 実行者 |
|-----|------|--------|-----------|----------|-------|--------|
| `grant` | ポイント付与 | +N | +N | - | +N | system_admin |
| `reserve` | 入札時の予約 | -N | -N | +N | - | bidder (自動) |
| `release` | 落札失敗時の解放 | +N | +N | -N | - | system (自動) |
| `consume` | 落札時の消費 | -N | - | -N | - | system (自動) |
| `refund` | オークション中止時の返金 | +N | +N | -N | - | auctioneer/admin |

**履歴記録の詳細**

各ポイント操作は以下の情報を記録：
1. **変更前後の全残高**: `balance_before/after`, `reserved_before/after`, `total_before/after`
2. **関連エンティティ**: `related_auction_id`, `related_bid_id`
3. **実行者**: `admin_id` (付与・返金の場合)
4. **理由**: `reason` (自由記述)

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
CREATE TRIGGER update_bidders_updated_at BEFORE UPDATE ON bidders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_admins_updated_at BEFORE UPDATE ON admins
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bidder_points_updated_at BEFORE UPDATE ON bidder_points
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_auctions_updated_at BEFORE UPDATE ON auctions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- item_media テーブルには updated_at がないため、トリガー不要
```

### 2. bidder_points 自動作成

**注意**: このトリガーは削除されました。バックエンドのアプリケーション層で実装されます。

~~削除済み: `trigger_create_bidder_points` および `create_bidder_points()` 関数~~

### 3. 入札時の is_winning 更新

**注意**: このトリガーは削除されました。バックエンドのアプリケーション層で実装されます。

~~削除済み: `trigger_update_bid_winning_status` および `update_bid_winning_status()` 関数~~

### 4. ポイント履歴自動記録

**注意**: このトリガーは削除されました。バックエンドのアプリケーション層で実装されます。

~~削除済み: `trigger_record_point_history` および `record_point_history()` 関数~~

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
    (SELECT url FROM item_media WHERE item_id = i.id ORDER BY display_order LIMIT 1) AS main_image_url,
    COUNT(DISTINCT b.bidder_id) AS bidder_count,
    MAX(b.bid_at) AS last_bid_at,
    bd.display_name AS winner_name
FROM auctions a
LEFT JOIN items i ON a.id = i.auction_id
LEFT JOIN bids b ON a.id = b.auction_id
LEFT JOIN bidders bd ON a.winner_id = bd.id
WHERE a.status IN ('active', 'ended')
GROUP BY a.id, i.id, i.name, bd.display_name
ORDER BY a.started_at DESC;
```

### 2. bidder_auction_summary (入札者ごとのオークション参加状況)

```sql
CREATE VIEW bidder_auction_summary AS
SELECT 
    bd.id AS bidder_id,
    bd.email,
    bd.display_name,
    bp.available_points,
    bp.reserved_points,
    COUNT(DISTINCT b.auction_id) AS participated_auctions,
    COUNT(DISTINCT CASE WHEN a.winner_id = bd.id THEN a.id END) AS won_auctions,
    COALESCE(SUM(CASE WHEN a.winner_id = bd.id THEN a.current_price END), 0) AS total_won_amount
FROM bidders bd
LEFT JOIN bidder_points bp ON bd.id = bp.bidder_id
LEFT JOIN bids b ON bd.id = b.bidder_id
LEFT JOIN auctions a ON b.auction_id = a.id AND a.status = 'ended'
GROUP BY bd.id, bd.email, bd.display_name, bp.available_points, bp.reserved_points;
```

### 3. point_history_detailed (ポイント履歴詳細ビュー)

```sql
CREATE VIEW point_history_detailed AS
SELECT 
    ph.id,
    ph.bidder_id,
    bd.email AS bidder_email,
    bd.display_name AS bidder_name,
    ph.amount,
    ph.type,
    ph.reason,
    ph.balance_before,
    ph.balance_after,
    ph.reserved_before,
    ph.reserved_after,
    ph.total_before,
    ph.total_after,
    a.title AS auction_title,
    b.price AS bid_price,
    ad.display_name AS admin_name,
    ph.created_at
FROM point_history ph
JOIN bidders bd ON ph.bidder_id = bd.id
LEFT JOIN auctions a ON ph.related_auction_id = a.id
LEFT JOIN bids b ON ph.related_bid_id = b.id
LEFT JOIN admins ad ON ph.admin_id = ad.id
ORDER BY ph.created_at DESC;
```

### 4. bidder_point_balance_view (入札者ポイント残高ビュー)

```sql
CREATE VIEW bidder_point_balance_view AS
SELECT 
    bd.id AS bidder_id,
    bd.email,
    bd.display_name,
    bp.total_points,
    bp.available_points,
    bp.reserved_points,
    (bp.total_points - bp.available_points - bp.reserved_points) AS consumed_points,
    COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0) AS total_granted,
    COALESCE(SUM(CASE WHEN ph.type = 'consume' THEN ABS(ph.amount) END), 0) AS total_consumed,
    COALESCE(COUNT(CASE WHEN ph.type = 'reserve' THEN 1 END), 0) AS total_bids,
    bp.updated_at AS last_updated
FROM bidders bd
LEFT JOIN bidder_points bp ON bd.id = bp.bidder_id
LEFT JOIN point_history ph ON bd.id = ph.bidder_id
GROUP BY bd.id, bd.email, bd.display_name, bp.total_points, 
         bp.available_points, bp.reserved_points, bp.updated_at;
```

### 5. point_transactions_by_auction (オークション別ポイント取引)

```sql
CREATE VIEW point_transactions_by_auction AS
SELECT 
    a.id AS auction_id,
    a.title AS auction_title,
    a.status AS auction_status,
    COUNT(DISTINCT ph.bidder_id) AS unique_bidders,
    COALESCE(SUM(CASE WHEN ph.type = 'reserve' THEN ABS(ph.amount) END), 0) AS total_reserved,
    COALESCE(SUM(CASE WHEN ph.type = 'consume' THEN ABS(ph.amount) END), 0) AS total_consumed,
    COALESCE(SUM(CASE WHEN ph.type = 'release' THEN ph.amount END), 0) AS total_released,
    COALESCE(SUM(CASE WHEN ph.type = 'refund' THEN ph.amount END), 0) AS total_refunded,
    MAX(ph.created_at) AS last_transaction_at
FROM auctions a
LEFT JOIN point_history ph ON a.id = ph.related_auction_id
GROUP BY a.id, a.title, a.status
ORDER BY a.id DESC;
```

## 初期データ

### システム管理者アカウント

```sql
-- パスワード: admin123 (本番環境では必ず変更)
INSERT INTO admins (email, password_hash, role, display_name, status)
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
INSERT INTO admins (email, password_hash, role, display_name, status)
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
INSERT INTO bidders (id, email, password_hash, display_name, status)
VALUES 
    (gen_random_uuid(), 'bidder1@example.com', '$2a$12$9kQwLMRPqZrV6wS9hG5Id.MyL0cF3sX8kU6nD4eQ2gO7iH9mB5rCj', '入札者1', 'active'),
    (gen_random_uuid(), 'bidder2@example.com', '$2a$12$0lRxMNSQraW7xT0iH6Je.NzM1dG4tY9lV7oE5fR3hP8jI0nC6sDk', '入札者2', 'active'),
    (gen_random_uuid(), 'bidder3@example.com', '$2a$12$1mSyNOTRsbX8yU1jI7Kf.O0N2eH5uZ0mW8pF6gS4iQ9kJ1oD7tEl', '入札者3', 'active');

-- テスト用ポイント付与 (各入札者に10,000pt)
UPDATE bidder_points SET total_points = 10000, available_points = 10000 
WHERE bidder_id IN (SELECT id FROM bidders);
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
-- 入札履歴取得クエリ (入札者ごと)
EXPLAIN ANALYZE
SELECT b.id, b.price, b.bid_at, b.is_winning, a.title, i.name
FROM bids b
JOIN auctions a ON b.auction_id = a.id
JOIN items i ON a.id = i.auction_id
WHERE b.bidder_id = $1
ORDER BY b.bid_at DESC
LIMIT 20;

-- オークション詳細取得クエリ（メディア含む）
EXPLAIN ANALYZE
SELECT 
    a.*,
    i.id AS item_id,
    i.name AS item_name,
    i.description AS item_description,
    i.metadata AS item_metadata,
    COUNT(DISTINCT b.id) AS total_bids,
    MAX(b.bid_at) AS last_bid_at,
    json_agg(
        json_build_object(
            'id', im.id,
            'media_type', im.media_type,
            'url', im.url,
            'thumbnail_url', im.thumbnail_url,
            'display_order', im.display_order
        ) ORDER BY im.display_order
    ) FILTER (WHERE im.id IS NOT NULL) AS media
FROM auctions a
LEFT JOIN items i ON a.id = i.auction_id
LEFT JOIN bids b ON a.id = b.auction_id
LEFT JOIN item_media im ON i.id = im.item_id
WHERE a.id = $1
GROUP BY a.id, i.id, i.name, i.description, i.metadata;

-- ポイント履歴取得クエリ (入札者ごと、ページング対応)
EXPLAIN ANALYZE
SELECT 
    ph.id,
    ph.type,
    ph.amount,
    ph.balance_after,
    ph.reserved_after,
    ph.reason,
    a.title AS auction_title,
    ph.created_at
FROM point_history ph
LEFT JOIN auctions a ON ph.related_auction_id = a.id
WHERE ph.bidder_id = $1
ORDER BY ph.created_at DESC
LIMIT 50 OFFSET $2;

-- ポイント残高整合性チェッククエリ
EXPLAIN ANALYZE
SELECT 
    bp.bidder_id,
    bp.total_points,
    bp.available_points,
    bp.reserved_points,
    COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0) AS calculated_total,
    (bp.available_points + bp.reserved_points) AS current_sum
FROM bidder_points bp
LEFT JOIN point_history ph ON bp.bidder_id = ph.bidder_id
GROUP BY bp.bidder_id, bp.total_points, bp.available_points, bp.reserved_points
HAVING bp.total_points != COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0);
```

### 3. パーティショニング (将来的な対応)

大量のデータが蓄積された場合、以下のテーブルをパーティション化：

```sql
-- bids テーブルを月ごとにパーティション化
CREATE TABLE bids_partitioned (
    LIKE bids INCLUDING ALL
) PARTITION BY RANGE (bid_at);

CREATE TABLE bids_2025_01 PARTITION OF bids_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE bids_2025_02 PARTITION OF bids_partitioned
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

-- point_history テーブルを月ごとにパーティション化
CREATE TABLE point_history_partitioned (
    LIKE point_history INCLUDING ALL
) PARTITION BY RANGE (created_at);

CREATE TABLE point_history_2025_01 PARTITION OF point_history_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE point_history_2025_02 PARTITION OF point_history_partitioned
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
```

### ポイント履歴記録の実装例 (Go)

```go
// internal/service/point_service.go
package service

import (
    "context"
    "fmt"
    "gorm.io/gorm"
)

type PointService struct {
    db *gorm.DB
}

// ポイント付与
func (s *PointService) GrantPoints(ctx context.Context, bidderID string, amount int64, reason string, adminID int64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 現在のポイント残高を取得（FOR UPDATE）
        var points BidderPoint
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("bidder_id = ?", bidderID).
            First(&points).Error; err != nil {
            return err
        }

        // 変更前の状態を保存
        balanceBefore := points.AvailablePoints
        reservedBefore := points.ReservedPoints
        totalBefore := points.TotalPoints

        // ポイント付与
        points.TotalPoints += amount
        points.AvailablePoints += amount

        if err := tx.Save(&points).Error; err != nil {
            return err
        }

        // 履歴記録（トリガーで自動作成されるが、手動でも記録可能）
        history := PointHistory{
            BidderID:       bidderID,
            Amount:         amount,
            Type:           "grant",
            Reason:         reason,
            AdminID:        &adminID,
            BalanceBefore:  balanceBefore,
            BalanceAfter:   points.AvailablePoints,
            ReservedBefore: reservedBefore,
            ReservedAfter:  points.ReservedPoints,
            TotalBefore:    totalBefore,
            TotalAfter:     points.TotalPoints,
        }

        return tx.Create(&history).Error
    })
}

// 入札時のポイント予約
func (s *PointService) ReservePoints(ctx context.Context, bidderID string, auctionID int64, bidID int64, amount int64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        var points BidderPoint
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("bidder_id = ?", bidderID).
            First(&points).Error; err != nil {
            return err
        }

        // 残高チェック
        if points.AvailablePoints < amount {
            return fmt.Errorf("insufficient points: available=%d, required=%d", 
                points.AvailablePoints, amount)
        }

        balanceBefore := points.AvailablePoints
        reservedBefore := points.ReservedPoints
        totalBefore := points.TotalPoints

        // ポイント予約（available → reserved）
        points.AvailablePoints -= amount
        points.ReservedPoints += amount

        if err := tx.Save(&points).Error; err != nil {
            return err
        }

        // 履歴記録
        history := PointHistory{
            BidderID:          bidderID,
            Amount:            -amount,
            Type:              "reserve",
            RelatedAuctionID:  &auctionID,
            RelatedBidID:      &bidID,
            BalanceBefore:     balanceBefore,
            BalanceAfter:      points.AvailablePoints,
            ReservedBefore:    reservedBefore,
            ReservedAfter:     points.ReservedPoints,
            TotalBefore:       totalBefore,
            TotalAfter:        points.TotalPoints,
        }

        return tx.Create(&history).Error
    })
}
```

### 本番環境マイグレーション

```bash
# golang-migrate を使用
migrate -path migrations -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" up

# または環境変数を使用
export DB_PASSWORD="your_secure_password"
migrate -path migrations -database "postgresql://auction_user:${DB_PASSWORD}@localhost:5432/auction_db?sslmode=disable" up
```

## セキュリティ考慮事項

1. **パスワードハッシュ化**: bcrypt (cost factor: 12)
2. **SQL インジェクション対策**: GORM のプリペアドステートメント使用
3. **個人情報保護**: email, password_hash は暗号化推奨
4. **監査ログ**: point_history, price_history で全操作を記録
5. **アクセス制限**: PostgreSQL ユーザー権限で Read/Write を分離

## 監視項目

### データベース全般

- **接続数**: `SELECT count(*) FROM pg_stat_activity;`
- **スロークエリ**: `log_min_duration_statement = 1000` (1秒以上)
- **テーブルサイズ**: 
  ```sql
  SELECT 
      schemaname,
      tablename,
      pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
  FROM pg_tables 
  WHERE schemaname = 'public'
  ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
  ```
- **インデックス使用率**: `SELECT * FROM pg_stat_user_indexes;`

### ポイント履歴関連

- **日次ポイント変動量**:
  ```sql
  SELECT 
      DATE(created_at) AS date,
      type,
      COUNT(*) AS transaction_count,
      SUM(ABS(amount)) AS total_amount
  FROM point_history
  WHERE created_at >= NOW() - INTERVAL '7 days'
  GROUP BY DATE(created_at), type
  ORDER BY date DESC, type;
  ```

- **ポイント残高整合性チェック**:
  ```sql
  -- 履歴の合計と現在の残高が一致するか確認
  SELECT 
      bp.bidder_id,
      bd.email,
      bp.total_points AS current_total,
      COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0) AS granted_total,
      bp.available_points + bp.reserved_points AS sum_of_parts
  FROM bidder_points bp
  JOIN bidders bd ON bp.bidder_id = bd.id
  LEFT JOIN point_history ph ON bp.bidder_id = ph.bidder_id
  GROUP BY bp.bidder_id, bd.email, bp.total_points, bp.available_points, bp.reserved_points
  HAVING bp.total_points != COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0)
      OR bp.total_points < (bp.available_points + bp.reserved_points);
  ```

- **異常なポイント取引の検出**:
  ```sql
  -- 短時間に大量のポイント変動があった入札者を検出
  SELECT 
      bidder_id,
      COUNT(*) AS transaction_count,
      SUM(ABS(amount)) AS total_amount,
      MIN(created_at) AS first_transaction,
      MAX(created_at) AS last_transaction
  FROM point_history
  WHERE created_at >= NOW() - INTERVAL '1 hour'
  GROUP BY bidder_id
  HAVING COUNT(*) > 50 OR SUM(ABS(amount)) > 100000
  ORDER BY transaction_count DESC;
  ```

## 商品メディア管理の設計方針

### 複数メディア対応の理由

1. **柔軟な商品展示**
   - 1つの商品に対して複数の角度からの画像を登録可能
   - 動画による詳細な商品紹介が可能
   - 表示順序を制御可能

2. **パフォーマンスの最適化**
   - メディアファイルは別テーブルで管理し、必要に応じて取得
   - サムネイルURLを別途保存することで一覧表示を高速化
   - メディアの追加・削除が商品テーブルに影響しない

3. **スケーラビリティ**
   - 大量の画像・動画に対応可能
   - CDN（CloudFront等）との連携が容易
   - メディアのみのキャッシュ戦略が可能

4. **汎用性**
   - 競走馬以外の商品（美術品、骨董品、車両等）にも対応
   - 商品カテゴリごとに異なるメディア構成が可能

### metadata フィールドの活用

`items.metadata` はJSONB型で、商品カテゴリに応じて自由にフィールドを定義可能：

**競走馬の場合:**
```json
{
  "category": "horse",
  "birth_date": "2022-03-15",
  "gender": "male",
  "pedigree": {...}
}
```

**美術品の場合:**
```json
{
  "category": "art",
  "artist": "山田太郎",
  "year": "2023",
  "technique": "油彩"
}
```

**車両の場合:**
```json
{
  "category": "vehicle",
  "make": "Toyota",
  "model": "Land Cruiser",
  "year": 2023,
  "mileage": 5000
}
```

### メディア取得の最適化

```sql
-- 商品一覧表示用（メインメディアのみ）
SELECT 
    i.*,
    (SELECT url FROM item_media 
     WHERE item_id = i.id 
     ORDER BY display_order LIMIT 1) AS main_media_url
FROM items i;

-- 商品詳細表示用（全メディア取得）
SELECT 
    i.*,
    json_agg(
        json_build_object(
            'type', im.media_type,
            'url', im.url,
            'thumbnail', im.thumbnail_url
        ) ORDER BY im.display_order
    ) AS media
FROM items i
LEFT JOIN item_media im ON i.id = im.item_id
WHERE i.id = $1
GROUP BY i.id;
```

## テーブル分離の設計方針

### 入札者と管理者を分離する理由

1. **セキュリティの向上**
   - 管理者テーブルと入札者テーブルを物理的に分離することで、権限昇格攻撃のリスクを軽減
   - 入札者情報の漏洩が管理者情報に影響しない

2. **データ構造の最適化**
   - 入札者には`bidder_points`が必須だが、管理者には不要
   - 管理者には`role`(system_admin/auctioneer)が必要だが、入札者には不要
   - 各テーブルに必要なカラムのみを持つことで正規化を実現

3. **パフォーマンスの向上**
   - 入札者の検索クエリが管理者データを含まないため高速化
   - インデックスサイズの削減
   - 入札者テーブルのみをパーティション化可能

4. **スケーラビリティ**
   - 入札者テーブルは将来的に数万〜数十万レコードに成長する可能性
   - 管理者テーブルは数十〜数百レコード程度
   - 異なる成長率に応じた最適化戦略が可能

### UUIDを採用する理由

**入札者IDにUUIDを採用した理由:**

1. **プライバシー保護**
   - 連番IDでは入札者数や登録順序が推測可能
   - UUIDでは入札者の個人情報推測が困難

2. **分散システム対応**
   - 複数のアプリケーションサーバーで独立してIDを生成可能
   - データベースへの問い合わせなしでID生成が可能

3. **URLの推測防止**
   - `/bidders/1`, `/bidders/2` のような連番URLは容易に推測される
   - UUIDを使用することで不正アクセスのリスクを軽減

4. **マージ・統合の容易さ**
   - 異なる環境（開発・ステージング・本番）でIDの衝突が発生しない
   - データ移行時の主キー競合を回避

**管理者IDをBIGSERIALのままにした理由:**

1. **管理者数は限定的**
   - 数十〜数百程度の規模で、UUIDのメリットが小さい

2. **内部管理の利便性**
   - 管理者は内部スタッフのみなので連番でも問題ない
   - ログやデバッグ時にIDが短く読みやすい

3. **パフォーマンス**
   - 整数型の方がUUIDより検索・JOIN性能が高い
   - インデックスサイズが小さい

## まとめ

このデータベース設計は以下の要件を満たします：

✅ **ポイント制**: 仮想ポイントによる安全な入札管理  
✅ **主催者主導**: 価格開示履歴の完全な記録  
✅ **整合性**: トランザクションとトリガーによるデータ整合性保証  
✅ **パフォーマンス**: 適切なインデックスとビューによる高速クエリ  
✅ **監査**: 全ポイント操作・価格開示の履歴記録  
✅ **スケーラビリティ**: パーティショニング対応可能な設計  
✅ **セキュリティ**: ロールベースのアクセス制御とデータ保護  
✅ **テーブル分離**: 入札者と管理者の明確な分離による最適化  
✅ **プライバシー**: UUID採用による入札者情報の保護
