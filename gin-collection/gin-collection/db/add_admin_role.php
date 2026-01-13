<?php
// Migration script to add admin role to existing database
// This adds the is_admin column to users table and marks the admin user

require_once __DIR__ . '/../api/Database.php';

try {
    $db = Database::getInstance()->getConnection();

    echo "🔧 Adding admin role to database...\n\n";

    // Start transaction
    $db->beginTransaction();

    // Check if is_admin column already exists
    $columns = $db->query("PRAGMA table_info(users)")->fetchAll();
    $hasIsAdminColumn = false;

    foreach ($columns as $column) {
        if ($column['name'] === 'is_admin') {
            $hasIsAdminColumn = true;
            break;
        }
    }

    if ($hasIsAdminColumn) {
        echo "✓ is_admin column already exists\n";
    } else {
        echo "Adding is_admin column to users table...\n";
        $db->exec("ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0");
        echo "✓ is_admin column added\n";
    }

    // Create index for is_admin
    echo "Creating index for is_admin...\n";
    $db->exec("CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin)");
    echo "✓ Index created\n";

    // Mark admin user as admin
    echo "\nMarking admin user as admin...\n";
    $stmt = $db->prepare("UPDATE users SET is_admin = 1 WHERE username = ?");
    $stmt->execute(['admin']);

    if ($stmt->rowCount() > 0) {
        echo "✓ Admin user marked as admin\n";
    } else {
        echo "⚠️  Warning: No user with username 'admin' found\n";
        echo "   You can create an admin user by registering and then manually updating the database.\n";
    }

    // Commit transaction
    $db->commit();

    echo "\n✅ Migration completed successfully!\n\n";
    echo "Admin user now has admin privileges.\n";
    echo "Login at: https://atlas-bergedorf.de/GinVault/login.html\n";

} catch (Exception $e) {
    if (isset($db) && $db->inTransaction()) {
        $db->rollBack();
    }
    echo "\n❌ Migration failed: " . $e->getMessage() . "\n";
    echo "\nStack trace:\n";
    echo $e->getTraceAsString() . "\n";
    exit(1);
}
