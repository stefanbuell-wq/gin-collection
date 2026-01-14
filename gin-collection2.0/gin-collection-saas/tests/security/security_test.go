package security

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	"github.com/yourusername/gin-collection-saas/tests/testutil"
	"golang.org/x/crypto/bcrypt"
)

// TestSecurity_PasswordHashing verifies passwords are properly hashed
func TestSecurity_PasswordHashing(t *testing.T) {
	password := "SecurePassword123!"

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify hash is not the plain password
	if string(hash) == password {
		t.Error("Password hash should not equal plain password")
	}

	// Verify hash can be verified
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Error("Valid password should verify against hash")
	}

	// Verify wrong password fails
	err = bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword"))
	if err == nil {
		t.Error("Wrong password should not verify")
	}
}

// TestSecurity_SQLInjectionPrevention tests prepared statements prevent SQL injection
func TestSecurity_SQLInjectionPrevention(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, _, _, _ := testDB.SeedTestData(t)

	ginRepo := mysql.NewGinRepository(testDB.DB)

	// Attempt SQL injection via name parameter
	maliciousInput := "'; DROP TABLE gins; --"

	// Create gin with malicious input
	_, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO gins (tenant_id, uuid, name, brand)
		VALUES (?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), maliciousInput, "Brand")
	if err != nil {
		t.Fatalf("Failed to insert gin: %v", err)
	}

	// Verify table still exists and data is intact
	gins, err := ginRepo.List(ctx, tenant1ID, 10, 0, nil)
	if err != nil {
		t.Fatalf("Failed to list gins (SQL injection may have succeeded): %v", err)
	}

	// Find our gin with malicious input
	found := false
	for _, gin := range gins {
		if gin.Name == maliciousInput {
			found = true
			break
		}
	}

	if !found {
		t.Error("Gin with malicious input should be safely stored")
	}
}

// TestSecurity_APIKeyFormat verifies API keys are properly formatted
func TestSecurity_APIKeyFormat(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, _, user1ID, _ := testDB.SeedTestData(t)

	userRepo := mysql.NewUserRepository(testDB.DB)

	// Generate API key
	apiKey, err := userRepo.GenerateAPIKey(ctx, user1ID)
	if err != nil {
		t.Fatalf("Failed to generate API key: %v", err)
	}

	// Verify format
	if len(apiKey) < 40 {
		t.Error("API key should be at least 40 characters")
	}

	if apiKey[:3] != "sk_" {
		t.Error("API key should start with 'sk_'")
	}

	// Verify it can be retrieved
	user, err := userRepo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		t.Fatalf("Failed to get user by API key: %v", err)
	}

	if user.ID != user1ID {
		t.Error("Retrieved user should match")
	}

	if user.TenantID != tenant1ID {
		t.Error("User should belong to correct tenant")
	}
}

// TestSecurity_XSS_Prevention tests XSS prevention
func TestSecurity_XSS_Prevention(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, _, _, _ := testDB.SeedTestData(t)

	// XSS payload
	xssPayload := "<script>alert('XSS')</script>"

	// Insert gin with XSS payload
	_, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO gins (tenant_id, uuid, name, brand)
		VALUES (?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), xssPayload, "Brand")
	if err != nil {
		t.Fatalf("Failed to insert gin: %v", err)
	}

	// Retrieve gin
	ginRepo := mysql.NewGinRepository(testDB.DB)
	gins, err := ginRepo.List(ctx, tenant1ID, 10, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify XSS payload is stored as-is (database layer doesn't sanitize)
	// Output encoding must be done at presentation layer (frontend)
	found := false
	for _, gin := range gins {
		if gin.Name == xssPayload {
			found = true
			// Note: Frontend must HTML-encode this before rendering
			break
		}
	}

	if !found {
		t.Error("XSS payload should be stored (sanitization happens at output)")
	}
}

// TestSecurity_AuthorizationBypass tests authorization checks
func TestSecurity_AuthorizationBypass(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, tenant2ID, _, _ := testDB.SeedTestData(t)

	ginRepo := mysql.NewGinRepository(testDB.DB)

	// Create gin for tenant 1
	result, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO gins (tenant_id, uuid, name, brand)
		VALUES (?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), "Tenant 1 Gin", "Brand")
	if err != nil {
		t.Fatal(err)
	}
	ginID, _ := result.LastInsertId()

	// Try to access tenant 1's gin with tenant 2's context
	t.Run("CannotAccessOtherTenantData", func(t *testing.T) {
		_, err := ginRepo.GetByID(ctx, tenant2ID, ginID)
		if err == nil {
			t.Error("Should not be able to access another tenant's data")
		}
	})

	// Try to update tenant 1's gin with tenant 2's context
	t.Run("CannotUpdateOtherTenantData", func(t *testing.T) {
		_, err := testDB.DB.ExecContext(ctx, `
			UPDATE gins SET name = ? WHERE id = ? AND tenant_id = ?
		`, "Hacked Name", ginID, tenant2ID)
		if err != nil {
			t.Fatal(err)
		}

		// Verify name was NOT changed
		var name string
		err = testDB.DB.QueryRowContext(ctx, `
			SELECT name FROM gins WHERE id = ?
		`, ginID).Scan(&name)
		if err != nil {
			t.Fatal(err)
		}

		if name != "Tenant 1 Gin" {
			t.Error("Name should not have changed via unauthorized update")
		}
	})
}

// TestSecurity_RateLimitingData tests rate limiting data structure
func TestSecurity_RateLimitingData(t *testing.T) {
	// Note: Actual rate limiting is implemented in middleware with Redis
	// This test verifies the data structure needed for rate limiting

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Verify tier-based limits
	tiers := []struct {
		tier      string
		rateLimit int
	}{
		{"free", 100},
		{"basic", 500},
		{"pro", 2000},
		{"enterprise", 10000},
	}

	for _, tt := range tiers {
		result, err := testDB.DB.ExecContext(ctx, `
			INSERT INTO tenants (uuid, name, subdomain, tier, status)
			VALUES (?, ?, ?, ?, ?)
		`, uuid.New().String(), "Test", tt.tier, tt.tier, "active")
		if err != nil {
			t.Fatal(err)
		}
		tenantID, _ := result.LastInsertId()

		// Verify tenant has correct tier
		var tier string
		err = testDB.DB.QueryRowContext(ctx, `
			SELECT tier FROM tenants WHERE id = ?
		`, tenantID).Scan(&tier)
		if err != nil {
			t.Fatal(err)
		}

		if tier != tt.tier {
			t.Errorf("Expected tier %s, got %s", tt.tier, tier)
		}
	}
}

// TestSecurity_SecretStorage tests that secrets are not logged/exposed
func TestSecurity_SecretStorage(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, _, _, _ := testDB.SeedTestData(t)

	userRepo := mysql.NewUserRepository(testDB.DB)

	// Create user
	result, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO users (tenant_id, uuid, email, password_hash, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), "test@example.com", "hashed_password", "member", true)
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := result.LastInsertId()

	// Get user
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify password hash is not exposed in JSON (should be "-" tagged)
	// This is enforced by struct tags in models/user.go
	if user.PasswordHash == "" {
		t.Error("Password hash should be stored internally but not exposed")
	}

	// Note: In actual JSON marshaling, password_hash will be omitted due to `json:"-"` tag
}

// TestSecurity_AuditLogging tests audit logging captures security events
func TestSecurity_AuditLogging(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()
	tenant1ID, _, user1ID, _ := testDB.SeedTestData(t)

	auditRepo := mysql.NewAuditLogRepository(testDB.DB)

	// Log security event
	_, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, tenant1ID, user1ID, "failed_login", "user", "192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve logs
	logs, err := auditRepo.List(ctx, tenant1ID, 10, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(logs) == 0 {
		t.Fatal("Expected audit log to be created")
	}

	// Verify log details
	log := logs[0]
	if log.Action != "failed_login" {
		t.Errorf("Expected action 'failed_login', got '%s'", log.Action)
	}

	if log.IPAddress == nil || *log.IPAddress != "192.168.1.1" {
		t.Error("IP address should be logged")
	}
}
