-- Drop triggers
DROP TRIGGER IF EXISTS trigger_record_point_history ON bidder_points;
DROP TRIGGER IF EXISTS trigger_update_bid_winning_status ON bids;
DROP TRIGGER IF EXISTS trigger_create_bidder_points ON bidders;
DROP TRIGGER IF EXISTS update_items_updated_at ON items;
DROP TRIGGER IF EXISTS update_auctions_updated_at ON auctions;
DROP TRIGGER IF EXISTS update_bidder_points_updated_at ON bidder_points;
DROP TRIGGER IF EXISTS update_admins_updated_at ON admins;
DROP TRIGGER IF EXISTS update_bidders_updated_at ON bidders;

-- Drop functions
DROP FUNCTION IF EXISTS record_point_history();
DROP FUNCTION IF EXISTS update_bid_winning_status();
DROP FUNCTION IF EXISTS create_bidder_points();
DROP FUNCTION IF EXISTS update_updated_at_column();
