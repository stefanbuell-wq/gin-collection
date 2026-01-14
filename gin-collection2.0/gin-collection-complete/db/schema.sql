-- Gin Collection Database Schema

CREATE TABLE IF NOT EXISTS gins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brand TEXT,
    country TEXT,
    region TEXT,
    gin_type TEXT, -- London Dry, Old Tom, New Western, Plymouth, etc.
    abv REAL,
    bottle_size INTEGER DEFAULT 700,
    fill_level INTEGER DEFAULT 100, -- Füllstand in Prozent (100, 75, 50, 25, 0)
    price REAL,
    current_market_value REAL, -- Aktueller Marktwert
    purchase_date DATE,
    purchase_location TEXT, -- Händler/Shop
    barcode TEXT UNIQUE,
    rating INTEGER CHECK(rating >= 1 AND rating <= 5),
    
    -- Strukturierte Tasting-Notizen
    nose_notes TEXT, -- Aroma/Nase
    palate_notes TEXT, -- Geschmack/Gaumen  
    finish_notes TEXT, -- Abgang
    general_notes TEXT, -- Allgemeine Notizen
    
    description TEXT,
    photo_url TEXT,
    is_finished INTEGER DEFAULT 0,
    
    -- Serviervorschläge
    recommended_tonic TEXT, -- Empfohlenes Tonic
    recommended_garnish TEXT, -- Empfohlene Garnitur
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS botanicals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    category TEXT, -- Zitrus, Gewürz, Blüten, Kräuter, Wurzeln
    description TEXT
);

CREATE TABLE IF NOT EXISTS gin_botanicals (
    gin_id INTEGER,
    botanical_id INTEGER,
    prominence TEXT, -- dominant, notable, subtle
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (botanical_id) REFERENCES botanicals(id) ON DELETE CASCADE,
    PRIMARY KEY (gin_id, botanical_id)
);

-- Vorgefertigte Botanicals
INSERT OR IGNORE INTO botanicals (name, category, description) VALUES
('Wacholder', 'Kräuter', 'Hauptbestandteil von Gin'),
('Koriander', 'Gewürz', 'Würzige, zitrusartige Note'),
('Angelikawurzel', 'Wurzeln', 'Erdige, holzige Note'),
('Zitronenschale', 'Zitrus', 'Frische Zitrusnote'),
('Orangenschale', 'Zitrus', 'Süße Zitrusnote'),
('Grapefruitschale', 'Zitrus', 'Herbe Zitrusnote'),
('Zimt', 'Gewürz', 'Warme, süße Würze'),
('Kardamom', 'Gewürz', 'Würzig-süße Note'),
('Kubebenpfeffer', 'Gewürz', 'Pfeffrige, leicht mentholige Note'),
('Süßholzwurzel', 'Wurzeln', 'Süße, lakritzig-holzige Note'),
('Iriswurzel', 'Wurzeln', 'Blumige, pudrige Note'),
('Lavendel', 'Blüten', 'Blumige, krautige Note'),
('Rosenblüten', 'Blüten', 'Zarte, florale Note'),
('Kamille', 'Blüten', 'Milde, apfelartige Note'),
('Gurke', 'Gemüse', 'Frische, grüne Note'),
('Pfeffer', 'Gewürz', 'Scharfe, würzige Note'),
('Ingwer', 'Wurzeln', 'Scharfe, zitrusartige Note'),
('Thymian', 'Kräuter', 'Kräuterige Note'),
('Salbei', 'Kräuter', 'Würzige Kräuternote'),
('Minze', 'Kräuter', 'Frische, kühle Note');

CREATE TABLE IF NOT EXISTS tasting_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gin_id INTEGER,
    date DATE NOT NULL,
    notes TEXT,
    rating INTEGER CHECK(rating >= 1 AND rating <= 5),
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE
);

-- Cocktail-Rezepte
CREATE TABLE IF NOT EXISTS cocktails (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    instructions TEXT,
    glass_type TEXT,
    ice_type TEXT,
    difficulty TEXT, -- easy, medium, hard
    prep_time INTEGER, -- in Minuten
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cocktail_ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cocktail_id INTEGER,
    ingredient TEXT NOT NULL,
    amount TEXT, -- z.B. "50ml", "2 dash", "1 TL"
    unit TEXT,
    is_gin INTEGER DEFAULT 0, -- 1 wenn Gin, 0 wenn andere Zutat
    FOREIGN KEY (cocktail_id) REFERENCES cocktails(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS gin_cocktails (
    gin_id INTEGER,
    cocktail_id INTEGER,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (cocktail_id) REFERENCES cocktails(id) ON DELETE CASCADE,
    PRIMARY KEY (gin_id, cocktail_id)
);

-- Foto-Galerie für mehrere Fotos pro Gin
CREATE TABLE IF NOT EXISTS gin_photos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gin_id INTEGER,
    photo_url TEXT NOT NULL,
    photo_type TEXT, -- bottle, label, moment, tasting
    caption TEXT,
    is_primary INTEGER DEFAULT 0, -- Hauptfoto
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE
);

-- Vorgefertigte Cocktails
INSERT OR IGNORE INTO cocktails (name, description, instructions, glass_type, ice_type, difficulty, prep_time) VALUES
('Gin & Tonic', 'Der Klassiker - simpel und erfrischend', '1. Glas mit Eiswürfeln füllen\n2. 50ml Gin zugeben\n3. Mit 150ml Tonic auffüllen\n4. Mit Zitrone oder Gurke garnieren', 'Highball', 'Eiswürfel', 'easy', 2),
('Negroni', 'Italienischer Klassiker - bitter und aromatisch', '1. Alle Zutaten zu gleichen Teilen in ein Glas mit Eis geben\n2. Umrühren\n3. Mit Orangenzeste garnieren', 'Tumbler', 'Eiswürfel', 'easy', 3),
('Martini', 'Der ultimative Gin-Cocktail', '1. Gin und Wermut im Verhältnis 6:1 mit Eis rühren\n2. In gekühltes Martiniglas abseihen\n3. Mit Olive oder Zitronenzeste garnieren', 'Martini', 'gerührt', 'medium', 5),
('Gin Fizz', 'Spritzig und erfrischend', '1. Gin, Zitronensaft und Zuckersirup mit Eis shaken\n2. In Glas abseihen\n3. Mit Sodawasser auffüllen', 'Highball', 'Eiswürfel', 'easy', 3),
('Tom Collins', 'Perfekt für heiße Tage', '1. Gin, Zitronensaft und Zuckersirup mit Eis shaken\n2. In Highball-Glas mit Eis abseihen\n3. Mit Soda auffüllen und mit Zitrone garnieren', 'Highball', 'Eiswürfel', 'easy', 4);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_gins_name ON gins(name);
CREATE INDEX IF NOT EXISTS idx_gins_brand ON gins(brand);
CREATE INDEX IF NOT EXISTS idx_gins_country ON gins(country);
CREATE INDEX IF NOT EXISTS idx_gins_barcode ON gins(barcode);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_gins_timestamp 
AFTER UPDATE ON gins
BEGIN
    UPDATE gins SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
