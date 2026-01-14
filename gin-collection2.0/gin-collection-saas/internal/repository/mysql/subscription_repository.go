package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// SubscriptionRepository implements subscription data access
type SubscriptionRepository struct {
	db *sql.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		subscription.TenantID,
		uuid.New().String(),
		subscription.PlanID,
		subscription.Status,
		subscription.BillingCycle,
		subscription.PayPalSubscriptionID,
		subscription.PayPalPlanID,
		subscription.Amount,
		subscription.Currency,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.NextBillingDate,
		subscription.CancelledAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get subscription ID: %w", err)
	}

	subscription.ID = id
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	return nil
}

// GetByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (*models.Subscription, error) {
	query := `
		SELECT
			id, tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE id = ?
	`

	subscription := &models.Subscription{}
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.UUID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.BillingCycle,
		&subscription.PayPalSubscriptionID,
		&subscription.PayPalPlanID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.NextBillingDate,
		&cancelledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if cancelledAt.Valid {
		subscription.CancelledAt = &cancelledAt.Time
	}

	return subscription, nil
}

// GetByTenantID retrieves the current subscription for a tenant
func (r *SubscriptionRepository) GetByTenantID(ctx context.Context, tenantID int64) (*models.Subscription, error) {
	query := `
		SELECT
			id, tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	subscription := &models.Subscription{}
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.UUID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.BillingCycle,
		&subscription.PayPalSubscriptionID,
		&subscription.PayPalPlanID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.NextBillingDate,
		&cancelledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if cancelledAt.Valid {
		subscription.CancelledAt = &cancelledAt.Time
	}

	return subscription, nil
}

// GetByPayPalSubscriptionID retrieves a subscription by PayPal subscription ID
func (r *SubscriptionRepository) GetByPayPalSubscriptionID(ctx context.Context, paypalSubscriptionID string) (*models.Subscription, error) {
	query := `
		SELECT
			id, tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE paypal_subscription_id = ?
	`

	subscription := &models.Subscription{}
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, paypalSubscriptionID).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.UUID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.BillingCycle,
		&subscription.PayPalSubscriptionID,
		&subscription.PayPalPlanID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.NextBillingDate,
		&cancelledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if cancelledAt.Valid {
		subscription.CancelledAt = &cancelledAt.Time
	}

	return subscription, nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET plan_id = ?, status = ?, billing_cycle = ?,
		    paypal_subscription_id = ?, paypal_plan_id = ?,
		    amount = ?, currency = ?,
		    current_period_start = ?, current_period_end = ?,
		    next_billing_date = ?, cancelled_at = ?,
		    updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		subscription.PlanID,
		subscription.Status,
		subscription.BillingCycle,
		subscription.PayPalSubscriptionID,
		subscription.PayPalPlanID,
		subscription.Amount,
		subscription.Currency,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.NextBillingDate,
		subscription.CancelledAt,
		subscription.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	subscription.UpdatedAt = time.Now()

	return nil
}

// UpdateStatus updates subscription status
func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, id int64, status models.SubscriptionStatus) error {
	query := `
		UPDATE subscriptions
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

// List retrieves all subscriptions for a tenant
func (r *SubscriptionRepository) List(ctx context.Context, tenantID int64) ([]*models.Subscription, error) {
	query := `
		SELECT
			id, tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription

	for rows.Next() {
		subscription := &models.Subscription{}
		var cancelledAt sql.NullTime

		err := rows.Scan(
			&subscription.ID,
			&subscription.TenantID,
			&subscription.UUID,
			&subscription.PlanID,
			&subscription.Status,
			&subscription.BillingCycle,
			&subscription.PayPalSubscriptionID,
			&subscription.PayPalPlanID,
			&subscription.Amount,
			&subscription.Currency,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.NextBillingDate,
			&cancelledAt,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		if cancelledAt.Valid {
			subscription.CancelledAt = &cancelledAt.Time
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// GetActiveSubscription retrieves the active subscription for a tenant
func (r *SubscriptionRepository) GetActiveSubscription(ctx context.Context, tenantID int64) (*models.Subscription, error) {
	query := `
		SELECT
			id, tenant_id, uuid, plan_id, status, billing_cycle,
			paypal_subscription_id, paypal_plan_id, amount, currency,
			current_period_start, current_period_end,
			next_billing_date, cancelled_at, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = ? AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	subscription := &models.Subscription{}
	var cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&subscription.ID,
		&subscription.TenantID,
		&subscription.UUID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.BillingCycle,
		&subscription.PayPalSubscriptionID,
		&subscription.PayPalPlanID,
		&subscription.Amount,
		&subscription.Currency,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.NextBillingDate,
		&cancelledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active subscription found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}

	if cancelledAt.Valid {
		subscription.CancelledAt = &cancelledAt.Time
	}

	return subscription, nil
}
