-- Gin Reference Catalog
-- A comprehensive list of popular gins for users to quickly add to their collection

-- Create reference table if not exists
CREATE TABLE IF NOT EXISTS gin_references (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    country VARCHAR(100),
    region VARCHAR(100),
    gin_type VARCHAR(50),
    abv DECIMAL(4,1),
    bottle_size INT DEFAULT 700,
    description TEXT,
    nose_notes TEXT,
    palate_notes TEXT,
    finish_notes TEXT,
    recommended_tonic VARCHAR(255),
    recommended_garnish VARCHAR(255),
    image_url VARCHAR(512),
    barcode VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_gin (name, brand),
    INDEX idx_gin_ref_name (name),
    INDEX idx_gin_ref_brand (brand),
    INDEX idx_gin_ref_country (country),
    INDEX idx_gin_ref_type (gin_type)
);

-- Clear existing reference data
TRUNCATE TABLE gin_references;

-- =====================================================
-- LONDON DRY GINS
-- =====================================================
INSERT INTO gin_references (name, brand, country, region, gin_type, abv, bottle_size, description, nose_notes, palate_notes, finish_notes, recommended_tonic, recommended_garnish) VALUES
('Bombay Sapphire', 'Bombay', 'England', 'Hampshire', 'London Dry', 40.0, 700, 'Iconic blue bottle, 10 hand-selected botanicals', 'Juniper, citrus peel, coriander', 'Balanced, peppery, light citrus', 'Clean, dry finish', 'Fever-Tree Indian', 'Lemon wedge'),
('Tanqueray', 'Tanqueray', 'England', 'Cameronbridge', 'London Dry', 43.1, 700, 'Classic 4-times distilled gin since 1830', 'Bold juniper, citrus, hint of pepper', 'Crisp, juniper-forward, citrus', 'Long, dry, peppery', 'Schweppes', 'Lime wedge'),
('Tanqueray No. Ten', 'Tanqueray', 'England', 'Cameronbridge', 'London Dry', 47.3, 700, 'Small batch, fresh citrus, Tiny Ten still', 'Fresh grapefruit, chamomile, juniper', 'Creamy, citrus-forward, complex', 'Long, smooth, citrus', 'Fever-Tree Mediterranean', 'Grapefruit slice'),
('Gordon''s', 'Gordon''s', 'England', 'Cameronbridge', 'London Dry', 37.5, 700, 'World''s best-selling London Dry since 1769', 'Juniper, coriander, angelica', 'Classic, juniper-dominant', 'Clean, dry', 'Schweppes', 'Lime wedge'),
('Beefeater', 'Beefeater', 'England', 'London', 'London Dry', 40.0, 700, 'Only premium gin still distilled in London', 'Juniper, citrus, angelica', 'Balanced, citrus notes, juniper', 'Dry, lingering', 'Fever-Tree Indian', 'Lemon twist'),
('Beefeater 24', 'Beefeater', 'England', 'London', 'London Dry', 45.0, 700, 'Japanese sencha and Chinese green tea infused', 'Japanese sencha tea, grapefruit, juniper', 'Smooth, tea notes, citrus', 'Long, refined', 'Fever-Tree Indian', 'Grapefruit peel'),
('Plymouth Gin', 'Plymouth', 'England', 'Plymouth', 'Plymouth', 41.2, 700, 'Only gin with Protected Geographical Indication', 'Earthy, juniper, cardamom', 'Full-bodied, earthy, slight sweetness', 'Smooth, rounded', 'Fever-Tree Mediterranean', 'Orange peel'),
('Sipsmith', 'Sipsmith', 'England', 'London', 'London Dry', 41.6, 700, 'First copper pot distillery in London since 1820', 'Classic juniper, citrus, floral', 'Dry, citrus-forward, bold juniper', 'Crisp, clean', 'Fever-Tree Indian', 'Lemon twist'),
('Broker''s', 'Broker''s', 'England', 'Birmingham', 'London Dry', 40.0, 700, 'Recognizable bowler hat cap, 10 botanicals', 'Juniper, citrus, spice', 'Dry, juniper-heavy, peppery', 'Long, spicy finish', 'Fever-Tree Indian', 'Lemon wedge'),
('Hayman''s London Dry', 'Hayman''s', 'England', 'London', 'London Dry', 41.2, 700, 'Family recipe since 1863, true English gin', 'Juniper, coriander, citrus peel', 'Full, juniper, citrus', 'Dry, balanced', 'Fever-Tree Indian', 'Lemon peel'),
('Bulldog', 'Bulldog', 'England', 'London', 'London Dry', 40.0, 700, 'Bold British attitude, dragon eye, poppy', 'Citrus, floral, juniper', 'Smooth, floral, hint of spice', 'Clean, citrus finish', 'Fever-Tree Indian', 'Grapefruit wedge'),
('Portobello Road', 'Portobello Road', 'England', 'London', 'London Dry', 42.0, 700, 'From the famous Notting Hill market', 'Citrus, juniper, nutmeg', 'Balanced, spicy, citrus', 'Long, warm, spicy', 'Fever-Tree Indian', 'Lemon twist'),
('City of London', 'City of London', 'England', 'London', 'London Dry', 41.3, 700, 'First new distillery in the City since 1820', 'Orange blossom, juniper, coriander', 'Citrus-forward, juniper, floral', 'Smooth, citrus finish', 'Fever-Tree Indian', 'Orange peel'),
('No. 3 London Dry', 'Berry Bros & Rudd', 'Netherlands', 'Schiedam', 'London Dry', 46.0, 700, 'From Britain''s oldest wine merchant', 'Juniper, citrus, cardamom', 'Bold, juniper, complex spices', 'Long, warming', 'Fever-Tree Indian', 'Grapefruit twist'),

-- =====================================================
-- NEW WESTERN / CONTEMPORARY GINS
-- =====================================================
('Hendrick''s', 'Hendrick''s', 'Scotland', 'Girvan', 'New Western', 41.4, 700, 'Infused with rose petals and cucumber', 'Rose petals, cucumber, citrus', 'Floral, cucumber, subtle juniper', 'Refreshing, smooth', 'Fever-Tree Elderflower', 'Cucumber slice'),
('Hendrick''s Orbium', 'Hendrick''s', 'Scotland', 'Girvan', 'New Western', 43.4, 700, 'Quinine, wormwood, blue lotus blossom', 'Quinine, floral, herbaceous', 'Bitter, complex, botanical', 'Long, bitter finish', 'Plain soda', 'Cucumber'),
('Monkey 47', 'Black Forest Distillers', 'Germany', 'Black Forest', 'New Western', 47.0, 500, '47 botanicals from the Black Forest', 'Complex herbal, pine, citrus', 'Juniper, citrus, lingonberry, spice', 'Long, complex, herbal', 'Thomas Henry', 'Lemon twist'),
('The Botanist', 'Bruichladdich', 'Scotland', 'Islay', 'New Western', 46.0, 700, '22 hand-foraged Islay botanicals', 'Floral, menthol, honey', 'Smooth, floral, wild herbs', 'Long, minty, herbal', 'Fever-Tree Mediterranean', 'Apple slice'),
('Aviation', 'Aviation', 'USA', 'Portland', 'New Western', 42.0, 700, 'Ryan Reynolds'' gin, lavender-forward', 'Lavender, cardamom, anise', 'Smooth, floral, spice', 'Subtle, dry', 'Fever-Tree Elderflower', 'Lavender sprig'),
('Uncle Val''s Botanical', 'Uncle Val''s', 'USA', 'California', 'New Western', 45.0, 750, 'Farm-to-bottle, cucumber and sage', 'Cucumber, sage, lemon', 'Crisp, garden fresh, herbal', 'Clean, herbal finish', 'Q Tonic', 'Cucumber ribbon'),
('St. George Terroir', 'St. George Spirits', 'USA', 'California', 'New Western', 45.0, 750, 'Tastes like a California forest hike', 'Douglas fir, sage, bay laurel', 'Forest floor, pine, herbs', 'Earthy, pine finish', 'Fever-Tree Elderflower', 'Rosemary sprig'),

-- =====================================================
-- JAPANESE GINS
-- =====================================================
('Roku', 'Suntory', 'Japan', 'Osaka', 'Japanese', 43.0, 700, '6 Japanese botanicals, hexagonal bottle', 'Sakura, yuzu, sencha tea', 'Delicate, balanced, citrus', 'Smooth, tea notes', 'Fever-Tree Japanese Yuzu', 'Ginger slice'),
('Ki No Bi', 'Kyoto Distillery', 'Japan', 'Kyoto', 'Japanese', 45.7, 700, 'Rice spirit base, 11 Japanese botanicals', 'Hinoki, yuzu, sansho pepper', 'Complex, layered, citrus', 'Warm, spicy finish', 'Fever-Tree Mediterranean', 'Shiso leaf'),
('Nikka Coffey Gin', 'Nikka', 'Japan', 'Sendai', 'Japanese', 47.0, 700, 'Made in Coffey stills, citrus-forward', 'Yuzu, kabosu, amanatsu', 'Bright citrus, juniper, sansho', 'Citrus, peppery', 'Fever-Tree Mediterranean', 'Yuzu peel'),
('Etsu', 'Etsu', 'Japan', 'Hokkaido', 'Japanese', 43.0, 700, 'Hokkaido botanicals, yuzu and green tea', 'Yuzu, green tea, juniper', 'Bright, citrus, tea notes', 'Clean, refreshing', 'Fever-Tree Indian', 'Yuzu slice'),

-- =====================================================
-- SPANISH & MEDITERRANEAN GINS
-- =====================================================
('Gin Mare', 'Gin Mare', 'Spain', 'Costa Brava', 'Mediterranean', 42.7, 700, 'Mediterranean botanicals, olive, basil', 'Olive, basil, rosemary, thyme', 'Herbal, savory, citrus', 'Mediterranean warmth', 'Fever-Tree Mediterranean', 'Rosemary sprig'),
('Nordés Atlantic Galician', 'Nordés', 'Spain', 'Galicia', 'New Western', 40.0, 700, 'Albariño grape base, Atlantic botanicals', 'White grape, eucalyptus, laurel', 'Herbal, fresh, fruity', 'Long, aromatic finish', 'Fever-Tree Mediterranean', 'Grape slice'),
('Malfy Originale', 'Malfy', 'Italy', 'Moncalieri', 'Italian', 41.0, 700, 'Italian coastal juniper, sun-ripened citrus', 'Italian lemon, juniper', 'Bright citrus, juniper', 'Clean, citrus finish', 'Fever-Tree Mediterranean', 'Lemon wheel'),
('Malfy Con Limone', 'Malfy', 'Italy', 'Moncalieri', 'Italian', 41.0, 700, 'Amalfi Coast lemons, vibrant citrus', 'Amalfi lemon, juniper, coriander', 'Bold lemon, smooth juniper', 'Zesty, refreshing', 'Fever-Tree Mediterranean', 'Lemon wheel'),
('Malfy Rosa', 'Malfy', 'Italy', 'Moncalieri', 'Pink Gin', 41.0, 700, 'Sicilian pink grapefruit and rhubarb', 'Pink grapefruit, rhubarb', 'Sweet, citrusy, fruity', 'Refreshing, fruity', 'Fever-Tree Aromatic', 'Grapefruit slice'),
('Malfy Con Arancia', 'Malfy', 'Italy', 'Moncalieri', 'Italian', 41.0, 700, 'Sicilian blood oranges', 'Blood orange, juniper', 'Sweet orange, balanced juniper', 'Citrus, warming', 'Fever-Tree Mediterranean', 'Orange slice'),

-- =====================================================
-- GERMAN GINS
-- =====================================================
('Ferdinand''s Saar', 'Ferdinand''s', 'Germany', 'Saar', 'New Western', 44.0, 500, 'Riesling grape infusion, 30 botanicals', 'Riesling, lavender, rose', 'Elegant, wine-like, floral', 'Smooth, floral finish', 'Schweppes Dry', 'Lemon twist'),
('Siegfried Rheinland', 'Siegfried', 'Germany', 'Rhineland', 'London Dry', 41.0, 500, 'Inspired by the Nibelungen saga', 'Linden blossom, thyme, juniper', 'Herbal, floral, balanced', 'Dry, herbal finish', 'Thomas Henry', 'Thyme sprig'),
('Windspiel Premium Dry', 'Windspiel', 'Germany', 'Eifel', 'London Dry', 47.0, 500, 'Volcanic Eifel region, potato base', 'Juniper, lavender, lemon', 'Smooth, lavender, citrus', 'Long, floral', 'Fever-Tree Mediterranean', 'Lemon twist'),

-- =====================================================
-- DUTCH GINS
-- =====================================================
('Nolet''s Silver', 'Nolet', 'Netherlands', 'Schiedam', 'New Western', 47.6, 700, '11th generation family recipe', 'Turkish rose, peach, raspberry', 'Silky, floral, fruit', 'Long, elegant', 'Fever-Tree Premium Indian', 'Rose petal'),
('Bobby''s Schiedam', 'Bobby''s', 'Netherlands', 'Schiedam', 'Genever Style', 42.0, 700, 'Indonesian-Dutch fusion, lemongrass', 'Lemongrass, clove, cubeb', 'Indonesian spices, citrus', 'Exotic, warming', 'Fever-Tree Ginger Ale', 'Lemongrass'),
('Rutte Celery', 'Rutte', 'Netherlands', 'Dordrecht', 'New Western', 43.0, 700, '1872 recipe, celery-forward', 'Celery, juniper, cardamom', 'Vegetal, herbal, juniper', 'Celery, dry finish', 'Fever-Tree Mediterranean', 'Celery stick'),

-- =====================================================
-- SCOTTISH GINS
-- =====================================================
('Edinburgh Gin', 'Edinburgh Gin', 'Scotland', 'Edinburgh', 'London Dry', 43.0, 700, 'Capital city gin, Scottish botanicals', 'Juniper, citrus, Scottish herbs', 'Balanced, citrus, herbal', 'Clean, dry', 'Fever-Tree Indian', 'Orange peel'),
('Edinburgh Seaside', 'Edinburgh Gin', 'Scotland', 'Edinburgh', 'New Western', 43.0, 700, 'Coastal botanicals, scurvy grass', 'Sea buckthorn, coastal herbs', 'Saline, herbal, citrus', 'Maritime, fresh', 'Fever-Tree Mediterranean', 'Lemon twist'),
('Rock Rose', 'Dunnet Bay', 'Scotland', 'Caithness', 'New Western', 41.5, 700, 'Hand-foraged Caithness botanicals', 'Rhodiola rosea, rowan berry', 'Floral, fruity, juniper', 'Long, aromatic', 'Fever-Tree Mediterranean', 'Orange peel'),
('Caorunn', 'Caorunn', 'Scotland', 'Speyside', 'New Western', 41.8, 700, '5 Celtic botanicals, rowan berry', 'Rowan berry, apple, heather', 'Crisp apple, floral', 'Dry, clean', 'Fever-Tree Indian', 'Apple slice'),
('Harris Gin', 'Isle of Harris', 'Scotland', 'Outer Hebrides', 'New Western', 45.0, 700, 'Sugar kelp from Outer Hebrides', 'Sea kelp, juniper, citrus', 'Maritime, herbal, citrus', 'Long, saline', 'Fever-Tree Mediterranean', 'Grapefruit'),

-- =====================================================
-- IRISH GINS
-- =====================================================
('Drumshanbo Gunpowder', 'Drumshanbo', 'Ireland', 'Leitrim', 'Irish', 43.0, 700, 'Chinese gunpowder tea infused', 'Gunpowder tea, citrus, juniper', 'Smooth, oriental, spice', 'Tea-like, dry', 'Fever-Tree Indian', 'Grapefruit wedge'),
('Dingle Original', 'Dingle', 'Ireland', 'Kerry', 'London Dry', 42.5, 700, 'From Ireland''s first whiskey distillery', 'Rowan berry, bog myrtle, heather', 'Floral, herbal, juniper', 'Dry, herbal', 'Fever-Tree Mediterranean', 'Rosemary'),
('Mór Irish Gin', 'Mór', 'Ireland', 'Galway', 'Irish', 40.0, 700, 'Wild Atlantic botanicals', 'Honey, heather, wild flowers', 'Floral, honey, citrus', 'Smooth, sweet', 'Fever-Tree Mediterranean', 'Lemon twist'),

-- =====================================================
-- FRENCH GINS
-- =====================================================
('Citadelle', 'Citadelle', 'France', 'Cognac', 'French', 44.0, 700, '19 botanicals, Cognac region', 'Juniper, violet, cinnamon', 'Complex, spiced, floral', 'Warming, spicy', 'Fever-Tree Indian', 'Lemon peel'),
('Citadelle Réserve', 'Citadelle', 'France', 'Cognac', 'Barrel Aged', 44.0, 700, 'Aged in Cognac barrels', 'Oak, vanilla, juniper', 'Rich, complex, woody', 'Warm, oaky finish', 'Schweppes 1783', 'Orange zest'),
('G''Vine Floraison', 'G''Vine', 'France', 'Cognac', 'French', 40.0, 700, 'Grape spirit, vine flower infusion', 'Vine flower, citrus, ginger', 'Floral, fresh, soft', 'Elegant, smooth', 'Fever-Tree Mediterranean', 'Lime twist'),
('Grey Goose Originale', 'Grey Goose', 'France', 'Cognac', 'French', 40.0, 700, 'From the makers of the famous vodka', 'Almond, juniper, citrus', 'Smooth, nutty, citrus', 'Clean, dry', 'Fever-Tree Indian', 'Lemon twist'),

-- =====================================================
-- AUSTRALIAN GINS
-- =====================================================
('Four Pillars Rare Dry', 'Four Pillars', 'Australia', 'Yarra Valley', 'New Western', 41.8, 700, 'Native Australian botanicals', 'Tasmanian pepper, lemon myrtle', 'Citrus, spice, pepper', 'Long, spicy', 'Fever-Tree Indian', 'Orange slice'),
('Four Pillars Bloody Shiraz', 'Four Pillars', 'Australia', 'Yarra Valley', 'Sloe/Fruit', 37.8, 700, 'Steeped in Yarra Valley Shiraz', 'Shiraz grapes, spice, berry', 'Rich, berry, wine-like', 'Long, fruity', 'Plain soda', 'Orange slice'),
('Archie Rose White Rye', 'Archie Rose', 'Australia', 'Sydney', 'New Western', 40.0, 700, 'White rye base, native botanicals', 'Blood lime, pepper berry', 'Spicy, citrus, complex', 'Long, peppery', 'Fever-Tree Indian', 'Lime wheel'),
('Adelaide Hills 78 Degrees', '78 Degrees', 'Australia', 'Adelaide Hills', 'London Dry', 42.0, 700, 'Cool climate botanicals', 'Juniper, citrus, pepperberry', 'Balanced, peppery, citrus', 'Dry, spicy', 'Fever-Tree Indian', 'Lemon twist'),

-- =====================================================
-- AMERICAN GINS
-- =====================================================
('Bluecoat', 'Bluecoat', 'USA', 'Philadelphia', 'American', 47.0, 750, 'American citrus-forward gin', 'Citrus, juniper, organic botanicals', 'Bright citrus, juniper', 'Clean, citrus finish', 'Fever-Tree Indian', 'Lemon peel'),
('Death''s Door', 'Death''s Door', 'USA', 'Wisconsin', 'American', 47.0, 750, 'Only 3 botanicals: juniper, fennel, coriander', 'Juniper, fennel, coriander', 'Bold juniper, anise', 'Dry, herbal', 'Fever-Tree Indian', 'Fennel frond'),
('Junipero', 'Anchor Distilling', 'USA', 'San Francisco', 'American', 49.3, 750, 'America''s original craft gin since 1996', 'Bold juniper, citrus, spice', 'Intense juniper, complex', 'Long, dry, warming', 'Fever-Tree Indian', 'Lemon twist'),
('FEW American Gin', 'FEW Spirits', 'USA', 'Illinois', 'American', 46.5, 750, 'Grain-to-glass production', 'Vanilla, citrus, hops', 'Vanilla, citrus, herbal', 'Smooth, hoppy finish', 'Fever-Tree Indian', 'Lemon peel'),

-- =====================================================
-- PINK / FLAVORED GINS
-- =====================================================
('Gordon''s Pink', 'Gordon''s', 'England', 'Cameronbridge', 'Pink Gin', 37.5, 700, 'Raspberry, strawberry, redcurrant', 'Berry, floral, juniper', 'Sweet berry, juniper', 'Fruity, refreshing', 'Schweppes Slimline', 'Strawberry'),
('Beefeater Pink', 'Beefeater', 'England', 'London', 'Pink Gin', 37.5, 700, 'Natural strawberry flavor', 'Strawberry, juniper', 'Sweet strawberry, juniper', 'Fruity, dry', 'Lemonade', 'Strawberry'),
('Whitley Neill Rhubarb & Ginger', 'Whitley Neill', 'England', 'Birmingham', 'Flavored', 43.0, 700, 'English rhubarb and ginger essence', 'Rhubarb, ginger warmth', 'Sweet, tangy, spicy', 'Warming, ginger finish', 'Ginger ale', 'Rhubarb ribbon'),
('Whitley Neill Blood Orange', 'Whitley Neill', 'England', 'Birmingham', 'Flavored', 43.0, 700, 'Sicilian blood orange essence', 'Blood orange, citrus', 'Sweet orange, juniper', 'Citrus, refreshing', 'Fever-Tree Mediterranean', 'Orange slice'),
('Tanqueray Flor de Sevilla', 'Tanqueray', 'England', 'Cameronbridge', 'Flavored', 41.3, 700, 'Seville orange essence', 'Seville orange, vanilla', 'Bittersweet orange, juniper', 'Orange marmalade', 'Fever-Tree Mediterranean', 'Orange wheel'),
('Tanqueray Rangpur', 'Tanqueray', 'England', 'Cameronbridge', 'Flavored', 41.3, 700, 'Rangpur lime from India', 'Rangpur lime, juniper', 'Zesty lime, juniper', 'Citrus, dry', 'Fever-Tree Indian', 'Lime wheel'),

-- =====================================================
-- OLD TOM & GENEVER GINS
-- =====================================================
('Hayman''s Old Tom', 'Hayman''s', 'England', 'London', 'Old Tom', 41.4, 700, 'Victorian-era recipe, slightly sweet', 'Juniper, citrus, sweetness', 'Sweet, botanical, juniper', 'Smooth, sweet finish', 'Fever-Tree Indian', 'Lemon twist'),
('Hernö Old Tom', 'Hernö', 'Sweden', 'Härnösand', 'Old Tom', 43.0, 500, 'World''s best Old Tom, meadowsweet', 'Meadowsweet, juniper, cassia', 'Sweet, balanced, floral', 'Lingering, botanical', 'Fever-Tree Aromatic', 'Lemon twist'),
('Jensen''s Old Tom', 'Jensen''s', 'England', 'London', 'Old Tom', 43.0, 700, 'Recreation of 1840s recipe', 'Liquorice, juniper, citrus', 'Sweet, herbal, juniper', 'Smooth, sweet', 'Fever-Tree Mediterranean', 'Lemon peel'),
('Bols Genever', 'Bols', 'Netherlands', 'Amsterdam', 'Genever', 42.0, 700, 'Original Dutch genever since 1820', 'Malt wine, juniper, herbs', 'Rich, malty, juniper', 'Smooth, malty finish', 'Best neat or cocktails', 'None'),
('Bobby''s Genever', 'Bobby''s', 'Netherlands', 'Schiedam', 'Genever', 38.0, 700, 'Indonesian spiced genever', 'Lemongrass, clove, juniper', 'Malty, spiced, complex', 'Warming, exotic', 'Best neat', 'None'),

-- =====================================================
-- NAVY STRENGTH GINS
-- =====================================================
('Plymouth Navy Strength', 'Plymouth', 'England', 'Plymouth', 'Navy Strength', 57.0, 700, 'Original Navy strength recipe', 'Intense juniper, citrus', 'Bold, juniper-heavy, citrus', 'Long, warming', 'Fever-Tree Indian', 'Lime wedge'),
('Sipsmith VJOP', 'Sipsmith', 'England', 'London', 'Navy Strength', 57.7, 700, 'Very Junipery Over Proof', 'Intense juniper, citrus peel', 'Powerful juniper, bold', 'Long, intense', 'Fever-Tree Indian', 'Lemon peel'),
('Four Pillars Navy Strength', 'Four Pillars', 'Australia', 'Yarra Valley', 'Navy Strength', 58.8, 700, 'Australian native botanicals', 'Intense citrus, pepper', 'Bold, spicy, citrus', 'Long, peppery', 'Fever-Tree Indian', 'Grapefruit'),
('Tarquin''s Cornish Pastis', 'Tarquin''s', 'England', 'Cornwall', 'Navy Strength', 57.0, 700, 'Cornish Navy Strength', 'Bold juniper, violet', 'Intense floral, juniper', 'Long, warming', 'Fever-Tree Indian', 'Orange peel');

-- Show summary
SELECT gin_type, COUNT(*) as count FROM gin_references GROUP BY gin_type ORDER BY count DESC;
SELECT country, COUNT(*) as count FROM gin_references GROUP BY country ORDER BY count DESC LIMIT 10;
SELECT 'Total reference gins:' as info, COUNT(*) as count FROM gin_references;
