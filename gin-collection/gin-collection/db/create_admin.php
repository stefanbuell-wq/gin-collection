<?php
// Create admin user for existing database
// Use this if the database exists but has no admin user

require_once __DIR__ . '/../api/Database.php';

try {
    $db = Database::getInstance()->getConnection();

    echo "Creating admin user...\n";

    // Check if admin already exists
    $stmt = $db->prepare("SELECT id FROM users WHERE username = ?");
    $stmt->execute(['admin']);

    if ($stmt->fetch()) {
        echo "❌ Admin user already exists!\n";
        echo "Username: admin\n";
        echo "If you forgot the password, you need to reset it manually.\n";
        exit(1);
    }

    // Create admin user
    $passwordHash = password_hash('Admin123!', PASSWORD_BCRYPT, ['cost' => 12]);

    $stmt = $db->prepare("
        INSERT INTO users (username, email, password_hash, full_name)
        VALUES (?, ?, ?, ?)
    ");

    $stmt->execute(['admin', 'admin@gin-collection.local', $passwordHash, 'Administrator']);

    $adminId = $db->lastInsertId();

    echo "✅ Admin user created successfully!\n\n";
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n";
    echo "Admin credentials:\n";
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n";
    echo "  Username: admin\n";
    echo "  Email:    admin@gin-collection.local\n";
    echo "  Password: Admin123!\n";
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n";
    echo "⚠️  IMPORTANT: Change the password after first login!\n\n";
    echo "Login at: https://atlas-bergedorf.de/GinVault/login.html\n";

} catch (Exception $e) {
    echo "❌ Error: " . $e->getMessage() . "\n";
    exit(1);
}
