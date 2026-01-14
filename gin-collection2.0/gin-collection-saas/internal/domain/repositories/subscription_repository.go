package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// SubscriptionRepository defines subscription data access
type SubscriptionRepository interface {
	// Create creates a new subscription
	Create(ctx context.Context, subscription *models.Subscription) error

	// GetByID retrieves a subscription by ID
	GetByID(ctx context.Context, id int64) (*models.Subscription, error)

	// GetByTenantID retrieves the current subscription for a tenant
	GetByTenantID(ctx context.Context, tenantID int64) (*models.Subscription, error)

	// GetByPayPalSubscriptionID retrieves a subscription by PayPal subscription ID
	GetByPayPalSubscriptionID(ctx context.Context, paypalSubscriptionID string) (*models.Subscription, error)

	// Update updates a subscription
	Update(ctx context.Context, subscription *models.Subscription) error

	// UpdateStatus updates subscription status
	UpdateStatus(ctx context.Context, id int64, status models.SubscriptionStatus) error

	// List retrieves all subscriptions for a tenant
	List(ctx context.Context, tenantID int64) ([]*models.Subscription, error)

	// GetActiveSubscription retrieves the active subscription for a tenant
	GetActiveSubscription(ctx context.Context, tenantID int64) (*models.Subscription, error)
}
