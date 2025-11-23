-- Migration: 011_add_lot_number_to_items
-- Description: itemsテーブルにlot_numberカラムを追加
-- Date: 2025-11-23

BEGIN;

-- Step 1: itemsテーブルにlot_numberカラムを追加
ALTER TABLE items ADD COLUMN lot_number INT;

-- Step 2: 既存のデータにlot_numberを設定（auction_id別に連番を付与）
WITH numbered_items AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (PARTITION BY auction_id ORDER BY created_at, id) AS lot_num
    FROM items
)
UPDATE items
SET lot_number = numbered_items.lot_num
FROM numbered_items
WHERE items.id = numbered_items.id;

-- Step 3: lot_numberをNOT NULLに変更
ALTER TABLE items ALTER COLUMN lot_number SET NOT NULL;

-- Step 4: インデックスを追加（auction_id + lot_numberで検索する場合に有効）
CREATE INDEX idx_items_auction_lot ON items(auction_id, lot_number);

-- Step 5: ユニーク制約を追加（同じオークション内でlot_numberが重複しないように）
ALTER TABLE items ADD CONSTRAINT uk_items_auction_lot UNIQUE (auction_id, lot_number);

-- Step 6: チェック制約を追加（lot_numberは1以上）
ALTER TABLE items ADD CONSTRAINT chk_items_lot_number_positive CHECK (lot_number >= 1);

COMMIT;
