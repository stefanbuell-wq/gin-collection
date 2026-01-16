-- Test Data for GinVault
-- Password for all test users: TEST123456
-- Hash: $2b$10$dEzAIfnOwTf26../Sj4GgOR8S1A076/VDz2xX4i.LHRy1ZyBbFA32

-- Clear existing test data (optional)
-- DELETE FROM users WHERE email IN ('test@test.com', 'basic@demo.local', 'pro@demo.local', 'enterprise@demo.local');
-- DELETE FROM tenants WHERE subdomain IN ('test', 'basic', 'pro', 'enterprise');

-- Create Tenants
INSERT INTO tenants (uuid, name, subdomain, tier, is_active, created_at, updated_at) VALUES
(UUID(), 'Test Account', 'test', 'free', 1, NOW(), NOW()),
(UUID(), 'Basic Demo', 'basic', 'basic', 1, NOW(), NOW()),
(UUID(), 'Pro Demo', 'pro', 'pro', 1, NOW(), NOW()),
(UUID(), 'Enterprise Demo', 'enterprise', 'enterprise', 1, NOW(), NOW());

-- Create Users (owner role for each tenant)
INSERT INTO users (tenant_id, uuid, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at) VALUES
((SELECT id FROM tenants WHERE subdomain='test'), UUID(), 'test@test.com', '$2b$10$dEzAIfnOwTf26../Sj4GgOR8S1A076/VDz2xX4i.LHRy1ZyBbFA32', 'Test', 'User', 'owner', 1, NOW(), NOW()),
((SELECT id FROM tenants WHERE subdomain='basic'), UUID(), 'basic@demo.local', '$2b$10$dEzAIfnOwTf26../Sj4GgOR8S1A076/VDz2xX4i.LHRy1ZyBbFA32', 'Basic', 'Demo', 'owner', 1, NOW(), NOW()),
((SELECT id FROM tenants WHERE subdomain='pro'), UUID(), 'pro@demo.local', '$2b$10$dEzAIfnOwTf26../Sj4GgOR8S1A076/VDz2xX4i.LHRy1ZyBbFA32', 'Pro', 'Demo', 'owner', 1, NOW(), NOW()),
((SELECT id FROM tenants WHERE subdomain='enterprise'), UUID(), 'enterprise@demo.local', '$2b$10$dEzAIfnOwTf26../Sj4GgOR8S1A076/VDz2xX4i.LHRy1ZyBbFA32', 'Enterprise', 'Demo', 'owner', 1, NOW(), NOW());

-- Verify
SELECT t.subdomain, t.tier, u.email, u.role FROM tenants t JOIN users u ON t.id = u.tenant_id ORDER BY t.id;
