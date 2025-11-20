-- Add suspended admin account for testing
-- Email: suspended@example.com
-- Password: password123 (bcrypt cost 10)
INSERT INTO admins (email, password_hash, role, display_name, status)
VALUES (
    'suspended@example.com',
    '$2a$10$sqoQigDbroFLmNhNjT3.A.xFku.TJ9et2cGPsRwdsQ0r1f1.HJi0W',
    'auctioneer',
    '停止中のアカウント',
    'suspended'
);
