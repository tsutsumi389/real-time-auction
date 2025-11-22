-- Restore trigger function: Automatically create bidder_points when new bidder is created
CREATE OR REPLACE FUNCTION create_bidder_points()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO bidder_points (bidder_id, total_points, available_points, reserved_points)
    VALUES (NEW.id, 0, 0, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_bidder_points
AFTER INSERT ON bidders
FOR EACH ROW
EXECUTE FUNCTION create_bidder_points();

-- Restore trigger function: Update is_winning status for bids
CREATE OR REPLACE FUNCTION update_bid_winning_status()
RETURNS TRIGGER AS $$
BEGIN
    -- Set all previous bids for this auction to is_winning = FALSE
    UPDATE bids
    SET is_winning = FALSE
    WHERE auction_id = NEW.auction_id AND id != NEW.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_bid_winning_status
AFTER INSERT ON bids
FOR EACH ROW
EXECUTE FUNCTION update_bid_winning_status();

-- Restore trigger function: Automatically record point history
CREATE OR REPLACE FUNCTION record_point_history()
RETURNS TRIGGER AS $$
DECLARE
    point_type VARCHAR(50);
    point_amount BIGINT;
BEGIN
    -- Determine point type and amount based on changes
    IF NEW.total_points > OLD.total_points THEN
        -- Point grant
        point_type := 'grant';
        point_amount := NEW.total_points - OLD.total_points;
    ELSIF NEW.reserved_points > OLD.reserved_points THEN
        -- Reserve points for bidding
        point_type := 'reserve';
        point_amount := -(NEW.reserved_points - OLD.reserved_points);
    ELSIF NEW.reserved_points < OLD.reserved_points AND NEW.available_points > OLD.available_points THEN
        -- Release reserved points
        point_type := 'release';
        point_amount := NEW.available_points - OLD.available_points;
    ELSIF NEW.reserved_points < OLD.reserved_points AND NEW.available_points = OLD.available_points THEN
        -- Consume points on winning
        point_type := 'consume';
        point_amount := -(OLD.reserved_points - NEW.reserved_points);
    ELSE
        -- Other changes (normally should not happen)
        RETURN NEW;
    END IF;

    -- Insert history record
    INSERT INTO point_history (
        bidder_id,
        amount,
        type,
        balance_before,
        balance_after,
        reserved_before,
        reserved_after,
        total_before,
        total_after
    ) VALUES (
        NEW.bidder_id,
        point_amount,
        point_type,
        OLD.available_points,
        NEW.available_points,
        OLD.reserved_points,
        NEW.reserved_points,
        OLD.total_points,
        NEW.total_points
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_record_point_history
AFTER UPDATE ON bidder_points
FOR EACH ROW
WHEN (OLD.* IS DISTINCT FROM NEW.*)
EXECUTE FUNCTION record_point_history();
