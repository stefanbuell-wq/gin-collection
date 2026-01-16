package models

import (
	"encoding/json"
	"time"
)

// Tenant represents a SaaS tenant (customer organization)
type Tenant struct {
	ID                 int64           `json:"id"`
	UUID               string          `json:"uuid"`
	Name               string          `json:"name"`
	Subdomain          string          `json:"subdomain"`
	Tier               SubscriptionTier `json:"tier"`
	IsEnterprise       bool            `json:"is_enterprise"`
	DBConnectionString *string         `json:"-"` // Hidden from JSON, only for Enterprise
	Status             TenantStatus    `json:"status"`
	Settings           json.RawMessage `json:"settings,omitempty"`
	Branding           *TenantBranding `json:"branding,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// TenantStatus represents the status of a tenant
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusCancelled TenantStatus = "cancelled"
)

// TenantBranding holds custom branding for Enterprise tenants
type TenantBranding struct {
	LogoURL      string `json:"logo_url,omitempty"`
	PrimaryColor string `json:"primary_color,omitempty"`
	CustomDomain string `json:"custom_domain,omitempty"`
}

// SubscriptionTier represents the subscription tier
type SubscriptionTier string

const (
	TierFree       SubscriptionTier = "free"
	TierBasic      SubscriptionTier = "basic"
	TierPro        SubscriptionTier = "pro"
	TierEnterprise SubscriptionTier = "enterprise"
)

// PlanLimits defines limits for each subscription tier
type PlanLimits struct {
	MaxGins          *int  // nil = unlimited
	MaxPhotosPerGin  int
	HasBotanicals    bool
	HasCocktails     bool
	HasAISuggestions bool
	HasExport        bool
	HasImport        bool
	HasMultiUser     bool
	HasAPIAccess     bool
	APIRateLimit     int  // requests per hour
	StorageLimitMB   *int // nil = unlimited
}

// PlanLimitsMap defines limits for each tier
var PlanLimitsMap = map[SubscriptionTier]PlanLimits{
	TierFree: {
		MaxGins:          intPtr(25),
		MaxPhotosPerGin:  3,
		HasBotanicals:    false,
		HasCocktails:     false,
		HasAISuggestions: false,
		HasExport:        false,
		HasImport:        false,
		HasMultiUser:     false,
		HasAPIAccess:     false,
		APIRateLimit:     100,
		StorageLimitMB:   intPtr(100),
	},
	TierBasic: {
		MaxGins:          intPtr(100),
		MaxPhotosPerGin:  10,
		HasBotanicals:    false,
		HasCocktails:     false,
		HasAISuggestions: false,
		HasExport:        true,
		HasImport:        false,
		HasMultiUser:     false,
		HasAPIAccess:     false,
		APIRateLimit:     500,
		StorageLimitMB:   intPtr(1000),
	},
	TierPro: {
		MaxGins:          intPtr(500),
		MaxPhotosPerGin:  25,
		HasBotanicals:    true,
		HasCocktails:     true,
		HasAISuggestions: true,
		HasExport:        true,
		HasImport:        true,
		HasMultiUser:     false,
		HasAPIAccess:     true,
		APIRateLimit:     5000,
		StorageLimitMB:   intPtr(5000),
	},
	TierEnterprise: {
		MaxGins:          nil, // unlimited
		MaxPhotosPerGin:  -1,  // unlimited (use -1 for unlimited photos)
		HasBotanicals:    true,
		HasCocktails:     true,
		HasAISuggestions: true,
		HasExport:        true,
		HasImport:        true,
		HasMultiUser:     true,
		HasAPIAccess:     true,
		APIRateLimit:     10000,
		StorageLimitMB:   nil, // unlimited
	},
}

// GetLimits returns the plan limits for this tenant's tier
func (t *Tenant) GetLimits() PlanLimits {
	return PlanLimitsMap[t.Tier]
}

// helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
