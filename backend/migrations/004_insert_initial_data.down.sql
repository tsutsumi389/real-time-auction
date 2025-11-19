-- Delete test data
DELETE FROM bidders WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333'
);

DELETE FROM admins WHERE email IN (
    'admin@example.com',
    'auctioneer@example.com'
);
