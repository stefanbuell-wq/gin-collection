-- Enable strict mode
SET SQL_MODE='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- =======================
-- MULTI-TENANCY TABLES
-- =======================

CREATE TABLE tenants (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(63) UNIQUE NOT NULL,
    tier ENUM('free', 'basic', 'pro', 'enterprise') NOT NULL DEFAULT 'free',
    is_enterprise BOOLEAN DEFAULT FALSE,
    db_connection_string VARCHAR(512) COMMENT 'For Enterprise: separate DB connection',
    status ENUM('active', 'suspended', 'cancelled') NOT NULL DEFAULT 'active',
    settings JSON COMMENT 'Tenant-specific settings',
    branding JSON COMMENT 'Custom branding: logo, colors, domain',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_subdomain (subdomain),
    INDEX idx_tier (tier),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    uuid CHAR(36) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role ENUM('owner', 'admin', 'member', 'viewer') NOT NULL DEFAULT 'member',
    api_key VARCHAR(64) UNIQUE COMMENT 'Enterprise only',
    is_active BOOLEAN DEFAULT TRUE,
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_email_per_tenant (tenant_id, email),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_email (email),
    INDEX idx_api_key (api_key),
    INDEX idx_role (role),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE subscriptions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL UNIQUE,
    plan_id VARCHAR(50) NOT NULL COMMENT 'free, basic, pro, enterprise',
    status ENUM('active', 'past_due', 'cancelled', 'trialing') NOT NULL DEFAULT 'active',
    billing_cycle ENUM('monthly', 'yearly') NOT NULL DEFAULT 'monthly',
    current_period_start DATE NOT NULL,
    current_period_end DATE NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    paypal_customer_id VARCHAR(255),
    paypal_subscription_id VARCHAR(255),
    trial_ends_at TIMESTAMP NULL,
    cancelled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_plan_id (plan_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE usage_metrics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    metric_name VARCHAR(50) NOT NULL COMMENT 'gin_count, api_calls, storage_mb',
    current_value INT UNSIGNED NOT NULL DEFAULT 0,
    limit_value INT UNSIGNED COMMENT 'NULL = unlimited',
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_metric_per_tenant_period (tenant_id, metric_name, period_start),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_metric_name (metric_name),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =======================
-- CORE GIN TABLES
-- =======================

CREATE TABLE gins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    uuid CHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    country VARCHAR(100),
    region VARCHAR(100),
    gin_type VARCHAR(50),
    abv DECIMAL(4,2) COMMENT 'Alcohol by volume percentage',
    bottle_size INT UNSIGNED DEFAULT 700 COMMENT 'Size in ml',
    fill_level TINYINT UNSIGNED DEFAULT 100 COMMENT 'Percentage 0-100',
    price DECIMAL(10,2) COMMENT 'Purchase price',
    current_market_value DECIMAL(10,2),
    purchase_date DATE,
    purchase_location VARCHAR(255),
    barcode VARCHAR(50),
    rating TINYINT CHECK(rating >= 1 AND rating <= 5),
    nose_notes TEXT COMMENT 'Tasting: aroma',
    palate_notes TEXT COMMENT 'Tasting: taste',
    finish_notes TEXT COMMENT 'Tasting: finish',
    general_notes TEXT,
    description TEXT,
    photo_url VARCHAR(512),
    is_finished BOOLEAN DEFAULT FALSE,
    recommended_tonic VARCHAR(255),
    recommended_garnish VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_tenant_name (tenant_id, name),
    INDEX idx_tenant_brand (tenant_id, brand),
    INDEX idx_tenant_country (tenant_id, country),
    INDEX idx_tenant_barcode (tenant_id, barcode),
    INDEX idx_tenant_rating (tenant_id, rating),
    INDEX idx_tenant_finished (tenant_id, is_finished),
    UNIQUE KEY unique_barcode_per_tenant (tenant_id, barcode),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FULLTEXT INDEX ft_search (name, brand, nose_notes, palate_notes, finish_notes, general_notes)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =======================
-- SHARED REFERENCE DATA
-- =======================

CREATE TABLE botanicals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50) COMMENT 'Citrus, Spices, Flowers, Herbs, Roots, Vegetables',
    description TEXT,
    INDEX idx_category (category),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE cocktails (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    instructions TEXT,
    glass_type VARCHAR(50),
    ice_type VARCHAR(50),
    difficulty ENUM('easy', 'medium', 'hard') DEFAULT 'easy',
    prep_time INT UNSIGNED COMMENT 'Preparation time in minutes',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_difficulty (difficulty),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE cocktail_ingredients (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cocktail_id BIGINT UNSIGNED NOT NULL,
    ingredient VARCHAR(255) NOT NULL,
    amount VARCHAR(50) COMMENT 'e.g., 50ml, 2 dashes',
    unit VARCHAR(20),
    is_gin BOOLEAN DEFAULT FALSE,
    INDEX idx_cocktail_id (cocktail_id),
    FOREIGN KEY (cocktail_id) REFERENCES cocktails(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =======================
-- TENANT-SCOPED RELATIONSHIPS
-- =======================

CREATE TABLE gin_botanicals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    gin_id BIGINT UNSIGNED NOT NULL,
    botanical_id BIGINT UNSIGNED NOT NULL,
    prominence ENUM('dominant', 'notable', 'subtle') DEFAULT 'notable',
    UNIQUE KEY unique_gin_botanical_per_tenant (tenant_id, gin_id, botanical_id),
    INDEX idx_tenant_gin (tenant_id, gin_id),
    INDEX idx_tenant_botanical (tenant_id, botanical_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (botanical_id) REFERENCES botanicals(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE gin_cocktails (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    gin_id BIGINT UNSIGNED NOT NULL,
    cocktail_id BIGINT UNSIGNED NOT NULL,
    UNIQUE KEY unique_gin_cocktail_per_tenant (tenant_id, gin_id, cocktail_id),
    INDEX idx_tenant_gin (tenant_id, gin_id),
    INDEX idx_cocktail_id (cocktail_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (cocktail_id) REFERENCES cocktails(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE gin_photos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    gin_id BIGINT UNSIGNED NOT NULL,
    photo_url VARCHAR(512) NOT NULL,
    photo_type ENUM('bottle', 'label', 'moment', 'tasting') DEFAULT 'bottle',
    caption VARCHAR(500),
    is_primary BOOLEAN DEFAULT FALSE,
    storage_key VARCHAR(255) COMMENT 'S3 object key',
    file_size_kb INT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_gin (tenant_id, gin_id),
    INDEX idx_tenant_primary (tenant_id, is_primary),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tasting_sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    gin_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED,
    date DATE NOT NULL,
    notes TEXT,
    rating TINYINT CHECK(rating >= 1 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_gin (tenant_id, gin_id),
    INDEX idx_tenant_date (tenant_id, date),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =======================
-- AUDIT & COMPLIANCE
-- =======================

CREATE TABLE audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED,
    action VARCHAR(100) NOT NULL COMMENT 'create, update, delete, login, etc.',
    entity_type VARCHAR(50) NOT NULL COMMENT 'gin, user, subscription, etc.',
    entity_id BIGINT UNSIGNED,
    changes JSON COMMENT 'Before/after values',
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =======================
-- SEED DATA: Botanicals
-- =======================

INSERT INTO botanicals (name, category, description) VALUES
('Wacholder', 'Basis', 'Hauptzutat in allen Gins, verleiht den charakteristischen Gin-Geschmack'),
('Koriandersamen', 'Gewürze', 'Häufige Zutat, verleiht würzige und zitrusartige Noten'),
('Angelikawurzel', 'Wurzeln', 'Erdige, holzige Aromen mit leichter Süße'),
('Zitronenschale', 'Zitrus', 'Frische, lebendige Zitrusnoten'),
('Orangenschale', 'Zitrus', 'Süßere, weichere Zitrusnoten'),
('Grapefruitschale', 'Zitrus', 'Bittere, frische Zitrusnoten'),
('Zimt', 'Gewürze', 'Warme, süße Gewürznoten'),
('Kardamom', 'Gewürze', 'Komplexe, süß-würzige Aromen'),
('Kubebenpfeffer', 'Gewürze', 'Pfeffrig, würzig mit leichter Bitterkeit'),
('Süßholzwurzel', 'Wurzeln', 'Süße, leicht anisartige Noten'),
('Iriswurzel', 'Wurzeln', 'Blumig, holzig, fixiert andere Aromen'),
('Lavendel', 'Blüten', 'Blumige, leicht süße Noten'),
('Rosenblätter', 'Blüten', 'Delicate, parfümierte florale Noten'),
('Kamille', 'Blüten', 'Sanfte, apfelartige florale Noten'),
('Gurke', 'Gemüse', 'Frische, kühle, grüne Noten'),
('Pfeffer', 'Gewürze', 'Scharfe, würzige Noten'),
('Ingwer', 'Wurzeln', 'Scharfe, wärmende, würzige Noten'),
('Thymian', 'Kräuter', 'Erdige, leicht minzige Kräuternoten'),
('Salbei', 'Kräuter', 'Erdige, würzige Kräuternoten'),
('Minze', 'Kräuter', 'Frische, kühle, mentholhaltige Noten');

-- =======================
-- SEED DATA: Cocktails
-- =======================

INSERT INTO cocktails (name, description, instructions, glass_type, ice_type, difficulty, prep_time) VALUES
('Gin & Tonic', 'Der Klassiker - einfach und erfrischend', 'Gin in ein Highball-Glas mit Eis geben. Mit Tonic Water auffüllen. Mit Zitrone oder Gurke garnieren.', 'Highball', 'Würfel', 'easy', 2),
('Negroni', 'Bitter-süßer italienischer Klassiker', 'Gin, Campari und süßen Wermut zu gleichen Teilen in ein Rührglas mit Eis geben. Rühren und in ein Rocks-Glas abseihen. Mit Orangenzeste garnieren.', 'Rocks', 'Würfel', 'easy', 3),
('Martini', 'Der ultimative Gin-Cocktail', 'Gin und trockenen Wermut in ein Rührglas mit Eis geben. Gut rühren und in eine gekühlte Martinischale abseihen. Mit Olive oder Zitronenzeste garnieren.', 'Martini', 'keine', 'medium', 5),
('Gin Fizz', 'Spritzig und erfrischend', 'Gin, Zitronensaft und Zuckersirup in einen Shaker mit Eis geben. Kräftig schütteln und in ein Highball-Glas abseihen. Mit Sodawasser auffüllen.', 'Highball', 'Würfel', 'easy', 3),
('Tom Collins', 'Klassischer langer Drink', 'Gin, Zitronensaft und Zuckersirup in einen Shaker mit Eis geben. Schütteln und in ein Collins-Glas mit Eis abseihen. Mit Sodawasser auffüllen und mit Zitrone garnieren.', 'Collins', 'Würfel', 'easy', 4);

-- Cocktail Ingredients
INSERT INTO cocktail_ingredients (cocktail_id, ingredient, amount, unit, is_gin) VALUES
-- Gin & Tonic
(1, 'Gin', '50', 'ml', TRUE),
(1, 'Tonic Water', '150', 'ml', FALSE),
(1, 'Zitrone oder Gurke', '1', 'Scheibe', FALSE),

-- Negroni
(2, 'Gin', '30', 'ml', TRUE),
(2, 'Campari', '30', 'ml', FALSE),
(2, 'Süßer Wermut', '30', 'ml', FALSE),
(2, 'Orangenzeste', '1', 'Stück', FALSE),

-- Martini
(3, 'Gin', '60', 'ml', TRUE),
(3, 'Trockener Wermut', '10', 'ml', FALSE),
(3, 'Olive oder Zitronenzeste', '1', 'Stück', FALSE),

-- Gin Fizz
(4, 'Gin', '50', 'ml', TRUE),
(4, 'Zitronensaft', '25', 'ml', FALSE),
(4, 'Zuckersirup', '15', 'ml', FALSE),
(4, 'Sodawasser', '100', 'ml', FALSE),

-- Tom Collins
(5, 'Gin', '50', 'ml', TRUE),
(5, 'Zitronensaft', '25', 'ml', FALSE),
(5, 'Zuckersirup', '15', 'ml', FALSE),
(5, 'Sodawasser', '100', 'ml', FALSE),
(5, 'Zitronenscheibe', '1', 'Stück', FALSE);
