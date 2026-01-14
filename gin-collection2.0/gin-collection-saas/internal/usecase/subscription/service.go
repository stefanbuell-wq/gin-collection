package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/external"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles subscription business logic
type Service struct {
	subscriptionRepo repositories.SubscriptionRepository
	tenantRepo       repositories.TenantRepository
	paypalClient     *external.PayPalClient
	baseURL          string
}

// NewService creates a new subscription service
func NewService(
	subscriptionRepo repositories.SubscriptionRepository,
	tenantRepo repositories.TenantRepository,
	paypalClient *external.PayPalClient,
	baseURL string,
) *Service {
	return &Service{
		subscriptionRepo: subscriptionRepo,
		tenantRepo:       tenantRepo,
		paypalClient:     paypalClient,
		baseURL:          baseURL,
	}
}

// GetCurrentSubscription retrieves the current subscription for a tenant
func (s *Service) GetCurrentSubscription(ctx context.Context, tenantID int64) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		// If no subscription exists, tenant is on Free tier
		return nil, nil
	}

	return subscription, nil
}

// GetAvailablePlans returns all available subscription plans
func (s *Service) GetAvailablePlans() []models.SubscriptionPlan {
	return models.AvailablePlans
}

// InitiateUpgrade initiates a subscription upgrade to a new plan
func (s *Service) InitiateUpgrade(ctx context.Context, tenantID int64, planID string, billingCycle models.BillingCycle) (*UpgradeResponse, error) {
	logger.Info("Initiating subscription upgrade", "tenant_id", tenantID, "plan_id", planID, "billing_cycle", billingCycle)

	// Verify tenant exists
	_, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Validate plan
	plan := models.GetPlanByID(planID)
	if plan == nil {
		return nil, fmt.Errorf("invalid plan ID: %s", planID)
	}

	// Determine tier from plan ID
	tier := s.getTierFromPlanID(planID)

	// Get PayPal plan ID based on billing cycle
	paypalPlanID := s.getPayPalPlanID(planID, billingCycle)
	if paypalPlanID == "" {
		return nil, fmt.Errorf("no PayPal plan configured for %s %s", planID, billingCycle)
	}

	// Create PayPal subscription
	paypalReq := &external.PayPalSubscriptionRequest{
		PlanID: paypalPlanID,
		ApplicationContext: &external.PayPalApplicationContext{
			BrandName: "Gin Collection SaaS",
			ReturnURL: fmt.Sprintf("%s/subscription/success", s.baseURL),
			CancelURL: fmt.Sprintf("%s/subscription/cancel", s.baseURL),
			UserAction: "SUBSCRIBE_NOW",
			PaymentMethod: &external.PayPalPaymentMethod{
				PayeePreferred: "IMMEDIATE_PAYMENT_REQUIRED",
			},
		},
	}

	paypalResp, err := s.paypalClient.CreateSubscription(paypalReq)
	if err != nil {
		logger.Error("Failed to create PayPal subscription", "error", err.Error())
		return nil, fmt.Errorf("failed to create PayPal subscription: %w", err)
	}

	// Calculate amount based on billing cycle
	amount := plan.PriceMonthly
	if billingCycle == models.BillingCycleYearly {
		amount = plan.PriceYearly
	}

	// Create subscription record
	subscription := &models.Subscription{
		TenantID:              tenantID,
		PlanID:                planID,
		Status:                models.SubscriptionStatusPending,
		BillingCycle:          billingCycle,
		PayPalSubscriptionID:  &paypalResp.ID,
		PayPalPlanID:          &paypalPlanID,
		Amount:                amount,
		Currency:              "EUR",
		CurrentPeriodStart:    nil, // Set when activated
		CurrentPeriodEnd:      nil,
		NextBillingDate:       nil,
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		logger.Error("Failed to create subscription record", "error", err.Error())
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	logger.Info("Subscription created", "subscription_id", subscription.ID, "paypal_subscription_id", paypalResp.ID)

	return &UpgradeResponse{
		SubscriptionID:    subscription.ID,
		PayPalApprovalURL: paypalResp.GetApprovalURL(),
		Plan:              plan,
		Tier:              tier,
		Amount:            amount,
		Currency:          "EUR",
		BillingCycle:      billingCycle,
	}, nil
}

// ActivateSubscription activates a subscription after PayPal approval
func (s *Service) ActivateSubscription(ctx context.Context, paypalSubscriptionID string) error {
	logger.Info("Activating subscription", "paypal_subscription_id", paypalSubscriptionID)

	// Get subscription from DB
	subscription, err := s.subscriptionRepo.GetByPayPalSubscriptionID(ctx, paypalSubscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Get PayPal subscription details
	paypalSub, err := s.paypalClient.GetSubscription(paypalSubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get PayPal subscription: %w", err)
	}

	// Parse timestamps
	startTime, _ := time.Parse(time.RFC3339, paypalSub.StartTime)

	var nextBillingTime *time.Time
	if paypalSub.BillingInfo != nil && paypalSub.BillingInfo.NextBillingTime != "" {
		t, _ := time.Parse(time.RFC3339, paypalSub.BillingInfo.NextBillingTime)
		nextBillingTime = &t
	}

	// Calculate period end
	var periodEnd *time.Time
	if subscription.BillingCycle == models.BillingCycleMonthly {
		t := startTime.AddDate(0, 1, 0)
		periodEnd = &t
	} else if subscription.BillingCycle == models.BillingCycleYearly {
		t := startTime.AddDate(1, 0, 0)
		periodEnd = &t
	}

	// Update subscription
	subscription.Status = models.SubscriptionStatusActive
	subscription.CurrentPeriodStart = &startTime
	subscription.CurrentPeriodEnd = periodEnd
	subscription.NextBillingDate = nextBillingTime

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Update tenant tier
	tier := s.getTierFromPlanID(subscription.PlanID)
	tenant, err := s.tenantRepo.GetByID(ctx, subscription.TenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	tenant.Tier = tier
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return fmt.Errorf("failed to update tenant tier: %w", err)
	}

	logger.Info("Subscription activated", "subscription_id", subscription.ID, "tenant_id", subscription.TenantID, "tier", tier)

	return nil
}

// CancelSubscription cancels a subscription
func (s *Service) CancelSubscription(ctx context.Context, tenantID int64, reason string) error {
	logger.Info("Cancelling subscription", "tenant_id", tenantID, "reason", reason)

	// Get active subscription
	subscription, err := s.subscriptionRepo.GetActiveSubscription(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("no active subscription found: %w", err)
	}

	// Cancel in PayPal
	if subscription.PayPalSubscriptionID != nil {
		if err := s.paypalClient.CancelSubscription(*subscription.PayPalSubscriptionID, reason); err != nil {
			logger.Error("Failed to cancel PayPal subscription", "error", err.Error())
			// Continue anyway to update local status
		}
	}

	// Update subscription status
	now := time.Now()
	subscription.Status = models.SubscriptionStatusCancelled
	subscription.CancelledAt = &now

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Downgrade tenant to Free tier
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	tenant.Tier = models.TierFree
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return fmt.Errorf("failed to downgrade tenant: %w", err)
	}

	logger.Info("Subscription cancelled", "subscription_id", subscription.ID, "tenant_id", tenantID)

	return nil
}

// HandleWebhookEvent processes PayPal webhook events
func (s *Service) HandleWebhookEvent(ctx context.Context, event *WebhookEvent) error {
	logger.Info("Processing PayPal webhook", "event_type", event.EventType, "resource_id", event.Resource.ID)

	switch event.EventType {
	case "BILLING.SUBSCRIPTION.ACTIVATED":
		return s.handleSubscriptionActivated(ctx, event)
	case "BILLING.SUBSCRIPTION.UPDATED":
		return s.handleSubscriptionUpdated(ctx, event)
	case "BILLING.SUBSCRIPTION.CANCELLED":
		return s.handleSubscriptionCancelled(ctx, event)
	case "BILLING.SUBSCRIPTION.SUSPENDED":
		return s.handleSubscriptionSuspended(ctx, event)
	case "BILLING.SUBSCRIPTION.EXPIRED":
		return s.handleSubscriptionExpired(ctx, event)
	case "PAYMENT.SALE.COMPLETED":
		return s.handlePaymentCompleted(ctx, event)
	default:
		logger.Debug("Unhandled webhook event type", "event_type", event.EventType)
		return nil
	}
}

// getTierFromPlanID maps plan ID to subscription tier
func (s *Service) getTierFromPlanID(planID string) models.SubscriptionTier {
	switch planID {
	case "PLAN_BASIC_MONTHLY", "PLAN_BASIC_YEARLY":
		return models.TierBasic
	case "PLAN_PRO_MONTHLY", "PLAN_PRO_YEARLY":
		return models.TierPro
	case "PLAN_ENTERPRISE":
		return models.TierEnterprise
	default:
		return models.TierFree
	}
}

// getPayPalPlanID returns the PayPal plan ID for a given plan and billing cycle
func (s *Service) getPayPalPlanID(planID string, billingCycle models.BillingCycle) string {
	// In production, these should come from config/environment
	// For now, return placeholder IDs
	switch planID {
	case "PLAN_BASIC_MONTHLY":
		return "P-BASIC-MONTHLY"
	case "PLAN_BASIC_YEARLY":
		return "P-BASIC-YEARLY"
	case "PLAN_PRO_MONTHLY":
		return "P-PRO-MONTHLY"
	case "PLAN_PRO_YEARLY":
		return "P-PRO-YEARLY"
	case "PLAN_ENTERPRISE":
		return "P-ENTERPRISE"
	default:
		return ""
	}
}

// handleSubscriptionActivated handles subscription activation webhook
func (s *Service) handleSubscriptionActivated(ctx context.Context, event *WebhookEvent) error {
	return s.ActivateSubscription(ctx, event.Resource.ID)
}

// handleSubscriptionUpdated handles subscription update webhook
func (s *Service) handleSubscriptionUpdated(ctx context.Context, event *WebhookEvent) error {
	subscription, err := s.subscriptionRepo.GetByPayPalSubscriptionID(ctx, event.Resource.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Get updated details from PayPal
	paypalSub, err := s.paypalClient.GetSubscription(event.Resource.ID)
	if err != nil {
		return fmt.Errorf("failed to get PayPal subscription: %w", err)
	}

	// Update next billing date
	if paypalSub.BillingInfo != nil && paypalSub.BillingInfo.NextBillingTime != "" {
		t, _ := time.Parse(time.RFC3339, paypalSub.BillingInfo.NextBillingTime)
		subscription.NextBillingDate = &t
	}

	return s.subscriptionRepo.Update(ctx, subscription)
}

// handleSubscriptionCancelled handles subscription cancellation webhook
func (s *Service) handleSubscriptionCancelled(ctx context.Context, event *WebhookEvent) error {
	subscription, err := s.subscriptionRepo.GetByPayPalSubscriptionID(ctx, event.Resource.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	now := time.Now()
	subscription.Status = models.SubscriptionStatusCancelled
	subscription.CancelledAt = &now

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return err
	}

	// Downgrade tenant
	tenant, _ := s.tenantRepo.GetByID(ctx, subscription.TenantID)
	if tenant != nil {
		tenant.Tier = models.TierFree
		s.tenantRepo.Update(ctx, tenant)
	}

	return nil
}

// handleSubscriptionSuspended handles subscription suspension webhook
func (s *Service) handleSubscriptionSuspended(ctx context.Context, event *WebhookEvent) error {
	subscription, err := s.subscriptionRepo.GetByPayPalSubscriptionID(ctx, event.Resource.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = models.SubscriptionStatusSuspended
	return s.subscriptionRepo.Update(ctx, subscription)
}

// handleSubscriptionExpired handles subscription expiration webhook
func (s *Service) handleSubscriptionExpired(ctx context.Context, event *WebhookEvent) error {
	subscription, err := s.subscriptionRepo.GetByPayPalSubscriptionID(ctx, event.Resource.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = models.SubscriptionStatusExpired

	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return err
	}

	// Downgrade tenant
	tenant, _ := s.tenantRepo.GetByID(ctx, subscription.TenantID)
	if tenant != nil {
		tenant.Tier = models.TierFree
		s.tenantRepo.Update(ctx, tenant)
	}

	return nil
}

// handlePaymentCompleted handles successful payment webhook
func (s *Service) handlePaymentCompleted(ctx context.Context, event *WebhookEvent) error {
	// Log successful payment for accounting/monitoring
	logger.Info("Payment completed", "billing_agreement_id", event.Resource.BillingAgreementID, "amount", event.Resource.Amount)
	return nil
}

// UpgradeResponse represents the response for an upgrade request
type UpgradeResponse struct {
	SubscriptionID    int64                    `json:"subscription_id"`
	PayPalApprovalURL string                   `json:"paypal_approval_url"`
	Plan              *models.SubscriptionPlan `json:"plan"`
	Tier              models.SubscriptionTier  `json:"tier"`
	Amount            float64                  `json:"amount"`
	Currency          string                   `json:"currency"`
	BillingCycle      models.BillingCycle      `json:"billing_cycle"`
}

// WebhookEvent represents a PayPal webhook event
type WebhookEvent struct {
	ID         string             `json:"id"`
	EventType  string             `json:"event_type"`
	CreateTime string             `json:"create_time"`
	Resource   WebhookResource    `json:"resource"`
}

// WebhookResource represents the resource in a webhook event
type WebhookResource struct {
	ID                  string  `json:"id"`
	Status              string  `json:"status"`
	BillingAgreementID  string  `json:"billing_agreement_id"`
	Amount              *Amount `json:"amount,omitempty"`
}

// Amount represents payment amount
type Amount struct {
	Total    string `json:"total"`
	Currency string `json:"currency"`
}
