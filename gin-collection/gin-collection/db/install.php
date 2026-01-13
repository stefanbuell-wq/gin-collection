<?php
// Fresh installation script for Gin Collection
// This script creates a new database from scratch using schema.sql
// Use this for first-time installations (NOT for migrations)

require_once __DIR__ . '/../api/Database.php';

class FreshInstall
{
    private $db;
    private $dbPath;

    public function __construct()
    {
        $this->dbPath = __DIR__ . '/gin_collection.db';
    }

    public function run()
    {
        echo "🍸 Gin Collection - Fresh Installation\n";
        echo "========================================\n\n";

        try {
            // Step 1: Check if database already exists
            if (file_exists($this->dbPath)) {
                echo "⚠️  WARNING: Database file already exists!\n";
                echo "Path: {$this->dbPath}\n\n";
                echo "Do you want to DELETE the existing database and create a fresh one?\n";
                echo "This will PERMANENTLY DELETE all existing data!\n\n";
                echo "Type 'YES' to continue or anything else to abort: ";

                $handle = fopen("php://stdin", "r");
                $line = trim(fgets($handle));
                fclose($handle);

                if ($line !== 'YES') {
                    echo "\n❌ Installation aborted.\n";
                    echo "If you want to migrate an existing database, use: php db/migrate.php\n";
                    return;
                }

                // Delete existing database
                unlink($this->dbPath);
                echo "\n✓ Existing database deleted.\n\n";
            }

            // Step 2: Initialize database connection
            echo "Creating new database...\n";
            $this->db = Database::getInstance()->getConnection();

            // Step 3: Load and execute schema.sql
            echo "Loading schema from schema.sql...\n";
            $schema = file_get_contents(__DIR__ . '/schema.sql');

            if ($schema === false) {
                throw new Exception("Could not read schema.sql file");
            }

            // Execute schema
            echo "Creating tables and indexes...\n";
            $this->db->exec($schema);

            // Step 4: Create default admin user
            echo "Creating default admin user...\n";
            $this->createAdminUser();

            echo "\n✅ Installation completed successfully!\n\n";
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n";
            echo "Default admin credentials:\n";
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n";
            echo "  Username: admin\n";
            echo "  Email:    admin@gin-collection.local\n";
            echo "  Password: Admin123!\n";
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n";
            echo "⚠️  IMPORTANT: Change the admin password after first login!\n\n";
            echo "Next steps:\n";
            echo "1. Open: https://atlas-bergedorf.de/GinVault/login.html\n";
            echo "2. Login with the credentials above\n";
            echo "3. Change your password immediately\n";
            echo "4. Start adding your gin collection!\n\n";

        } catch (Exception $e) {
            echo "\n❌ Installation failed: " . $e->getMessage() . "\n";
            echo "\nStack trace:\n";
            echo $e->getTraceAsString() . "\n";
            throw $e;
        }
    }

    private function createAdminUser()
    {
        // Default admin password: Admin123!
        $passwordHash = password_hash('Admin123!', PASSWORD_BCRYPT, ['cost' => 12]);

        $stmt = $this->db->prepare("
            INSERT INTO users (username, email, password_hash, full_name)
            VALUES (?, ?, ?, ?)
        ");

        $stmt->execute(['admin', 'admin@gin-collection.local', $passwordHash, 'Administrator']);

        echo "✓ Admin user created (ID: " . $this->db->lastInsertId() . ")\n";
    }

    public function checkPrerequisites()
    {
        echo "Checking prerequisites...\n\n";

        $allGood = true;

        // Check PHP version
        $phpVersion = phpversion();
        echo "PHP Version: $phpVersion ";
        if (version_compare($phpVersion, '7.4.0', '>=')) {
            echo "✓\n";
        } else {
            echo "✗ (7.4+ required)\n";
            $allGood = false;
        }

        // Check SQLite extension
        echo "SQLite3 Extension: ";
        if (extension_loaded('sqlite3')) {
            echo "✓\n";
        } else {
            echo "✗ (required)\n";
            $allGood = false;
        }

        // Check PDO SQLite
        echo "PDO SQLite: ";
        if (extension_loaded('pdo_sqlite')) {
            echo "✓\n";
        } else {
            echo "✗ (required)\n";
            $allGood = false;
        }

        // Check directory permissions
        $dbDir = __DIR__;
        echo "Database directory writable: ";
        if (is_writable($dbDir)) {
            echo "✓\n";
        } else {
            echo "✗ (chmod 755 $dbDir)\n";
            $allGood = false;
        }

        // Check uploads directory
        $uploadsDir = dirname(__DIR__) . '/uploads';
        echo "Uploads directory writable: ";
        if (is_writable($uploadsDir)) {
            echo "✓\n";
        } else {
            echo "✗ (chmod 755 $uploadsDir)\n";
            $allGood = false;
        }

        echo "\n";

        if (!$allGood) {
            echo "❌ Prerequisites not met. Please fix the issues above.\n";
            exit(1);
        }

        echo "✅ All prerequisites met!\n\n";
    }
}

// Run installation if called directly
if (php_sapi_name() === 'cli') {
    $install = new FreshInstall();
    $install->checkPrerequisites();
    $install->run();
} else {
    echo "This script must be run from the command line.\n";
    echo "Usage: php db/install.php\n";
}
