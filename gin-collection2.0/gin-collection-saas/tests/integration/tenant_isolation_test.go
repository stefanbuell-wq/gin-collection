package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	"github.com/yourusername/gin-collection-saas/tests/testutil"
)

// TestTenantIsolation_GinRepository verifies that gins are isolated by tenant
func TestTenantIsolation_GinRepository(t *testing.T) {
	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	// Seed test data
	tenant1ID, tenant2ID, _, _ := testDB.SeedTestData(t)

	// Create gin repository
	ginRepo := mysql.NewGinRepository(testDB.DB)

	ctx := context.Background()

	// Insert gin for tenant 1
	_, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO gins (tenant_id, uuid, name, brand, country, is_available)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenant1ID, uuid.New().String(), "Gin A", "Brand A", "UK", true)
	if err != nil {
		t.Fatalf("Failed to insert gin for tenant 1: %v", err)
	}

	// Insert gin for tenant 2
	_, err = testDB.DB.ExecContext(ctx, `
		INSERT INTO gins (tenant_id, uuid, name, brand, country, is_available)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenant2ID, uuid.New().String(), "Gin B", "Brand B", "Spain", true)
	if err != nil {
		t.Fatalf("Failed to insert gin for tenant 2: %v", err)
	}

	// Test: Tenant 1 should only see their own gin
	t.Run("Tenant1_CanOnlySeeOwnGins", func(t *testing.T) {
		gins, err := ginRepo.List(ctx, tenant1ID, 10, 0, nil)
		if err != nil {
			t.Fatalf("Failed to list gins for tenant 1: %v", err)
		}

		if len(gins) != 1 {
			t.Errorf("Expected 1 gin for tenant 1, got %d", len(gins))
		}

		if len(gins) > 0 && gins[0].Name != "Gin A" {
			t.Errorf("Expected gin name 'Gin A', got '%s'", gins[0].Name)
		}
	})

	// Test: Tenant 2 should only see their own gin
	t.Run("Tenant2_CanOnlySeeOwnGins", func(t *testing.T) {
		gins, err := ginRepo.List(ctx, tenant2ID, 10, 0, nil)
		if err != nil {
			t.Fatalf("Failed to list gins for tenant 2: %v", err)
		}

		if len(gins) != 1 {
			t.Errorf("Expected 1 gin for tenant 2, got %d", len(gins))
		}

		if len(gins) > 0 && gins[0].Name != "Gin B" {
			t.Errorf("Expected gin name 'Gin B', got '%s'", gins[0].Name)
		}
	})

	// Test: Cross-tenant data leak prevention
	t.Run("CrossTenant_DataLeakPrevention", func(t *testing.T) {
		// Try to get tenant 1's gin with tenant 2's context
		tenant1Gins, _ := ginRepo.List(ctx, tenant1ID, 10, 0, nil)
		if len(tenant1Gins) == 0 {
			t.Fatal("No gins found for tenant 1")
		}

		tenant1GinID := tenant1Gins[0].ID

		// Attempt to access with wrong tenant ID
		_, err := ginRepo.GetByID(ctx, tenant2ID, tenant1GinID)
		if err == nil {
			t.Error("Expected error when accessing another tenant's gin, got nil")
		}
	})
}

// TestTenantIsolation_UserRepository verifies that users are isolated by tenant
func TestTenantIsolation_UserRepository(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	tenant1ID, tenant2ID, _, _ := testDB.SeedTestData(t)

	userRepo := mysql.NewUserRepository(testDB.DB)
	ctx := context.Background()

	// Test: List users only returns users for specific tenant
	t.Run("List_OnlyReturnsOwnTenantUsers", func(t *testing.T) {
		users1, err := userRepo.List(ctx, tenant1ID)
		if err != nil {
			t.Fatalf("Failed to list users for tenant 1: %v", err)
		}

		users2, err := userRepo.List(ctx, tenant2ID)
		if err != nil {
			t.Fatalf("Failed to list users for tenant 2: %v", err)
		}

		// Each tenant should have exactly 1 user (seeded)
		if len(users1) != 1 {
			t.Errorf("Expected 1 user for tenant 1, got %d", len(users1))
		}

		if len(users2) != 1 {
			t.Errorf("Expected 1 user for tenant 2, got %d", len(users2))
		}

		// Verify users belong to correct tenants
		if len(users1) > 0 && users1[0].TenantID != tenant1ID {
			t.Errorf("User for tenant 1 has wrong tenant_id: %d", users1[0].TenantID)
		}

		if len(users2) > 0 && users2[0].TenantID != tenant2ID {
			t.Errorf("User for tenant 2 has wrong tenant_id: %d", users2[0].TenantID)
		}
	})

	// Test: GetByEmail respects tenant isolation
	t.Run("GetByEmail_RespectsIsolation", func(t *testing.T) {
		// Get user1 from tenant1
		user1, err := userRepo.GetByEmail(ctx, tenant1ID, "user1@tenant1.com")
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}

		if user1.TenantID != tenant1ID {
			t.Errorf("User has wrong tenant_id: %d", user1.TenantID)
		}

		// Try to get user1@tenant1.com from tenant2 context (should fail)
		_, err = userRepo.GetByEmail(ctx, tenant2ID, "user1@tenant1.com")
		if err == nil {
			t.Error("Expected error when accessing user from wrong tenant, got nil")
		}
	})
}

// TestTenantIsolation_AuditLogs verifies audit logs are tenant-scoped
func TestTenantIsolation_AuditLogs(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	tenant1ID, tenant2ID, user1ID, user2ID := testDB.SeedTestData(t)

	auditRepo := mysql.NewAuditLogRepository(testDB.DB)
	ctx := context.Background()

	// Create audit logs for both tenants
	action := "test_action"
	entityType := "gin"
	entityID := int64(1)

	// Tenant 1 audit log
	_, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, entity_id)
		VALUES (?, ?, ?, ?, ?)
	`, tenant1ID, user1ID, action, entityType, entityID)
	if err != nil {
		t.Fatalf("Failed to insert audit log for tenant 1: %v", err)
	}

	// Tenant 2 audit log
	_, err = testDB.DB.ExecContext(ctx, `
		INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, entity_id)
		VALUES (?, ?, ?, ?, ?)
	`, tenant2ID, user2ID, action, entityType, entityID)
	if err != nil {
		t.Fatalf("Failed to insert audit log for tenant 2: %v", err)
	}

	// Test: Each tenant only sees their own audit logs
	t.Run("List_OnlyOwnTenantLogs", func(t *testing.T) {
		logs1, err := auditRepo.List(ctx, tenant1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list audit logs for tenant 1: %v", err)
		}

		logs2, err := auditRepo.List(ctx, tenant2ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list audit logs for tenant 2: %v", err)
		}

		if len(logs1) != 1 {
			t.Errorf("Expected 1 audit log for tenant 1, got %d", len(logs1))
		}

		if len(logs2) != 1 {
			t.Errorf("Expected 1 audit log for tenant 2, got %d", len(logs2))
		}

		// Verify tenant IDs
		if len(logs1) > 0 && logs1[0].TenantID != tenant1ID {
			t.Errorf("Audit log has wrong tenant_id: %d", logs1[0].TenantID)
		}

		if len(logs2) > 0 && logs2[0].TenantID != tenant2ID {
			t.Errorf("Audit log has wrong tenant_id: %d", logs2[0].TenantID)
		}
	})

	// Test: ListByUser respects tenant isolation
	t.Run("ListByUser_RespectsIsolation", func(t *testing.T) {
		// Get logs for user1 in tenant1
		logs, err := auditRepo.ListByUser(ctx, tenant1ID, user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list audit logs by user: %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("Expected 1 audit log, got %d", len(logs))
		}

		// Try to get logs for user1 from tenant2 context (should return 0 logs)
		logs, err = auditRepo.ListByUser(ctx, tenant2ID, user1ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list audit logs by user: %v", err)
		}

		if len(logs) != 0 {
			t.Errorf("Expected 0 audit logs when accessing with wrong tenant, got %d", len(logs))
		}
	})
}

// TestTenantIsolation_Count verifies count queries are tenant-scoped
func TestTenantIsolation_Count(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	tenant1ID, tenant2ID, _, _ := testDB.SeedTestData(t)

	ctx := context.Background()

	// Insert 3 gins for tenant 1
	for i := 0; i < 3; i++ {
		_, err := testDB.DB.ExecContext(ctx, `
			INSERT INTO gins (tenant_id, uuid, name, brand)
			VALUES (?, ?, ?, ?)
		`, tenant1ID, uuid.New().String(), "Gin", "Brand")
		if err != nil {
			t.Fatalf("Failed to insert gin: %v", err)
		}
	}

	// Insert 5 gins for tenant 2
	for i := 0; i < 5; i++ {
		_, err := testDB.DB.ExecContext(ctx, `
			INSERT INTO gins (tenant_id, uuid, name, brand)
			VALUES (?, ?, ?, ?)
		`, tenant2ID, uuid.New().String(), "Gin", "Brand")
		if err != nil {
			t.Fatalf("Failed to insert gin: %v", err)
		}
	}

	userRepo := mysql.NewUserRepository(testDB.DB)

	// Test: Count returns correct tenant-scoped counts
	t.Run("Count_ReturnsCorrectTenantCount", func(t *testing.T) {
		count1, err := userRepo.CountByTenant(ctx, tenant1ID)
		if err != nil {
			t.Fatalf("Failed to count for tenant 1: %v", err)
		}

		count2, err := userRepo.CountByTenant(ctx, tenant2ID)
		if err != nil {
			t.Fatalf("Failed to count for tenant 2: %v", err)
		}

		// Each tenant has 1 seeded user
		if count1 != 1 {
			t.Errorf("Expected count 1 for tenant 1, got %d", count1)
		}

		if count2 != 1 {
			t.Errorf("Expected count 1 for tenant 2, got %d", count2)
		}
	})
}
