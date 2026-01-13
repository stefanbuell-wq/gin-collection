<?php
// Database migration script for multi-tenancy
// This script migrates existing database to support multiple users

require_once __DIR__ . '/../api/Database.php';

class Migration {
    private $db;
    
    public function __construct() {
        $this->db = Database::getInstance()->getConnection();
    }
    
    public function run() {
        echo "Starting database migration...\n\n";
        
        try {
            // Start transaction
            $this->db->beginTransaction();
            
            // Step 1: Check if migration is needed
            if ($this->isMigrated()) {
                echo "✓ Database is already migrated.\n";
                return;
            }
            
            // Step 2: Create users table
            echo "Creating users table...\n";
            $this->createUsersTable();
            
            // Step 3: Create default admin user
            echo "Creating default admin user...\n";
            $adminId = $this->createAdminUser();
            
            // Step 4: Add user_id column to gins table
            echo "Adding user_id column to gins table...\n";
            $this->addUserIdToGins($adminId);
            
            // Step 5: Add user_id column to tasting_sessions table
            echo "Adding user_id column to tasting_sessions table...\n";
            $this->addUserIdToTastingSessions($adminId);
            
            // Step 6: Create indexes
            echo "Creating indexes...\n";
            $this->createIndexes();
            
            // Commit transaction
            $this->db->commit();
            
            echo "\n✓ Migration completed successfully!\n\n";
            echo "Default admin credentials:\n";
            echo "  Username: admin\n";
            echo "  Email: admin@gin-collection.local\n";
            echo "  Password: Admin123!\n\n";
            echo "⚠️  IMPORTANT: Please change the admin password after first login!\n";
            
        } catch (Exception $e) {
            $this->db->rollBack();
            echo "\n✗ Migration failed: " . $e->getMessage() . "\n";
            throw $e;
        }
    }
    
    private function isMigrated() {
        // Check if users table exists
        $result = $this->db->query("SELECT name FROM sqlite_master WHERE type='table' AND name='users'");
        return $result->fetch() !== false;
    }
    
    private function createUsersTable() {
        $this->db->exec("
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT UNIQUE NOT NULL,
                email TEXT UNIQUE NOT NULL,
                password_hash TEXT NOT NULL,
                full_name TEXT,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )
        ");
        
        // Create trigger for users table
        $this->db->exec("
            CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
            AFTER UPDATE ON users
            BEGIN
                UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
            END
        ");
    }
    
    private function createAdminUser() {
        // Default admin password: Admin123!
        $passwordHash = password_hash('Admin123!', PASSWORD_BCRYPT, ['cost' => 12]);
        
        $stmt = $this->db->prepare("
            INSERT INTO users (username, email, password_hash, full_name)
            VALUES (?, ?, ?, ?)
        ");
        
        $stmt->execute(['admin', 'admin@gin-collection.local', $passwordHash, 'Administrator']);
        
        return $this->db->lastInsertId();
    }
    
    private function addUserIdToGins($adminId) {
        // Check if column already exists
        $columns = $this->db->query("PRAGMA table_info(gins)")->fetchAll();
        $hasUserIdColumn = false;
        
        foreach ($columns as $column) {
            if ($column['name'] === 'user_id') {
                $hasUserIdColumn = true;
                break;
            }
        }
        
        if (!$hasUserIdColumn) {
            // SQLite doesn't support ALTER TABLE ADD COLUMN with FOREIGN KEY
            // So we need to recreate the table
            
            // Create temporary table with new schema
            $this->db->exec("
                CREATE TABLE gins_new (
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
                )
            ");
            
            // Copy data from old table to new table
            $this->db->exec("
                INSERT INTO gins_new (id, user_id, name, brand, country, region, abv, bottle_size, 
                                     price, purchase_date, barcode, rating, tasting_notes, 
                                     description, photo_url, is_finished, created_at, updated_at)
                SELECT id, $adminId, name, brand, country, region, abv, bottle_size, 
                       price, purchase_date, barcode, rating, tasting_notes, 
                       description, photo_url, is_finished, created_at, updated_at
                FROM gins
            ");
            
            // Drop old table
            $this->db->exec("DROP TABLE gins");
            
            // Rename new table
            $this->db->exec("ALTER TABLE gins_new RENAME TO gins");
        }
    }
    
    private function addUserIdToTastingSessions($adminId) {
        // Check if tasting_sessions table exists
        $result = $this->db->query("SELECT name FROM sqlite_master WHERE type='table' AND name='tasting_sessions'");
        
        if ($result->fetch() !== false) {
            // Check if column already exists
            $columns = $this->db->query("PRAGMA table_info(tasting_sessions)")->fetchAll();
            $hasUserIdColumn = false;
            
            foreach ($columns as $column) {
                if ($column['name'] === 'user_id') {
                    $hasUserIdColumn = true;
                    break;
                }
            }
            
            if (!$hasUserIdColumn) {
                // Recreate table with user_id
                $this->db->exec("
                    CREATE TABLE tasting_sessions_new (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        user_id INTEGER NOT NULL,
                        gin_id INTEGER,
                        date DATE NOT NULL,
                        notes TEXT,
                        rating INTEGER CHECK(rating >= 1 AND rating <= 5),
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                        FOREIGN KEY (gin_id) REFERENCES gins(id) ON DELETE CASCADE
                    )
                ");
                
                // Copy data if any exists
                $this->db->exec("
                    INSERT INTO tasting_sessions_new (id, user_id, gin_id, date, notes, rating)
                    SELECT id, $adminId, gin_id, date, notes, rating
                    FROM tasting_sessions
                ");
                
                $this->db->exec("DROP TABLE tasting_sessions");
                $this->db->exec("ALTER TABLE tasting_sessions_new RENAME TO tasting_sessions");
            }
        }
    }
    
    private function createIndexes() {
        $this->db->exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)");
        $this->db->exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)");
        $this->db->exec("CREATE INDEX IF NOT EXISTS idx_gins_user_id ON gins(user_id)");
        $this->db->exec("CREATE INDEX IF NOT EXISTS idx_tasting_sessions_user_id ON tasting_sessions(user_id)");
    }
}

// Run migration if called directly
if (php_sapi_name() === 'cli') {
    $migration = new Migration();
    $migration->run();
} else {
    echo "This script must be run from the command line.\n";
}
