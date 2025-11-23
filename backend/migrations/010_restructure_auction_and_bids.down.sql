-- Migration: 010_restructure_auction_and_bids (Rollback)
-- Description: オークション構造の変更をロールバック
-- Date: 2025-11-23

BEGIN;

-- Step 1: ビューを削除
DROP VIEW IF EXISTS bidder_auction_summary;
DROP VIEW IF EXISTS active_auctions_view;

-- Step 2: auctionsテーブルにカラムを復元
ALTER TABLE auctions
    ADD COLUMN starting_price BIGINT,
    ADD COLUMN current_price BIGINT,
    ADD COLUMN winner_id UUID,
    ADD COLUMN started_at TIMESTAMPTZ,
    ADD COLUMN ended_at TIMESTAMPTZ;

-- Step 3: auctionsテーブルに制約を復元
ALTER TABLE auctions 
    ADD CONSTRAINT fk_auctions_winner 
    FOREIGN KEY (winner_id) REFERENCES bidders(id) ON DELETE SET NULL;

ALTER TABLE auctions 
    ADD CONSTRAINT chk_auctions_price_positive 
    CHECK (starting_price IS NULL OR starting_price > 0);

ALTER TABLE auctions 
    ADD CONSTRAINT chk_auctions_dates 
    CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at);

-- Step 4: auctionsテーブルにインデックスを復元
CREATE INDEX idx_auctions_started_at ON auctions(started_at);
CREATE INDEX idx_auctions_ended_at ON auctions(ended_at);
CREATE INDEX idx_auctions_winner ON auctions(winner_id);

-- Step 5: データをitemsからauctionsに戻す
UPDATE auctions a
SET 
    starting_price = i.starting_price,
    current_price = i.current_price,
    winner_id = i.winner_id,
    started_at = i.started_at,
    ended_at = i.ended_at
FROM (
    SELECT DISTINCT ON (auction_id)
        auction_id,
        starting_price,
        current_price,
        winner_id,
        started_at,
        ended_at
    FROM items
    WHERE starting_price IS NOT NULL 
       OR current_price IS NOT NULL 
       OR winner_id IS NOT NULL 
       OR started_at IS NOT NULL 
       OR ended_at IS NOT NULL
    ORDER BY auction_id, created_at
) i
WHERE a.id = i.auction_id;

-- Step 6: bidsテーブルにauction_idカラムを復元
ALTER TABLE bids ADD COLUMN auction_id UUID;

-- Step 7: bidsデータを復元（item_idからauction_idを設定）
UPDATE bids b
SET auction_id = i.auction_id
FROM items i
WHERE b.item_id = i.id;

-- Step 8: auction_idをNOT NULLに変更
ALTER TABLE bids ALTER COLUMN auction_id SET NOT NULL;

-- Step 9: bidsテーブルの新しいインデックスを削除
DROP INDEX IF EXISTS idx_bids_item;
DROP INDEX IF EXISTS idx_bids_winning;
DROP INDEX IF EXISTS idx_bids_item_bidder;

-- Step 10: bidsテーブルの新しい制約を削除
ALTER TABLE bids DROP CONSTRAINT IF EXISTS fk_bids_item;

-- Step 11: bidsテーブルの古いインデックスを復元
CREATE INDEX idx_bids_auction ON bids(auction_id, bid_at DESC);
CREATE INDEX idx_bids_winning ON bids(auction_id, is_winning) WHERE is_winning = TRUE;
CREATE INDEX idx_bids_auction_bidder ON bids(auction_id, bidder_id);

-- Step 12: bidsテーブルの古い制約を復元
ALTER TABLE bids 
    ADD CONSTRAINT fk_bids_auction 
    FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;

-- Step 13: bidsテーブルからitem_idカラムを削除
ALTER TABLE bids DROP COLUMN item_id;

-- Step 14: itemsテーブルの制約を削除
ALTER TABLE items DROP CONSTRAINT IF EXISTS fk_items_winner;
ALTER TABLE items DROP CONSTRAINT IF EXISTS chk_items_price_positive;
ALTER TABLE items DROP CONSTRAINT IF EXISTS chk_items_dates;

-- Step 15: itemsテーブルのインデックスを削除
DROP INDEX IF EXISTS idx_items_started_at;
DROP INDEX IF EXISTS idx_items_ended_at;
DROP INDEX IF EXISTS idx_items_winner;

-- Step 16: itemsテーブルからカラムを削除
ALTER TABLE items 
    DROP COLUMN IF EXISTS starting_price,
    DROP COLUMN IF EXISTS current_price,
    DROP COLUMN IF EXISTS winner_id,
    DROP COLUMN IF EXISTS started_at,
    DROP COLUMN IF EXISTS ended_at;

-- Step 17: ビューを復元
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

COMMIT;
