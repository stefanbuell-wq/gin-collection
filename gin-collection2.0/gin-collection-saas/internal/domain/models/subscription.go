package models

import "time"

// Subscription represents a tenant's subscription
type Subscription struct {
	ID                   int64              `json:"id"`
	TenantID             int64              `json:"tenant_id"`
	UUID                 string             `json:"uuid"`
	PlanID               string             `json:"plan_id"` // free, basic_monthly, pro_yearly, etc.
	Status               SubscriptionStatus `json:"status"`
	BillingCycle         BillingCycle       `json:"billing_cycle"`
	CurrentPeriodStart   *time.Time         `json:"current_period_start,omitempty"`
	CurrentPeriodEnd     *time.Time         `json:"current_period_end,omitempty"`
	NextBillingDate      *time.Time         `json:"next_billing_date,omitempty"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	PayPalCustomerID     *string            `json:"paypal_customer_id,omitempty"`
	PayPalSubscriptionID *string            `json:"paypal_subscription_id,omitempty"`
	PayPalPlanID         *string            `json:"paypal_plan_id,omitempty"`
	Amount               float64            `json:"amount"`
	Currency             string             `json:"currency"`
	TrialEndsAt          *time.Time         `json:"trial_ends_at,omitempty"`
	CancelledAt          *time.Time         `json:"cancelled_at,omitempty"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "pending"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
)

// BillingCycle represents how often the subscription is billed
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

// SubscriptionPlan represents a subscription plan
type SubscriptionPlan struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Tier          SubscriptionTier `json:"tier"`
	BillingCycle  BillingCycle     `json:"billing_cycle"`
	PriceMonthly  float64          `json:"price_monthly"`
	PriceYearly   float64          `json:"price_yearly"`
	Currency      string           `json:"currency"`
	PayPalPlanID  string           `json:"paypal_plan_id"`
	Features      []string         `json:"features"`
	Limits        PlanLimits       `json:"limits"`
}

// Plan represents a subscription plan (alias for backward compatibility)
type Plan = SubscriptionPlan

// AvailablePlans defines all available subscription plans
var AvailablePlans = []SubscriptionPlan{
	{
		ID:           "PLAN_FREE",
		Name:         "Free",
		Tier:         TierFree,
		PriceMonthly: 0,
		PriceYearly:  0,
		Currency:     "EUR",
		PayPalPlanID: "",
		Features: []string{
			"Up to 10 gins",
			"1 photo per gin",
			"Basic search",
			"Mobile app (PWA)",
		},
		Limits: PlanLimitsMap[TierFree],
	},
	{
		ID:           "PLAN_BASIC_MONTHLY",
		Name:         "Basic",
		Tier:         TierBasic,
		PriceMonthly: 2.99,
		PriceYearly:  29.99,
		Currency:     "EUR",
		PayPalPlanID: "P-BASIC-MONTHLY",
		Features: []string{
			"Up to 50 gins",
			"3 photos per gin",
			"Tasting notes",
			"Advanced search & filters",
		},
		Limits: PlanLimitsMap[TierBasic],
	},
	{
		ID:           "PLAN_BASIC_YEARLY",
		Name:         "Basic (Yearly)",
		Tier:         TierBasic,
		PriceMonthly: 2.99,
		PriceYearly:  29.99,
		Currency:     "EUR",
		PayPalPlanID: "P-BASIC-YEARLY",
		Features: []string{
			"Up to 50 gins",
			"3 photos per gin",
			"Tasting notes",
			"Advanced search & filters",
			"Save 17% vs monthly",
		},
		Limits: PlanLimitsMap[TierBasic],
	},
	{
		ID:           "PLAN_PRO_MONTHLY",
		Name:         "Pro",
		Tier:         TierPro,
		PriceMonthly: 5.99,
		PriceYearly:  59.99,
		Currency:     "EUR",
		PayPalPlanID: "P-PRO-MONTHLY",
		Features: []string{
			"Unlimited gins",
			"10 photos per gin",
			"Botanicals tracking",
			"Cocktail recipes",
			"AI suggestions",
			"Export/Import",
		},
		Limits: PlanLimitsMap[TierPro],
	},
	{
		ID:           "PLAN_PRO_YEARLY",
		Name:         "Pro (Yearly)",
		Tier:         TierPro,
		PriceMonthly: 5.99,
		PriceYearly:  59.99,
		Currency:     "EUR",
		PayPalPlanID: "P-PRO-YEARLY",
		Features: []string{
			"Unlimited gins",
			"10 photos per gin",
			"Botanicals tracking",
			"Cocktail recipes",
			"AI suggestions",
			"Export/Import",
			"Save 17% vs monthly",
		},
		Limits: PlanLimitsMap[TierPro],
	},
	{
		ID:           "PLAN_ENTERPRISE",
		Name:         "Enterprise",
		Tier:         TierEnterprise,
		PriceMonthly: 0,
		PriceYearly:  0,
		Currency:     "EUR",
		PayPalPlanID: "P-ENTERPRISE",
		Features: []string{
			"Everything in Pro",
			"Multi-user support",
			"Separate database",
			"API access",
			"Custom branding",
			"SLA guarantee",
			"Priority support",
		},
		Limits: PlanLimitsMap[TierEnterprise],
	},
}

// GetPlanByID retrieves a plan by its ID
func GetPlanByID(planID string) *SubscriptionPlan {
	for i := range AvailablePlans {
		if AvailablePlans[i].ID == planID {
			return &AvailablePlans[i]
		}
	}
	return nil
}

// UsageMetric represents usage tracking for a tenant
type UsageMetric struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	MetricName  string    `json:"metric_name"` // gin_count, api_calls, storage_mb
	CurrentValue int      `json:"current_value"`
	LimitValue  *int      `json:"limit_value,omitempty"` // nil = unlimited
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpgradeRequest represents a subscription upgrade request
type UpgradeRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

// UpgradeResponse represents a subscription upgrade response
type UpgradeResponse struct {
	PayPalApprovalURL string        `json:"paypal_approval_url,omitempty"`
	Subscription      *Subscription `json:"subscription"`
	Message           string        `json:"message"`
}
