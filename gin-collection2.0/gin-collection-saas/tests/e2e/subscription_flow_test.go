package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	subscriptionUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/subscription"
	"github.com/yourusername/gin-collection-saas/tests/testutil"
)

// MockPayPalClient implements a mock PayPal client for testing
type MockPayPalClient struct {
	subscriptionID string
	shouldFail     bool
}

func (m *MockPayPalClient) CreateSubscription(req interface{}) (*struct {
	ID          string
	Status      string
	ApprovalURL string
}, error) {
	if m.shouldFail {
		return nil, &struct{ error }{error: &struct{ string }{"PayPal error"}}
	}

	m.subscriptionID = "I-" + uuid.New().String()[:8]
	return &struct {
		ID          string
		Status      string
		ApprovalURL string
	}{
		ID:          m.subscriptionID,
		Status:      "APPROVAL_PENDING",
		ApprovalURL: "https://paypal.com/approve/" + m.subscriptionID,
	}, nil
}

func (m *MockPayPalClient) GetSubscription(id string) (*struct {
	ID     string
	Status string
}, error) {
	return &struct {
		ID     string
		Status string
	}{
		ID:     id,
		Status: "ACTIVE",
	}, nil
}

func (m *MockPayPalClient) CancelSubscription(id string) error {
	return nil
}

// TestE2E_SubscriptionUpgradeFlow tests the complete subscription upgrade flow
func TestE2E_SubscriptionUpgradeFlow(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Create subscriptions table
	_, err := testDB.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS subscriptions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			tenant_id BIGINT NOT NULL,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			plan_id VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL,
			billing_cycle VARCHAR(20) NOT NULL,
			paypal_subscription_id VARCHAR(100),
			paypal_plan_id VARCHAR(100),
			amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'USD',
			current_period_start TIMESTAMP NULL,
			current_period_end TIMESTAMP NULL,
			next_billing_date TIMESTAMP NULL,
			cancelled_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_tenant_id (tenant_id)
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Setup repositories
	tenantRepo := mysql.NewTenantRepository(testDB.DB)
	subscriptionRepo := mysql.NewSubscriptionRepository(testDB.DB)

	// Create test tenant (Free tier)
	result, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO tenants (uuid, name, subdomain, tier, status)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), "Test Tenant", "test", "free", "active")
	if err != nil {
		t.Fatal(err)
	}
	tenantID, _ := result.LastInsertId()

	// Create mock PayPal client
	mockPayPal := &MockPayPalClient{}

	// Create subscription service
	subscriptionService := subscriptionUsecase.NewService(
		subscriptionRepo,
		tenantRepo,
		mockPayPal,
		"http://localhost:8080",
	)

	// Test Step 1: Initiate upgrade from Free to Pro
	t.Run("Step1_InitiateUpgrade", func(t *testing.T) {
		upgradeResp, err := subscriptionService.InitiateUpgrade(ctx, tenantID, "pro", models.BillingCycleMonthly)
		if err != nil {
			t.Fatalf("Failed to initiate upgrade: %v", err)
		}

		if upgradeResp.ApprovalURL == "" {
			t.Error("Expected approval URL, got empty string")
		}

		if upgradeResp.SubscriptionID == "" {
			t.Error("Expected subscription ID, got empty string")
		}

		// Verify subscription was created in database
		sub, err := subscriptionRepo.GetByTenantID(ctx, tenantID)
		if err != nil {
			t.Fatalf("Failed to get subscription: %v", err)
		}

		if sub.Status != models.SubscriptionStatusPending {
			t.Errorf("Expected status 'pending', got '%s'", sub.Status)
		}

		if sub.PlanID != "pro" {
			t.Errorf("Expected plan 'pro', got '%s'", sub.PlanID)
		}
	})

	// Test Step 2: User approves on PayPal (simulated)
	// In real flow, user is redirected to PayPal, approves, then redirected back

	// Test Step 3: Activate subscription after approval
	t.Run("Step3_ActivateSubscription", func(t *testing.T) {
		// Get subscription
		sub, err := subscriptionRepo.GetByTenantID(ctx, tenantID)
		if err != nil {
			t.Fatal(err)
		}

		// Activate subscription
		err = subscriptionService.ActivateSubscription(ctx, tenantID, sub.UUID)
		if err != nil {
			t.Fatalf("Failed to activate subscription: %v", err)
		}

		// Verify subscription is now active
		sub, err = subscriptionRepo.GetByTenantID(ctx, tenantID)
		if err != nil {
			t.Fatal(err)
		}

		if sub.Status != models.SubscriptionStatusActive {
			t.Errorf("Expected status 'active', got '%s'", sub.Status)
		}

		// Verify tenant tier was upgraded
		tenant, err := tenantRepo.GetByID(ctx, tenantID)
		if err != nil {
			t.Fatal(err)
		}

		if tenant.Tier != "pro" {
			t.Errorf("Expected tenant tier 'pro', got '%s'", tenant.Tier)
		}
	})

	// Test Step 4: Handle webhook event (payment successful)
	t.Run("Step4_HandleWebhook_PaymentSuccess", func(t *testing.T) {
		sub, _ := subscriptionRepo.GetByTenantID(ctx, tenantID)

		// Simulate PayPal webhook event
		webhookEvent := &subscriptionUsecase.WebhookEvent{
			EventType: "PAYMENT.SALE.COMPLETED",
			Resource: map[string]interface{}{
				"billing_agreement_id": sub.PayPalSubscriptionID,
				"amount": map[string]interface{}{
					"total": "5.99",
				},
			},
		}

		err := subscriptionService.HandleWebhookEvent(ctx, webhookEvent)
		if err != nil {
			t.Errorf("Failed to handle webhook: %v", err)
		}
	})

	// Test Step 5: Cancel subscription
	t.Run("Step5_CancelSubscription", func(t *testing.T) {
		err := subscriptionService.CancelSubscription(ctx, tenantID)
		if err != nil {
			t.Fatalf("Failed to cancel subscription: %v", err)
		}

		// Verify subscription is cancelled
		sub, err := subscriptionRepo.GetByTenantID(ctx, tenantID)
		if err != nil {
			t.Fatal(err)
		}

		if sub.Status != models.SubscriptionStatusCancelled {
			t.Errorf("Expected status 'cancelled', got '%s'", sub.Status)
		}

		if sub.CancelledAt == nil {
			t.Error("Expected cancelled_at to be set")
		}

		// Note: Tenant tier stays at 'pro' until subscription expires
		// (Grace period until end of billing period)
	})
}

// TestE2E_SubscriptionFailureHandling tests error scenarios
func TestE2E_SubscriptionFailureHandling(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Create subscriptions table
	_, err := testDB.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS subscriptions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			tenant_id BIGINT NOT NULL,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			plan_id VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL,
			billing_cycle VARCHAR(20) NOT NULL,
			paypal_subscription_id VARCHAR(100),
			amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'USD',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	tenantRepo := mysql.NewTenantRepository(testDB.DB)
	subscriptionRepo := mysql.NewSubscriptionRepository(testDB.DB)

	// Create tenant
	result, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO tenants (uuid, name, subdomain, tier, status)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), "Test", "test", "free", "active")
	if err != nil {
		t.Fatal(err)
	}
	tenantID, _ := result.LastInsertId()

	// Test: PayPal API failure
	t.Run("PayPalFailure_ReturnsError", func(t *testing.T) {
		mockPayPal := &MockPayPalClient{shouldFail: true}
		service := subscriptionUsecase.NewService(subscriptionRepo, tenantRepo, mockPayPal, "http://localhost")

		_, err := service.InitiateUpgrade(ctx, tenantID, "pro", models.BillingCycleMonthly)
		if err == nil {
			t.Error("Expected error when PayPal fails, got nil")
		}
	})

	// Test: Invalid plan ID
	t.Run("InvalidPlanID_ReturnsError", func(t *testing.T) {
		mockPayPal := &MockPayPalClient{}
		service := subscriptionUsecase.NewService(subscriptionRepo, tenantRepo, mockPayPal, "http://localhost")

		_, err := service.InitiateUpgrade(ctx, tenantID, "invalid_plan", models.BillingCycleMonthly)
		if err == nil {
			t.Error("Expected error for invalid plan, got nil")
		}
	})
}

// TestE2E_SubscriptionRenewal tests subscription renewal flow
func TestE2E_SubscriptionRenewal(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.RunMigrations(t)

	ctx := context.Background()

	// Create subscriptions table
	_, err := testDB.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS subscriptions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			tenant_id BIGINT NOT NULL,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			plan_id VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL,
			billing_cycle VARCHAR(20) NOT NULL,
			paypal_subscription_id VARCHAR(100),
			amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'USD',
			current_period_start TIMESTAMP NULL,
			current_period_end TIMESTAMP NULL,
			next_billing_date TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create active subscription that's about to renew
	result, err := testDB.DB.ExecContext(ctx, `
		INSERT INTO tenants (uuid, name, subdomain, tier, status)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), "Test", "test", "pro", "active")
	if err != nil {
		t.Fatal(err)
	}
	tenantID, _ := result.LastInsertId()

	now := time.Now()
	periodEnd := now.Add(24 * time.Hour) // Expires tomorrow
	nextBilling := periodEnd

	_, err = testDB.DB.ExecContext(ctx, `
		INSERT INTO subscriptions (tenant_id, uuid, plan_id, status, billing_cycle, paypal_subscription_id, amount, currency, current_period_end, next_billing_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tenantID, uuid.New().String(), "pro", "active", "monthly", "I-TEST123", 5.99, "USD", periodEnd, nextBilling)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("HandleRenewalWebhook", func(t *testing.T) {
		subscriptionRepo := mysql.NewSubscriptionRepository(testDB.DB)
		tenantRepo := mysql.NewTenantRepository(testDB.DB)
		mockPayPal := &MockPayPalClient{}
		service := subscriptionUsecase.NewService(subscriptionRepo, tenantRepo, mockPayPal, "http://localhost")

		// Simulate renewal webhook
		webhookEvent := &subscriptionUsecase.WebhookEvent{
			EventType: "BILLING.SUBSCRIPTION.RENEWED",
			Resource: map[string]interface{}{
				"id": "I-TEST123",
			},
		}

		err := service.HandleWebhookEvent(ctx, webhookEvent)
		if err != nil {
			t.Errorf("Failed to handle renewal webhook: %v", err)
		}

		// Verify subscription is still active
		sub, err := subscriptionRepo.GetByPayPalSubscriptionID(ctx, "I-TEST123")
		if err != nil {
			t.Fatal(err)
		}

		if sub.Status != models.SubscriptionStatusActive {
			t.Errorf("Expected status 'active' after renewal, got '%s'", sub.Status)
		}
	})
}
