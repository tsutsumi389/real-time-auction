-- Insert system administrator account
-- Password: admin123 (MUST be changed in production)
-- Generated using bcrypt cost factor 12
INSERT INTO admins (email, password_hash, role, display_name, status)
VALUES (
    'admin@example.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5lW5vJ4r.Wquy',
    'system_admin',
    'システム管理者',
    'active'
);

-- Insert auctioneer account
-- Password: auctioneer123
INSERT INTO admins (email, password_hash, role, display_name, status)
VALUES (
    'auctioneer@example.com',
    '$2a$12$8kPvKJQNnZqU5vR8gF4Hc.LxK9bE2rW7jT5mC3dP1fN6hG8lA4qBi',
    'auctioneer',
    '主催者',
    'active'
);

-- Insert test bidder accounts
-- Password: bidder123 for all
INSERT INTO bidders (id, email, password_hash, display_name, status)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'bidder1@example.com', '$2a$12$9kQwLMRPqZrV6wS9hG5Id.MyL0cF3sX8kU6nD4eQ2gO7iH9mB5rCj', '入札者1', 'active'),
    ('22222222-2222-2222-2222-222222222222', 'bidder2@example.com', '$2a$12$0lRxMNSQraW7xT0iH6Je.NzM1dG4tY9lV7oE5fR3hP8jI0nC6sDk', '入札者2', 'active'),
    ('33333333-3333-3333-3333-333333333333', 'bidder3@example.com', '$2a$12$1mSyNOTRsbX8yU1jI7Kf.O0N2eH5uZ0mW8pF6gS4iQ9kJ1oD7tEl', '入札者3', 'active');

-- Grant test points to bidders (10,000 points each)
-- Note: bidder_points records are automatically created by trigger
UPDATE bidder_points 
SET total_points = 10000, available_points = 10000 
WHERE bidder_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333'
);
