-- Migration: 012_add_started_at_to_auctions (Rollback)
-- Description: auctionsテーブルからstarted_atカラムを削除
-- Date: 2025-11-24

BEGIN;

-- Step 1: インデックスを削除
DROP INDEX IF EXISTS idx_auctions_started_at;

-- Step 2: カラムを削除
ALTER TABLE auctions DROP COLUMN IF EXISTS started_at;

COMMIT;
