-- Gin Collection Database Schema

-- Users table for multi-tenancy
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    is_admin INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    brand TEXT,
    country TEXT,
    region TEXT,
    abv REAL,
    bottle_size INTEGER DEFAULT 700,
    price REAL,
    purchase_date DATE,
    barcode TEXT,
    rating INTEGER CHECK(rating >= 1 AND rating <= 5),
    tasting_notes TEXT,
    description TEXT,
    photo_url TEXT,
    is_finished INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS botanicals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS gin_botanicals (
    gin_id INTEGER,
    botanical_id INTEGER,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE,
    FOREIGN KEY (botanical_id) REFERENCES botanicals(id) ON DELETE CASCADE,
    PRIMARY KEY (gin_id, botanical_id)
);

CREATE TABLE IF NOT EXISTS tasting_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    gin_id INTEGER,
    date DATE NOT NULL,
    notes TEXT,
    rating INTEGER CHECK(rating >= 1 AND rating <= 5),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);
CREATE INDEX IF NOT EXISTS idx_gins_user_id ON gins(user_id);
CREATE INDEX IF NOT EXISTS idx_gins_name ON gins(name);
CREATE INDEX IF NOT EXISTS idx_gins_brand ON gins(brand);
CREATE INDEX IF NOT EXISTS idx_gins_country ON gins(country);
CREATE INDEX IF NOT EXISTS idx_gins_barcode ON gins(barcode);
CREATE INDEX IF NOT EXISTS idx_tasting_sessions_user_id ON tasting_sessions(user_id);

-- Triggers to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_gins_timestamp 
AFTER UPDATE ON gins
BEGIN
    UPDATE gins SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
