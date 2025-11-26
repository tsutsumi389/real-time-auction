-- Migration rollback: 013_restructure_price_history_for_items
-- Description: Revert price_history from item_id back to auction_id

BEGIN;

-- Step 1: price_historyテーブルにauction_idカラムを追加
ALTER TABLE price_history ADD COLUMN auction_id UUID;

-- Step 2: item_idからauction_idにデータを移行
UPDATE price_history ph
SET auction_id = (
    SELECT i.auction_id
    FROM items i
    WHERE i.id = ph.item_id
)
WHERE ph.item_id IS NOT NULL;

-- Step 3: auction_idをNOT NULLに変更
ALTER TABLE price_history ALTER COLUMN auction_id SET NOT NULL;

-- Step 4: インデックスの名前を戻す
DROP INDEX IF EXISTS idx_price_history_disclosed_by;
CREATE INDEX idx_price_history_opened_by ON price_history(disclosed_by);

-- Step 5: 外部キー制約の名前を戻す
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_disclosed_by;
ALTER TABLE price_history
    ADD CONSTRAINT fk_price_history_opener
    FOREIGN KEY (disclosed_by) REFERENCES admins(id) ON DELETE CASCADE;

-- Step 6: カラム名を戻す
ALTER TABLE price_history RENAME COLUMN disclosed_by TO opened_by;
ALTER TABLE price_history RENAME COLUMN disclosed_at TO opened_at;

-- Step 7: 新しいインデックスを削除
DROP INDEX IF EXISTS idx_price_history_item;

-- Step 8: 新しい外部キー制約を削除
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_item;

-- Step 9: 古いインデックスを復元
CREATE INDEX idx_price_history_auction ON price_history(auction_id, opened_at DESC);

-- Step 10: 古い外部キー制約を復元
ALTER TABLE price_history
    ADD CONSTRAINT fk_price_history_auction
    FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;

-- Step 11: item_idカラムを削除
ALTER TABLE price_history DROP COLUMN item_id;

COMMIT;
