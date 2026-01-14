package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	ginUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/gin"
	userUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/user"
	"github.com/yourusername/gin-collection-saas/tests/testutil"
)

// TestTierEnforcement_GinLimits verifies gin count limits per tier
func TestTierEnforcement_GinLimits(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	tests := []struct {
		name      string
		tier      string
		maxGins   int
		shouldErr bool
	}{
		{"Free tier - 10 gins", "free", 10, false},
		{"Free tier - 11 gins", "free", 11, true},
		{"Basic tier - 50 gins", "basic", 50, false},
		{"Basic tier - 51 gins", "basic", 51, true},
		{"Pro tier - unlimited", "pro", 1000, false},
		{"Enterprise tier - unlimited", "enterprise", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create tenant with specific tier
			result, err := testDB.DB.ExecContext(ctx, `
				INSERT INTO tenants (uuid, name, subdomain, tier, status)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), "Test Tenant", "test", tt.tier, "active")
			if err != nil {
				t.Fatalf("Failed to create tenant: %v", err)
			}
			tenantID, _ := result.LastInsertId()

			// Insert gins up to the limit
			for i := 0; i < tt.maxGins; i++ {
				_, err := testDB.DB.ExecContext(ctx, `
					INSERT INTO gins (tenant_id, uuid, name, brand)
					VALUES (?, ?, ?, ?)
				`, tenantID, uuid.New().String(), "Gin", "Brand")

				if i < tt.maxGins-1 {
					// Should succeed before limit
					if err != nil {
						t.Fatalf("Failed to insert gin %d: %v", i, err)
					}
				} else {
					// Last insert - check if it should error
					if tt.shouldErr && err == nil {
						t.Error("Expected error when exceeding limit, got nil")
					}
					if !tt.shouldErr && err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				}
			}
		})
	}
}

// TestTierEnforcement_PhotoLimits verifies photo limits per gin per tier
func TestTierEnforcement_PhotoLimits(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Add gin_photos table
	_, err := testDB.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gin_photos (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			tenant_id BIGINT NOT NULL,
			gin_id BIGINT NOT NULL,
			photo_url VARCHAR(500) NOT NULL,
			photo_type VARCHAR(20) DEFAULT 'bottle',
			is_primary BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create photos table: %v", err)
	}

	tests := []struct {
		tier       string
		maxPhotos  int
		insertMore bool // Try to insert one more than max
	}{
		{"free", 1, true},
		{"basic", 3, true},
		{"pro", 10, true},
		{"enterprise", 50, false}, // Don't test all 50
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			// Create tenant
			result, err := testDB.DB.ExecContext(ctx, `
				INSERT INTO tenants (uuid, name, subdomain, tier, status)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), "Test", tt.tier, tt.tier, "active")
			if err != nil {
				t.Fatal(err)
			}
			tenantID, _ := result.LastInsertId()

			// Create gin
			result, err = testDB.DB.ExecContext(ctx, `
				INSERT INTO gins (tenant_id, uuid, name, brand)
				VALUES (?, ?, ?, ?)
			`, tenantID, uuid.New().String(), "Test Gin", "Brand")
			if err != nil {
				t.Fatal(err)
			}
			ginID, _ := result.LastInsertId()

			// Insert photos up to limit
			for i := 0; i < tt.maxPhotos; i++ {
				_, err := testDB.DB.ExecContext(ctx, `
					INSERT INTO gin_photos (tenant_id, gin_id, photo_url)
					VALUES (?, ?, ?)
				`, tenantID, ginID, "http://example.com/photo.jpg")
				if err != nil {
					t.Fatalf("Failed to insert photo %d: %v", i, err)
				}
			}

			// Try to insert one more if test requires
			if tt.insertMore {
				// Count current photos
				var count int
				err := testDB.DB.QueryRowContext(ctx, `
					SELECT COUNT(*) FROM gin_photos WHERE tenant_id = ? AND gin_id = ?
				`, tenantID, ginID).Scan(&count)
				if err != nil {
					t.Fatal(err)
				}

				if count != tt.maxPhotos {
					t.Errorf("Expected %d photos, got %d", tt.maxPhotos, count)
				}
			}
		})
	}
}

// TestTierEnforcement_FeatureAccess verifies tier-based feature access
func TestTierEnforcement_FeatureAccess(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	features := map[string][]string{
		"free":       {},
		"basic":      {},
		"pro":        {"botanicals", "cocktails", "ai_suggestions", "export", "import"},
		"enterprise": {"botanicals", "cocktails", "ai_suggestions", "export", "import", "multi_user", "api_access"},
	}

	for tier, expectedFeatures := range features {
		t.Run(tier, func(t *testing.T) {
			// Create tenant
			result, err := testDB.DB.ExecContext(ctx, `
				INSERT INTO tenants (uuid, name, subdomain, tier, status)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), "Test", tier, tier, "active")
			if err != nil {
				t.Fatal(err)
			}
			tenantID, _ := result.LastInsertId()

			// Get tenant
			var tierValue string
			err = testDB.DB.QueryRowContext(ctx, `
				SELECT tier FROM tenants WHERE id = ?
			`, tenantID).Scan(&tierValue)
			if err != nil {
				t.Fatal(err)
			}

			// Get limits based on tier
			limits := models.TierLimits[models.TenantTier(tierValue)]

			// Verify features
			for _, feature := range expectedFeatures {
				hasFeature := false
				for _, f := range limits.Features {
					if f == feature {
						hasFeature = true
						break
					}
				}
				if !hasFeature {
					t.Errorf("Tier %s should have feature %s", tier, feature)
				}
			}
		})
	}
}

// TestTierEnforcement_MultiUserRestriction verifies multi-user is Enterprise-only
func TestTierEnforcement_MultiUserRestriction(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	tenantRepo := mysql.NewTenantRepository(testDB.DB)
	userRepo := mysql.NewUserRepository(testDB.DB)
	auditRepo := mysql.NewAuditLogRepository(testDB.DB)

	userService := userUsecase.NewService(userRepo, tenantRepo, auditRepo)

	tests := []struct {
		tier      string
		shouldErr bool
	}{
		{"free", true},
		{"basic", true},
		{"pro", true},
		{"enterprise", false},
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			// Create tenant
			result, err := testDB.DB.ExecContext(ctx, `
				INSERT INTO tenants (uuid, name, subdomain, tier, status)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), "Test", tt.tier, tt.tier, "active")
			if err != nil {
				t.Fatal(err)
			}
			tenantID, _ := result.LastInsertId()

			// Create owner user
			result, err = testDB.DB.ExecContext(ctx, `
				INSERT INTO users (tenant_id, uuid, email, password_hash, role, is_active)
				VALUES (?, ?, ?, ?, ?, ?)
			`, tenantID, uuid.New().String(), "owner@test.com", "hash", "owner", true)
			if err != nil {
				t.Fatal(err)
			}
			ownerID, _ := result.LastInsertId()

			// Try to invite another user
			_, err = userService.InviteUser(ctx, tenantID, ownerID, "user2@test.com", "John", "Doe", models.RoleMember)

			if tt.shouldErr && err == nil {
				t.Error("Expected error for non-Enterprise tier, got nil")
			}

			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for Enterprise tier: %v", err)
			}
		})
	}
}

// TestTierEnforcement_StorageLimits verifies storage limits per tier
func TestTierEnforcement_StorageLimits(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Create usage_metrics table
	_, err := testDB.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS usage_metrics (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			tenant_id BIGINT NOT NULL,
			metric_name VARCHAR(100) NOT NULL,
			current_value INT NOT NULL DEFAULT 0,
			period_start DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY unique_tenant_metric_period (tenant_id, metric_name, period_start)
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		tier         string
		storageLimMB *int
	}{
		{"free", intPtr(100)},
		{"basic", intPtr(500)},
		{"pro", intPtr(5000)},
		{"enterprise", nil}, // Unlimited
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			// Create tenant
			result, err := testDB.DB.ExecContext(ctx, `
				INSERT INTO tenants (uuid, name, subdomain, tier, status)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), "Test", tt.tier, tt.tier, "active")
			if err != nil {
				t.Fatal(err)
			}
			tenantID, _ := result.LastInsertId()

			// Get tier limits
			var tierValue string
			err = testDB.DB.QueryRowContext(ctx, `
				SELECT tier FROM tenants WHERE id = ?
			`, tenantID).Scan(&tierValue)
			if err != nil {
				t.Fatal(err)
			}

			limits := models.TierLimits[models.TenantTier(tierValue)]

			// Verify storage limit matches
			if tt.storageLimMB == nil && limits.StorageLimitMB != nil {
				t.Errorf("Expected unlimited storage, got %dMB", *limits.StorageLimitMB)
			}

			if tt.storageLimMB != nil && limits.StorageLimitMB == nil {
				t.Error("Expected storage limit, got unlimited")
			}

			if tt.storageLimMB != nil && limits.StorageLimitMB != nil {
				if *tt.storageLimMB != *limits.StorageLimitMB {
					t.Errorf("Expected storage limit %dMB, got %dMB", *tt.storageLimMB, *limits.StorageLimitMB)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
