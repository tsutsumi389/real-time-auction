-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create bidders table
CREATE TABLE bidders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_bidders_status CHECK (status IN ('active', 'suspended', 'deleted')),
    CONSTRAINT chk_bidders_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE UNIQUE INDEX idx_bidders_email ON bidders(email) WHERE status != 'deleted';
CREATE INDEX idx_bidders_status ON bidders(status);

-- Create admins table
CREATE TABLE admins (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'auctioneer',
    display_name VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_admins_role CHECK (role IN ('system_admin', 'auctioneer')),
    CONSTRAINT chk_admins_status CHECK (status IN ('active', 'suspended', 'deleted')),
    CONSTRAINT chk_admins_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE UNIQUE INDEX idx_admins_email ON admins(email) WHERE status != 'deleted';
CREATE INDEX idx_admins_role ON admins(role);
CREATE INDEX idx_admins_status ON admins(status);

-- Create bidder_points table
CREATE TABLE bidder_points (
    bidder_id UUID PRIMARY KEY,
    total_points BIGINT NOT NULL DEFAULT 0,
    available_points BIGINT NOT NULL DEFAULT 0,
    reserved_points BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_bidder_points_bidder FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE,
    CONSTRAINT chk_bidder_points_non_negative CHECK (total_points >= 0 AND available_points >= 0 AND reserved_points >= 0),
    CONSTRAINT chk_bidder_points_balance CHECK (available_points + reserved_points <= total_points)
);

CREATE INDEX idx_bidder_points_available ON bidder_points(available_points);

-- Create auctions table
CREATE TABLE auctions (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    starting_price BIGINT,
    current_price BIGINT,
    winner_id UUID,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_auctions_winner FOREIGN KEY (winner_id) REFERENCES bidders(id) ON DELETE SET NULL,
    CONSTRAINT chk_auctions_status CHECK (status IN ('pending', 'active', 'ended', 'cancelled')),
    CONSTRAINT chk_auctions_price_positive CHECK (starting_price IS NULL OR starting_price > 0),
    CONSTRAINT chk_auctions_dates CHECK (ended_at IS NULL OR started_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_started_at ON auctions(started_at);
CREATE INDEX idx_auctions_ended_at ON auctions(ended_at);
CREATE INDEX idx_auctions_winner ON auctions(winner_id);

-- Create items table
CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_items_auction FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE
);

CREATE INDEX idx_items_auction ON items(auction_id);
CREATE INDEX idx_items_metadata ON items USING GIN(metadata);

-- Create item_media table
CREATE TABLE item_media (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    media_type VARCHAR(20) NOT NULL,
    url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_item_media_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    CONSTRAINT chk_item_media_type CHECK (media_type IN ('image', 'video')),
    CONSTRAINT chk_item_media_order_non_negative CHECK (display_order >= 0)
);

CREATE INDEX idx_item_media_item ON item_media(item_id, display_order);
CREATE INDEX idx_item_media_type ON item_media(media_type);

-- Create bids table
CREATE TABLE bids (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    bidder_id UUID NOT NULL,
    price BIGINT NOT NULL,
    bid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_winning BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT fk_bids_auction FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE,
    CONSTRAINT fk_bids_bidder FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE,
    CONSTRAINT chk_bids_price_positive CHECK (price > 0)
);

CREATE INDEX idx_bids_auction ON bids(auction_id, bid_at DESC);
CREATE INDEX idx_bids_bidder ON bids(bidder_id, bid_at DESC);
CREATE INDEX idx_bids_winning ON bids(auction_id, is_winning) WHERE is_winning = TRUE;
CREATE INDEX idx_bids_auction_bidder ON bids(auction_id, bidder_id);

-- Create price_history table
CREATE TABLE price_history (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    price BIGINT NOT NULL,
    opened_by BIGINT NOT NULL,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    had_bid BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_price_history_auction FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE,
    CONSTRAINT fk_price_history_opener FOREIGN KEY (opened_by) REFERENCES admins(id) ON DELETE CASCADE,
    CONSTRAINT chk_price_history_price_positive CHECK (price > 0)
);

CREATE INDEX idx_price_history_auction ON price_history(auction_id, opened_at DESC);
CREATE INDEX idx_price_history_opened_by ON price_history(opened_by);

-- Create point_history table
CREATE TABLE point_history (
    id BIGSERIAL PRIMARY KEY,
    bidder_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    reason TEXT,
    related_auction_id BIGINT,
    related_bid_id BIGINT,
    admin_id BIGINT,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reserved_before BIGINT NOT NULL,
    reserved_after BIGINT NOT NULL,
    total_before BIGINT NOT NULL,
    total_after BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_point_history_bidder FOREIGN KEY (bidder_id) REFERENCES bidders(id) ON DELETE CASCADE,
    CONSTRAINT fk_point_history_admin FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE SET NULL,
    CONSTRAINT fk_point_history_auction FOREIGN KEY (related_auction_id) REFERENCES auctions(id) ON DELETE SET NULL,
    CONSTRAINT fk_point_history_bid FOREIGN KEY (related_bid_id) REFERENCES bids(id) ON DELETE SET NULL,
    CONSTRAINT chk_point_history_type CHECK (type IN ('grant', 'reserve', 'release', 'consume', 'refund')),
    CONSTRAINT chk_point_history_amount_not_zero CHECK (amount != 0),
    CONSTRAINT chk_point_history_balance_non_negative CHECK (
        balance_before >= 0 AND balance_after >= 0 AND 
        reserved_before >= 0 AND reserved_after >= 0 AND
        total_before >= 0 AND total_after >= 0
    )
);

CREATE INDEX idx_point_history_bidder ON point_history(bidder_id, created_at DESC);
CREATE INDEX idx_point_history_type ON point_history(type);
CREATE INDEX idx_point_history_admin ON point_history(admin_id);
CREATE INDEX idx_point_history_auction ON point_history(related_auction_id);
CREATE INDEX idx_point_history_bid ON point_history(related_bid_id);
CREATE INDEX idx_point_history_created_at ON point_history(created_at DESC);
