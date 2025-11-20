-- Update admin test accounts with consistent password: password123
-- This migration updates the test accounts to use a consistent password for development
-- IMPORTANT: These accounts MUST NOT be used in production

-- Update system admin account
-- Email: admin@example.com
-- Password: password123 (bcrypt cost 10)
UPDATE admins
SET password_hash = '$2a$10$sqoQigDbroFLmNhNjT3.A.xFku.TJ9et2cGPsRwdsQ0r1f1.HJi0W'
WHERE email = 'admin@example.com';

-- Update auctioneer account
-- Email: auctioneer@example.com
-- Password: password123 (bcrypt cost 10)
UPDATE admins
SET password_hash = '$2a$10$sqoQigDbroFLmNhNjT3.A.xFku.TJ9et2cGPsRwdsQ0r1f1.HJi0W'
WHERE email = 'auctioneer@example.com';
