-- Rollback admin test accounts to original passwords
-- This reverts the password changes made in the up migration

-- Revert system admin account to original password
-- Password: admin123
UPDATE admins
SET password_hash = '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5lW5vJ4r.Wquy'
WHERE email = 'admin@example.com';

-- Revert auctioneer account to original password
-- Password: auctioneer123
UPDATE admins
SET password_hash = '$2a$12$8kPvKJQNnZqU5vR8gF4Hc.LxK9bE2rW7jT5mC3dP1fN6hG8lA4qBi'
WHERE email = 'auctioneer@example.com';
