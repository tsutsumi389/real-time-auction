-- Migration: 014_make_item_auction_nullable (rollback)
-- Description: itemsテーブルのauction_idをNOT NULLに戻す
-- Date: 2025-11-30

BEGIN;

-- Step 1: 部分インデックスを削除
DROP INDEX IF EXISTS uk_items_auction_lot_when_assigned;

-- Step 2: 新しいチェック制約を削除
ALTER TABLE items DROP CONSTRAINT IF EXISTS chk_items_lot_number_non_negative;

-- Step 3: lot_numberのデフォルト値を削除
ALTER TABLE items ALTER COLUMN lot_number DROP DEFAULT;

-- Step 4: auction_idがNULLのレコードを削除（ロールバック時にデータ損失の可能性あり）
DELETE FROM items WHERE auction_id IS NULL;

-- Step 5: auction_idをNOT NULLに戻す
ALTER TABLE items ALTER COLUMN auction_id SET NOT NULL;

-- Step 6: 元のチェック制約を追加（lot_number >= 1）
ALTER TABLE items ADD CONSTRAINT chk_items_lot_number_positive CHECK (lot_number >= 1);

-- Step 7: 元のユニーク制約を追加
ALTER TABLE items ADD CONSTRAINT uk_items_auction_lot UNIQUE (auction_id, lot_number);

-- Step 8: ビューを再作成
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
