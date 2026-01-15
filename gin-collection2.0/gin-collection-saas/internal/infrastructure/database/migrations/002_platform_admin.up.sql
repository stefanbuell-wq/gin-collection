-- Platform Admin table for super-admin access
CREATE TABLE IF NOT EXISTS platform_admins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_email (email),
    INDEX idx_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Default Platform Admin (Password: admin123)
-- IMPORTANT: Change this password after first login!
INSERT INTO platform_admins (email, password_hash, name, is_active) VALUES
('admin@gin-collection.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.Y6VELdBlKmzaqLrPVe', 'Platform Admin', true);
