-- Migration: 011_add_lot_number_to_items (rollback)
-- Description: itemsテーブルからlot_numberカラムを削除
-- Date: 2025-11-23

BEGIN;

-- Step 1: チェック制約を削除
ALTER TABLE items DROP CONSTRAINT IF EXISTS chk_items_lot_number_positive;

-- Step 2: ユニーク制約を削除
ALTER TABLE items DROP CONSTRAINT IF EXISTS uk_items_auction_lot;

-- Step 3: インデックスを削除
DROP INDEX IF EXISTS idx_items_auction_lot;

-- Step 4: lot_numberカラムを削除
ALTER TABLE items DROP COLUMN IF EXISTS lot_number;

COMMIT;
