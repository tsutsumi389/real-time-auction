-- Migration: 012_add_started_at_to_auctions
-- Description: auctionsテーブルにstarted_at（オークション開始日時）カラムを追加
--   オークション全体の開始日時を管理するために必要
--   入札者用一覧画面で表示・ソートに使用
-- Date: 2025-11-24

BEGIN;

-- Step 1: auctionsテーブルにstarted_atカラムを追加
ALTER TABLE auctions ADD COLUMN started_at TIMESTAMPTZ;

-- Step 2: インデックスを作成（ソート性能向上のため）
CREATE INDEX idx_auctions_started_at ON auctions(started_at);

-- Step 3: コメントを追加
COMMENT ON COLUMN auctions.started_at IS 'オークション全体の開始日時（入札者用一覧画面で表示・ソートに使用）';

COMMIT;
