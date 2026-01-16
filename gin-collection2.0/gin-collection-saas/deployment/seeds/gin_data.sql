-- Sample Gin Data for GinVault
-- Different gins for each tenant to showcase tier features

-- Get tenant IDs
SET @test_tenant = (SELECT id FROM tenants WHERE subdomain = 'test');
SET @basic_tenant = (SELECT id FROM tenants WHERE subdomain = 'basic');
SET @pro_tenant = (SELECT id FROM tenants WHERE subdomain = 'pro');
SET @enterprise_tenant = (SELECT id FROM tenants WHERE subdomain = 'enterprise');

-- =====================================================
-- FREE TIER (test@test.com) - 3 Gins, basic data
-- =====================================================
INSERT INTO gins (tenant_id, uuid, name, brand, country, gin_type, abv, bottle_size, fill_level, price, rating, created_at, updated_at) VALUES
(@test_tenant, UUID(), 'Bombay Sapphire', 'Bombay', 'England', 'London Dry', 40.0, 700, 75, 24.99, 4, NOW(), NOW()),
(@test_tenant, UUID(), 'Tanqueray', 'Tanqueray', 'England', 'London Dry', 43.1, 700, 50, 22.99, 4, NOW(), NOW()),
(@test_tenant, UUID(), 'Gordon''s', 'Gordon''s', 'England', 'London Dry', 37.5, 700, 90, 14.99, 3, NOW(), NOW());

-- =====================================================
-- BASIC TIER (basic@demo.local) - 5 Gins, more details
-- =====================================================
INSERT INTO gins (tenant_id, uuid, name, brand, country, region, gin_type, abv, bottle_size, fill_level, price, rating, nose_notes, palate_notes, created_at, updated_at) VALUES
(@basic_tenant, UUID(), 'Hendrick''s', 'Hendrick''s', 'Scotland', 'Girvan', 'New Western', 41.4, 700, 60, 34.99, 5, 'Rose petals, cucumber', 'Floral, refreshing', NOW(), NOW()),
(@basic_tenant, UUID(), 'Monkey 47', 'Black Forest', 'Germany', 'Black Forest', 'New Western', 47.0, 500, 40, 44.99, 5, 'Complex, herbal', 'Juniper, citrus, spice', NOW(), NOW()),
(@basic_tenant, UUID(), 'The Botanist', 'Bruichladdich', 'Scotland', 'Islay', 'New Western', 46.0, 700, 85, 39.99, 4, 'Floral, menthol', 'Smooth, botanical', NOW(), NOW()),
(@basic_tenant, UUID(), 'Roku', 'Suntory', 'Japan', 'Osaka', 'Japanese', 43.0, 700, 100, 32.99, 4, 'Sakura, yuzu', 'Delicate, balanced', NOW(), NOW()),
(@basic_tenant, UUID(), 'Sipsmith', 'Sipsmith', 'England', 'London', 'London Dry', 41.6, 700, 30, 29.99, 4, 'Classic juniper', 'Citrus forward', NOW(), NOW());

-- =====================================================
-- PRO TIER (pro@demo.local) - 8 Gins, full tasting notes
-- =====================================================
INSERT INTO gins (tenant_id, uuid, name, brand, country, region, gin_type, abv, bottle_size, fill_level, price, current_market_value, rating, nose_notes, palate_notes, finish_notes, general_notes, recommended_tonic, recommended_garnish, created_at, updated_at) VALUES
(@pro_tenant, UUID(), 'Nordes Atlantic Galician', 'Nordes', 'Spain', 'Galicia', 'New Western', 40.0, 700, 70, 38.99, 42.00, 5, 'White grape, eucalyptus', 'Herbal, fresh', 'Long, aromatic', 'Perfect for summer', 'Fever-Tree Mediterranean', 'Grape slice', NOW(), NOW()),
(@pro_tenant, UUID(), 'Ki No Bi', 'Kyoto Distillery', 'Japan', 'Kyoto', 'Japanese', 45.7, 700, 55, 59.99, 65.00, 5, 'Hinoki, yuzu, sansho', 'Complex, layered', 'Warm, spicy', 'Artisanal Japanese craft', 'Fever-Tree Yuzu', 'Shiso leaf', NOW(), NOW()),
(@pro_tenant, UUID(), 'Ferdinand''s Saar', 'Ferdinand''s', 'Germany', 'Saar', 'New Western', 44.0, 500, 45, 42.99, 48.00, 5, 'Riesling, lavender', 'Elegant, wine-like', 'Smooth, floral', 'Infused with Riesling grapes', 'Schweppes Dry', 'Lemon twist', NOW(), NOW()),
(@pro_tenant, UUID(), 'Malfy Rosa', 'Malfy', 'Italy', 'Moncalieri', 'Pink Gin', 41.0, 700, 80, 27.99, 29.00, 4, 'Pink grapefruit, rhubarb', 'Sweet, citrusy', 'Refreshing', 'Beautiful pink color', 'Fever-Tree Aromatic', 'Pink grapefruit', NOW(), NOW()),
(@pro_tenant, UUID(), 'Tarquin''s Cornish', 'Tarquin''s', 'England', 'Cornwall', 'London Dry', 42.0, 700, 95, 36.99, 38.00, 4, 'Violet, citrus', 'Floral, juniper', 'Crisp, clean', 'Handcrafted in small batches', 'Fever-Tree Indian', 'Orange peel', NOW(), NOW()),
(@pro_tenant, UUID(), 'Four Pillars Bloody Shiraz', 'Four Pillars', 'Australia', 'Yarra Valley', 'Sloe/Fruit', 37.8, 700, 25, 54.99, 58.00, 5, 'Shiraz grapes, spice', 'Rich, berry', 'Long, wine-like', 'Limited edition', 'Plain soda', 'Orange slice', NOW(), NOW()),
(@pro_tenant, UUID(), 'Aviation American', 'Aviation', 'USA', 'Portland', 'New Western', 42.0, 700, 65, 31.99, 34.00, 4, 'Lavender, anise', 'Smooth, floral', 'Subtle, dry', 'Ryan Reynolds'' gin', 'Fever-Tree Elderflower', 'Lavender sprig', NOW(), NOW()),
(@pro_tenant, UUID(), 'Whitley Neill Rhubarb & Ginger', 'Whitley Neill', 'England', 'Birmingham', 'Flavored', 43.0, 700, 100, 26.99, 28.00, 4, 'Rhubarb, ginger warmth', 'Sweet, tangy', 'Spicy finish', 'Great for cocktails', 'Ginger ale', 'Rhubarb ribbon', NOW(), NOW());

-- =====================================================
-- ENTERPRISE TIER (enterprise@demo.local) - 12 Premium Gins
-- =====================================================
INSERT INTO gins (tenant_id, uuid, name, brand, country, region, gin_type, abv, bottle_size, fill_level, price, current_market_value, purchase_date, purchase_location, rating, nose_notes, palate_notes, finish_notes, general_notes, recommended_tonic, recommended_garnish, is_finished, created_at, updated_at) VALUES
(@enterprise_tenant, UUID(), 'Nolet''s Silver', 'Nolet', 'Netherlands', 'Schiedam', 'New Western', 47.6, 700, 50, 49.99, 55.00, '2024-06-15', 'Whisky Exchange London', 5, 'Turkish rose, peach', 'Silky, floral', 'Long, elegant', 'Ultra-premium Dutch gin', 'Fever-Tree Premium Indian', 'Rose petal', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Gin Mare', 'Gin Mare', 'Spain', 'Costa Brava', 'Mediterranean', 42.7, 700, 85, 44.99, 48.00, '2024-08-20', 'El Corte Inglés Barcelona', 5, 'Olive, basil, rosemary', 'Herbal, savory', 'Mediterranean warmth', 'Perfect for Martini', 'Fever-Tree Mediterranean', 'Rosemary sprig', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Citadelle Réserve', 'Citadelle', 'France', 'Cognac', 'Barrel Aged', 44.0, 700, 40, 54.99, 62.00, '2024-03-10', 'La Maison du Whisky Paris', 5, 'Oak, vanilla, spice', 'Rich, complex', 'Warm, woody', 'Aged in Cognac barrels', 'Schweppes 1783', 'Orange zest', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Hernö Old Tom', 'Hernö', 'Sweden', 'Härnösand', 'Old Tom', 43.0, 500, 70, 46.99, 52.00, '2024-05-22', 'Systembolaget Stockholm', 5, 'Meadowsweet, juniper', 'Sweet, balanced', 'Lingering botanicals', 'World''s best Old Tom', 'Fever-Tree Aromatic', 'Lemon twist', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Uncle Val''s Botanical', 'Uncle Val''s', 'USA', 'California', 'New Western', 45.0, 750, 60, 34.99, 38.00, '2024-07-08', 'Total Wine San Francisco', 4, 'Cucumber, sage', 'Crisp, garden fresh', 'Clean, herbal', 'Farm to bottle concept', 'Q Tonic', 'Cucumber ribbon', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Ableforth''s Bathtub', 'Ableforth''s', 'England', 'Tunbridge Wells', 'Compound', 43.3, 700, 90, 36.99, 40.00, '2024-09-14', 'Master of Malt', 4, 'Warm spices, juniper', 'Cardamom, cinnamon', 'Long, spiced', 'Cold compounded gin', 'Fever-Tree Indian', 'Orange peel', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Drumshanbo Gunpowder', 'Drumshanbo', 'Ireland', 'Leitrim', 'Irish', 43.0, 700, 35, 38.99, 42.00, '2024-04-30', 'Celtic Whiskey Shop Dublin', 4, 'Gunpowder tea, citrus', 'Smooth, oriental', 'Tea-like finish', 'Chinese gunpowder tea infused', 'Fever-Tree Indian', 'Grapefruit wedge', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Bobby''s Schiedam', 'Bobby''s', 'Netherlands', 'Schiedam', 'Genever Style', 42.0, 700, 100, 39.99, 44.00, '2024-10-05', 'Gall & Gall Amsterdam', 4, 'Lemongrass, clove', 'Indonesian spices', 'Exotic, warming', 'Indonesian-Dutch fusion', 'Fever-Tree Ginger Ale', 'Lemongrass stalk', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Oxley Cold Distilled', 'Oxley', 'England', 'London', 'London Dry', 47.0, 700, 55, 52.99, 58.00, '2024-02-18', 'Hedonism Wines London', 5, 'Fresh citrus, juniper', 'Vibrant, pure', 'Exceptionally clean', 'Cold vacuum distilled', 'Fever-Tree Premium Indian', 'Grapefruit twist', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'St. George Terroir', 'St. George Spirits', 'USA', 'California', 'New Western', 45.0, 750, 45, 38.99, 44.00, '2024-08-01', 'K&L Wine Merchants', 4, 'Douglas fir, sage', 'Forest floor, herbs', 'Pine, earthy', 'Tastes like a hike in California', 'Fever-Tree Elderflower', 'Rosemary', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Cambridge Dry', 'Cambridge Distillery', 'England', 'Cambridge', 'London Dry', 42.0, 700, 20, 44.99, 50.00, '2024-01-25', 'Cambridge Wine Merchants', 5, 'Rose, lemon verbena', 'Elegant, refined', 'Delicate, floral', 'Tailored gin concept', 'Fever-Tree Naturally Light', 'Edible flowers', 0, NOW(), NOW()),
(@enterprise_tenant, UUID(), 'Beefeater 24', 'Beefeater', 'England', 'London', 'London Dry', 45.0, 700, 0, 34.99, 38.00, '2023-11-20', 'The Whisky Exchange', 4, 'Japanese sencha tea, citrus', 'Smooth, tea notes', 'Long, refined', 'FINISHED - need to reorder!', 'Fever-Tree Indian', 'Grapefruit peel', 1, NOW(), NOW());

-- Verify inserted data
SELECT
    t.subdomain as tenant,
    t.tier,
    COUNT(g.id) as gin_count,
    ROUND(AVG(g.rating), 1) as avg_rating,
    ROUND(SUM(g.price), 2) as total_value
FROM tenants t
LEFT JOIN gins g ON t.id = g.tenant_id
WHERE t.subdomain IN ('test', 'basic', 'pro', 'enterprise')
GROUP BY t.id, t.subdomain, t.tier
ORDER BY t.id;
