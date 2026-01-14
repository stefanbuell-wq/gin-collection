package testutil

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// TestDB represents a test database connection
type TestDB struct {
	DB     *sql.DB
	DBName string
}

// SetupTestDB creates a temporary test database
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Connect to MySQL server (not a specific database)
	db, err := sql.Open("mysql", "root:test_password@tcp(localhost:3306)/")
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Create a unique test database
	dbName := fmt.Sprintf("test_gin_collection_%s", uuid.New().String()[:8])
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Connect to the test database
	testDB, err := sql.Open("mysql", fmt.Sprintf("root:test_password@tcp(localhost:3306)/%s?parseTime=true", dbName))
	if err != nil {
		db.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{
		DB:     testDB,
		DBName: dbName,
	}
}

// Teardown cleans up the test database
func (tdb *TestDB) Teardown(t *testing.T) {
	t.Helper()

	// Close connection
	tdb.DB.Close()

	// Connect to server to drop database
	db, err := sql.Open("mysql", "root:test_password@tcp(localhost:3306)/")
	if err != nil {
		t.Logf("Warning: Failed to connect for cleanup: %v", err)
		return
	}
	defer db.Close()

	// Drop test database
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", tdb.DBName))
	if err != nil {
		t.Logf("Warning: Failed to drop test database: %v", err)
	}
}

// RunMigrations runs database migrations for testing
func (tdb *TestDB) RunMigrations(t *testing.T) {
	t.Helper()

	// Create basic schema for testing
	schema := `
	CREATE TABLE IF NOT EXISTS tenants (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		uuid VARCHAR(36) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		subdomain VARCHAR(63) UNIQUE NOT NULL,
		tier VARCHAR(20) NOT NULL DEFAULT 'free',
		status VARCHAR(20) NOT NULL DEFAULT 'active',
		db_connection_string TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		tenant_id BIGINT NOT NULL,
		uuid VARCHAR(36) UNIQUE NOT NULL,
		email VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		role VARCHAR(20) NOT NULL DEFAULT 'member',
		api_key VARCHAR(100),
		is_active BOOLEAN DEFAULT TRUE,
		email_verified_at TIMESTAMP NULL,
		last_login_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY unique_tenant_email (tenant_id, email),
		INDEX idx_tenant_id (tenant_id),
		INDEX idx_api_key (api_key)
	);

	CREATE TABLE IF NOT EXISTS gins (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		tenant_id BIGINT NOT NULL,
		uuid VARCHAR(36) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		brand VARCHAR(255),
		country VARCHAR(100),
		gin_type VARCHAR(50),
		abv DECIMAL(4,2),
		fill_level VARCHAR(20),
		price DECIMAL(10,2),
		barcode VARCHAR(50),
		rating INT,
		is_favorite BOOLEAN DEFAULT FALSE,
		is_available BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_tenant_id (tenant_id),
		INDEX idx_tenant_name (tenant_id, name)
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		tenant_id BIGINT NOT NULL,
		user_id BIGINT,
		action VARCHAR(100) NOT NULL,
		entity_type VARCHAR(50) NOT NULL,
		entity_id BIGINT,
		changes TEXT,
		ip_address VARCHAR(45),
		user_agent VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_tenant_id (tenant_id),
		INDEX idx_user_id (user_id),
		INDEX idx_entity (entity_type, entity_id)
	);
	`

	_, err := tdb.DB.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// SeedTestData inserts test data
func (tdb *TestDB) SeedTestData(t *testing.T) (tenant1ID, tenant2ID, user1ID, user2ID int64) {
	t.Helper()

	// Insert tenant 1
	result, err := tdb.DB.Exec(`
		INSERT INTO tenants (uuid, name, subdomain, tier, status)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), "Tenant 1", "tenant1", "free", "active")
	if err != nil {
		t.Fatalf("Failed to seed tenant 1: %v", err)
	}
	tenant1ID, _ = result.LastInsertId()

	// Insert tenant 2
	result, err = tdb.DB.Exec(`
		INSERT INTO tenants (uuid, name, subdomain, tier, status)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), "Tenant 2", "tenant2", "pro", "active")
	if err != nil {
		t.Fatalf("Failed to seed tenant 2: %v", err)
	}
	tenant2ID, _ = result.LastInsertId()

	// Insert user for tenant 1
	result, err = tdb.DB.Exec(`
		INSERT INTO users (tenant_id, uuid, email, password_hash, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), "user1@tenant1.com", "hash", "owner", true)
	if err != nil {
		t.Fatalf("Failed to seed user 1: %v", err)
	}
	user1ID, _ = result.LastInsertId()

	// Insert user for tenant 2
	result, err = tdb.DB.Exec(`
		INSERT INTO users (tenant_id, uuid, email, password_hash, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenant2ID, uuid.New().String(), "user1@tenant2.com", "hash", "owner", true)
	if err != nil {
		t.Fatalf("Failed to seed user 2: %v", err)
	}
	user2ID, _ = result.LastInsertId()

	return
}
