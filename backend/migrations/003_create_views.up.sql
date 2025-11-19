-- View: active_auctions_view (Active auctions list)
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

-- View: bidder_auction_summary (Bidder auction participation summary)
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

-- View: point_history_detailed (Detailed point history)
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

-- View: bidder_point_balance_view (Bidder point balance)
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

-- View: point_transactions_by_auction (Point transactions by auction)
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
