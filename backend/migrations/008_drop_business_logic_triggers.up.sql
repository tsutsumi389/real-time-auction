-- Drop trigger and function for automatic bidder_points creation
DROP TRIGGER IF EXISTS trigger_create_bidder_points ON bidders;
DROP FUNCTION IF EXISTS create_bidder_points();

-- Drop trigger and function for bid winning status update
DROP TRIGGER IF EXISTS trigger_update_bid_winning_status ON bids;
DROP FUNCTION IF EXISTS update_bid_winning_status();

-- Drop trigger and function for point history recording
DROP TRIGGER IF EXISTS trigger_record_point_history ON bidder_points;
DROP FUNCTION IF EXISTS record_point_history();
