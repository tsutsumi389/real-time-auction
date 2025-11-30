-- Migration: 014_make_item_auction_nullable
-- Description: itemsテーブルのauction_idをNULL許可に変更し、商品をオークション未割当でも登録可能にする
-- Date: 2025-11-30

BEGIN;

-- Step 1: 既存のユニーク制約を削除
ALTER TABLE items DROP CONSTRAINT IF EXISTS uk_items_auction_lot;

-- Step 2: 既存のチェック制約を削除（lot_number >= 1）
ALTER TABLE items DROP CONSTRAINT IF EXISTS chk_items_lot_number_positive;

-- Step 3: auction_idをNULL許可に変更
ALTER TABLE items ALTER COLUMN auction_id DROP NOT NULL;

-- Step 4: lot_numberのデフォルト値を設定（未割当の場合は0）
ALTER TABLE items ALTER COLUMN lot_number SET DEFAULT 0;

-- Step 5: 新しいチェック制約を追加（lot_numberは0以上）
-- 0: 未割当、1以上: オークションに割当済み
ALTER TABLE items ADD CONSTRAINT chk_items_lot_number_non_negative CHECK (lot_number >= 0);

-- Step 6: 部分インデックスを作成（auction_idがNULLでない場合のみユニーク）
CREATE UNIQUE INDEX uk_items_auction_lot_when_assigned 
    ON items(auction_id, lot_number) 
    WHERE auction_id IS NOT NULL;

-- Step 7: ビューを再作成（NULLのauction_idを考慮）
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

COMMIT;
