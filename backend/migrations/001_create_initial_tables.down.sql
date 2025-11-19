-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS point_history;
DROP TABLE IF EXISTS price_history;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS item_media;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS auctions;
DROP TABLE IF EXISTS bidder_points;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS bidders;

-- Drop extension
DROP EXTENSION IF EXISTS "pgcrypto";
