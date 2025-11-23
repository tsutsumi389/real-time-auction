-- Migration: 010_restructure_auction_and_bids
-- Description: オークション構造の変更
--   1. auctionsテーブルから商品固有の項目をitemsテーブルに移動
--   2. bidsテーブルのauction_idをitem_idに変更
-- Date: 2025-11-23

BEGIN;

-- Step 1: 既存のビューを削除（auction_idに依存しているため）
DROP VIEW IF EXISTS active_auctions_view;
DROP VIEW IF EXISTS bidder_auction_summary;

-- Step 2: itemsテーブルに新しいカラムを追加
ALTER TABLE items
    ADD COLUMN starting_price BIGINT,
    ADD COLUMN current_price BIGINT,
    ADD COLUMN winner_id UUID,
    ADD COLUMN started_at TIMESTAMPTZ,
    ADD COLUMN ended_at TIMESTAMPTZ;

-- Step 3: itemsテーブルにインデックスを追加
CREATE INDEX idx_items_started_at ON items(started_at);
CREATE INDEX idx_items_ended_at ON items(ended_at);
CREATE INDEX idx_items_winner ON items(winner_id);

-- Step 4: itemsテーブルに制約を追加
ALTER TABLE items 
    ADD CONSTRAINT fk_items_winner 
    FOREIGN KEY (winner_id) REFERENCES bidders(id) ON DELETE SET NULL;

ALTER TABLE items 
    ADD CONSTRAINT chk_items_price_positive 
    CHECK (starting_price IS NULL OR starting_price > 0);

ALTER TABLE items 
    ADD CONSTRAINT chk_items_dates 
    CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at);

-- Step 5: 既存データをauctionsからitemsに移行（もしデータが存在する場合）
UPDATE items i
SET 
    starting_price = a.starting_price,
    current_price = a.current_price,
    winner_id = a.winner_id,
    started_at = a.started_at,
    ended_at = a.ended_at
FROM auctions a
WHERE i.auction_id = a.id
    AND (a.starting_price IS NOT NULL 
         OR a.current_price IS NOT NULL 
         OR a.winner_id IS NOT NULL 
         OR a.started_at IS NOT NULL 
         OR a.ended_at IS NOT NULL);

-- Step 6: bidsテーブルに新しいitem_idカラムを追加
ALTER TABLE bids ADD COLUMN item_id UUID;

-- Step 7: 既存のbidsデータを移行（auction_idから対応するitem_idを設定）
-- 注意: 1つのオークションに複数の商品がある場合、最初の商品に紐付ける
UPDATE bids b
SET item_id = (
    SELECT i.id
    FROM items i
    WHERE i.auction_id = b.auction_id
    ORDER BY i.created_at
    LIMIT 1
)
WHERE b.auction_id IS NOT NULL;

-- Step 8: item_idをNOT NULLに変更
ALTER TABLE bids ALTER COLUMN item_id SET NOT NULL;

-- Step 9: bidsテーブルの古いインデックスを削除
DROP INDEX IF EXISTS idx_bids_auction;
DROP INDEX IF EXISTS idx_bids_winning;
DROP INDEX IF EXISTS idx_bids_auction_bidder;

-- Step 10: bidsテーブルの古い制約を削除
ALTER TABLE bids DROP CONSTRAINT IF EXISTS fk_bids_auction;

-- Step 11: bidsテーブルに新しいインデックスを作成
CREATE INDEX idx_bids_item ON bids(item_id, bid_at DESC);
CREATE INDEX idx_bids_winning ON bids(item_id, is_winning) WHERE is_winning = TRUE;
CREATE INDEX idx_bids_item_bidder ON bids(item_id, bidder_id);

-- Step 12: bidsテーブルに新しい制約を追加
ALTER TABLE bids 
    ADD CONSTRAINT fk_bids_item 
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE;

-- Step 13: bidsテーブルからauction_idカラムを削除
ALTER TABLE bids DROP COLUMN auction_id;

-- Step 14: auctionsテーブルのインデックスを削除
DROP INDEX IF EXISTS idx_auctions_started_at;
DROP INDEX IF EXISTS idx_auctions_ended_at;
DROP INDEX IF EXISTS idx_auctions_winner;

-- Step 15: auctionsテーブルの制約を削除
ALTER TABLE auctions DROP CONSTRAINT IF EXISTS fk_auctions_winner;
ALTER TABLE auctions DROP CONSTRAINT IF EXISTS chk_auctions_price_positive;
ALTER TABLE auctions DROP CONSTRAINT IF EXISTS chk_auctions_dates;

-- Step 16: auctionsテーブルからカラムを削除
ALTER TABLE auctions 
    DROP COLUMN IF EXISTS starting_price,
    DROP COLUMN IF EXISTS current_price,
    DROP COLUMN IF EXISTS winner_id,
    DROP COLUMN IF EXISTS started_at,
    DROP COLUMN IF EXISTS ended_at;

-- Step 17: ビューを再作成
DROP VIEW IF EXISTS active_auctions_view;
CREATE VIEW active_auctions_view AS
SELECT 
    a.id,
    a.title,
    a.description,
    a.status,
    i.id AS item_id,
    i.name AS item_name,
    i.current_price,
    i.started_at,
    i.ended_at,
    (SELECT url FROM item_media WHERE item_id = i.id ORDER BY display_order LIMIT 1) AS main_image_url,
    COUNT(DISTINCT b.bidder_id) AS bidder_count,
    MAX(b.bid_at) AS last_bid_at,
    bd.display_name AS winner_name
FROM auctions a
LEFT JOIN items i ON a.id = i.auction_id
LEFT JOIN bids b ON i.id = b.item_id
LEFT JOIN bidders bd ON i.winner_id = bd.id
WHERE a.status IN ('active', 'ended')
GROUP BY a.id, i.id, i.name, i.current_price, i.started_at, i.ended_at, bd.display_name
ORDER BY i.started_at DESC;

DROP VIEW IF EXISTS bidder_auction_summary;
CREATE VIEW bidder_auction_summary AS
SELECT 
    bd.id AS bidder_id,
    bd.email,
    bd.display_name,
    bp.available_points,
    bp.reserved_points,
    COUNT(DISTINCT i.auction_id) AS participated_auctions,
    COUNT(DISTINCT CASE WHEN i.winner_id = bd.id THEN i.id END) AS won_items,
    COALESCE(SUM(CASE WHEN i.winner_id = bd.id THEN i.current_price END), 0) AS total_won_amount
FROM bidders bd
LEFT JOIN bidder_points bp ON bd.id = bp.bidder_id
LEFT JOIN bids b ON bd.id = b.bidder_id
LEFT JOIN items i ON b.item_id = i.id AND i.ended_at IS NOT NULL
GROUP BY bd.id, bd.email, bd.display_name, bp.available_points, bp.reserved_points;

COMMIT;
