-- Migration: 013_restructure_price_history_for_items
-- Description: price_historyテーブルをauction_idからitem_idに変更し、disclosed_byフィールドを追加
-- Date: 2025-11-26

BEGIN;

-- Step 1: price_historyテーブルに新しいitem_idカラムを追加
ALTER TABLE price_history ADD COLUMN item_id UUID;

-- Step 2: 既存のprice_historyデータを移行（auction_idから対応するitem_idを設定）
-- 注意: 1つのオークションに複数の商品がある場合、最初の商品に紐付ける
UPDATE price_history ph
SET item_id = (
    SELECT i.id
    FROM items i
    WHERE i.auction_id = ph.auction_id
    ORDER BY i.lot_number ASC
    LIMIT 1
)
WHERE ph.auction_id IS NOT NULL;

-- Step 3: item_idをNOT NULLに変更
ALTER TABLE price_history ALTER COLUMN item_id SET NOT NULL;

-- Step 4: 古いインデックスを削除
DROP INDEX IF EXISTS idx_price_history_auction;
DROP INDEX IF EXISTS idx_price_history_auction_id;

-- Step 5: 古い外部キー制約を削除
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS price_history_auction_id_fkey;
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_auction;

-- Step 6: 新しいインデックスを作成
CREATE INDEX idx_price_history_item ON price_history(item_id, opened_at DESC);

-- Step 7: 新しい外部キー制約を追加
ALTER TABLE price_history
    ADD CONSTRAINT fk_price_history_item
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE;

-- Step 8: auction_idカラムを削除
ALTER TABLE price_history DROP COLUMN auction_id;

-- Step 9: opened_byカラムの名前をdisclosed_byに変更（意味をより明確に）
ALTER TABLE price_history RENAME COLUMN opened_by TO disclosed_by;
ALTER TABLE price_history RENAME COLUMN opened_at TO disclosed_at;

-- Step 10: 外部キー制約の名前を更新
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_opener;
ALTER TABLE price_history
    ADD CONSTRAINT fk_price_history_disclosed_by
    FOREIGN KEY (disclosed_by) REFERENCES admins(id) ON DELETE CASCADE;

-- Step 11: インデックスの名前を更新
DROP INDEX IF EXISTS idx_price_history_opened_by;
CREATE INDEX idx_price_history_disclosed_by ON price_history(disclosed_by);

COMMIT;
