-- Migration: Convert auctions and items table IDs from BIGSERIAL to UUID
-- This migration changes primary keys and all related foreign keys

-- Step 1: Enable uuid-ossp extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Step 2: Drop dependent views first
DROP VIEW IF EXISTS active_auctions_view CASCADE;
DROP VIEW IF EXISTS bidder_auction_summary CASCADE;
DROP VIEW IF EXISTS point_history_detailed CASCADE;
DROP VIEW IF EXISTS bidder_point_balance_view CASCADE;
DROP VIEW IF EXISTS point_transactions_by_auction CASCADE;

-- Step 3: Add temporary UUID columns
ALTER TABLE auctions ADD COLUMN id_uuid UUID DEFAULT gen_random_uuid();
ALTER TABLE items ADD COLUMN id_uuid UUID DEFAULT gen_random_uuid();
ALTER TABLE items ADD COLUMN auction_id_uuid UUID;

-- Step 4: Populate UUID columns in dependent tables
UPDATE items i SET auction_id_uuid = a.id_uuid 
FROM auctions a WHERE i.auction_id = a.id;

ALTER TABLE item_media ADD COLUMN item_id_uuid UUID;
UPDATE item_media im SET item_id_uuid = i.id_uuid 
FROM items i WHERE im.item_id = i.id;

ALTER TABLE bids ADD COLUMN auction_id_uuid UUID;
UPDATE bids b SET auction_id_uuid = a.id_uuid 
FROM auctions a WHERE b.auction_id = a.id;

ALTER TABLE price_history ADD COLUMN auction_id_uuid UUID;
UPDATE price_history ph SET auction_id_uuid = a.id_uuid 
FROM auctions a WHERE ph.auction_id = a.id;

ALTER TABLE point_history ADD COLUMN related_auction_id_uuid UUID;
UPDATE point_history ph SET related_auction_id_uuid = a.id_uuid 
FROM auctions a WHERE ph.related_auction_id = a.id;

-- Step 5: Drop old foreign key constraints
ALTER TABLE items DROP CONSTRAINT IF EXISTS items_auction_id_fkey;
ALTER TABLE items DROP CONSTRAINT IF EXISTS fk_items_auction;
ALTER TABLE item_media DROP CONSTRAINT IF EXISTS item_media_item_id_fkey;
ALTER TABLE item_media DROP CONSTRAINT IF EXISTS fk_item_media_item;
ALTER TABLE bids DROP CONSTRAINT IF EXISTS bids_auction_id_fkey;
ALTER TABLE bids DROP CONSTRAINT IF EXISTS fk_bids_auction;
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS price_history_auction_id_fkey;
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_auction;
ALTER TABLE point_history DROP CONSTRAINT IF EXISTS point_history_related_auction_id_fkey;
ALTER TABLE point_history DROP CONSTRAINT IF EXISTS fk_point_history_auction;

-- Step 5: Drop old foreign key constraints
ALTER TABLE items DROP CONSTRAINT IF EXISTS items_auction_id_fkey;
ALTER TABLE items DROP CONSTRAINT IF EXISTS fk_items_auction;
ALTER TABLE item_media DROP CONSTRAINT IF EXISTS item_media_item_id_fkey;
ALTER TABLE item_media DROP CONSTRAINT IF EXISTS fk_item_media_item;
ALTER TABLE bids DROP CONSTRAINT IF EXISTS bids_auction_id_fkey;
ALTER TABLE bids DROP CONSTRAINT IF EXISTS fk_bids_auction;
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS price_history_auction_id_fkey;
ALTER TABLE price_history DROP CONSTRAINT IF EXISTS fk_price_history_auction;
ALTER TABLE point_history DROP CONSTRAINT IF EXISTS point_history_related_auction_id_fkey;
ALTER TABLE point_history DROP CONSTRAINT IF EXISTS fk_point_history_auction;

-- Step 6: Drop old primary key constraints and columns
ALTER TABLE auctions DROP CONSTRAINT IF EXISTS auctions_pkey;
ALTER TABLE auctions DROP COLUMN id;

ALTER TABLE items DROP CONSTRAINT IF EXISTS items_pkey;
ALTER TABLE items DROP COLUMN id;
ALTER TABLE items DROP COLUMN auction_id;

ALTER TABLE item_media DROP COLUMN item_id;
ALTER TABLE bids DROP COLUMN auction_id;
ALTER TABLE price_history DROP COLUMN auction_id;
ALTER TABLE point_history DROP COLUMN related_auction_id;

-- Step 7: Rename UUID columns to id
ALTER TABLE auctions RENAME COLUMN id_uuid TO id;
ALTER TABLE items RENAME COLUMN id_uuid TO id;
ALTER TABLE items RENAME COLUMN auction_id_uuid TO auction_id;
ALTER TABLE item_media RENAME COLUMN item_id_uuid TO item_id;
ALTER TABLE bids RENAME COLUMN auction_id_uuid TO auction_id;
ALTER TABLE price_history RENAME COLUMN auction_id_uuid TO auction_id;
ALTER TABLE point_history RENAME COLUMN related_auction_id_uuid TO related_auction_id;

-- Step 8: Add new primary key constraints
ALTER TABLE auctions ADD PRIMARY KEY (id);
ALTER TABLE items ADD PRIMARY KEY (id);

-- Step 9: Add new foreign key constraints
ALTER TABLE items 
  ADD CONSTRAINT items_auction_id_fkey 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;

ALTER TABLE item_media 
  ADD CONSTRAINT item_media_item_id_fkey 
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE;

ALTER TABLE bids 
  ADD CONSTRAINT bids_auction_id_fkey 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;

ALTER TABLE price_history 
  ADD CONSTRAINT price_history_auction_id_fkey 
  FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE;

ALTER TABLE point_history 
  ADD CONSTRAINT point_history_related_auction_id_fkey 
  FOREIGN KEY (related_auction_id) REFERENCES auctions(id) ON DELETE SET NULL;

-- Step 10: Update default values for new records
ALTER TABLE auctions ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE items ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Step 11: Add NOT NULL constraints where appropriate
ALTER TABLE items ALTER COLUMN auction_id SET NOT NULL;
ALTER TABLE item_media ALTER COLUMN item_id SET NOT NULL;
ALTER TABLE bids ALTER COLUMN auction_id SET NOT NULL;
ALTER TABLE price_history ALTER COLUMN auction_id SET NOT NULL;
-- Note: point_history.related_auction_id is nullable

-- Step 12: Recreate indexes if they existed
CREATE INDEX IF NOT EXISTS idx_items_auction_id ON items(auction_id);
CREATE INDEX IF NOT EXISTS idx_bids_auction_id ON bids(auction_id);
CREATE INDEX IF NOT EXISTS idx_price_history_auction_id ON price_history(auction_id);
CREATE INDEX IF NOT EXISTS idx_point_history_related_auction_id ON point_history(related_auction_id);

-- Step 13: Recreate views
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

CREATE VIEW point_history_detailed AS
SELECT 
    ph.id,
    ph.bidder_id,
    bd.email AS bidder_email,
    bd.display_name AS bidder_name,
    ph.amount,
    ph.type,
    ph.reason,
    ph.balance_before,
    ph.balance_after,
    ph.reserved_before,
    ph.reserved_after,
    ph.total_before,
    ph.total_after,
    a.title AS auction_title,
    b.price AS bid_price,
    ad.display_name AS admin_name,
    ph.created_at
FROM point_history ph
JOIN bidders bd ON ph.bidder_id = bd.id
LEFT JOIN auctions a ON ph.related_auction_id = a.id
LEFT JOIN bids b ON ph.related_bid_id = b.id
LEFT JOIN admins ad ON ph.admin_id = ad.id
ORDER BY ph.created_at DESC;

CREATE VIEW bidder_point_balance_view AS
SELECT 
    bd.id AS bidder_id,
    bd.email,
    bd.display_name,
    bp.total_points,
    bp.available_points,
    bp.reserved_points,
    (bp.total_points - bp.available_points - bp.reserved_points) AS consumed_points,
    COALESCE(SUM(CASE WHEN ph.type = 'grant' THEN ph.amount END), 0) AS total_granted,
    COALESCE(SUM(CASE WHEN ph.type = 'consume' THEN ABS(ph.amount) END), 0) AS total_consumed,
    COALESCE(COUNT(CASE WHEN ph.type = 'reserve' THEN 1 END), 0) AS total_bids,
    bp.updated_at AS last_updated
FROM bidders bd
LEFT JOIN bidder_points bp ON bd.id = bp.bidder_id
LEFT JOIN point_history ph ON bd.id = ph.bidder_id
GROUP BY bd.id, bd.email, bd.display_name, bp.total_points, 
         bp.available_points, bp.reserved_points, bp.updated_at;

CREATE VIEW point_transactions_by_auction AS
SELECT 
    a.id AS auction_id,
    a.title AS auction_title,
    a.status AS auction_status,
    COUNT(DISTINCT ph.bidder_id) AS unique_bidders,
    COALESCE(SUM(CASE WHEN ph.type = 'reserve' THEN ABS(ph.amount) END), 0) AS total_reserved,
    COALESCE(SUM(CASE WHEN ph.type = 'consume' THEN ABS(ph.amount) END), 0) AS total_consumed,
    COALESCE(SUM(CASE WHEN ph.type = 'release' THEN ph.amount END), 0) AS total_released,
    COALESCE(SUM(CASE WHEN ph.type = 'refund' THEN ph.amount END), 0) AS total_refunded,
    MAX(ph.created_at) AS last_transaction_at
FROM auctions a
LEFT JOIN point_history ph ON a.id = ph.related_auction_id
GROUP BY a.id, a.title, a.status
ORDER BY a.id DESC;
