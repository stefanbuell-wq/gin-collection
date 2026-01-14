package models

import "time"

// AuditLog represents an audit log entry for tracking user actions
type AuditLog struct {
	ID         int64      `json:"id"`
	TenantID   int64      `json:"tenant_id"`
	UserID     *int64     `json:"user_id,omitempty"` // Can be null for system actions
	Action     string     `json:"action"`            // e.g., "create_gin", "delete_user", "update_subscription"
	EntityType string     `json:"entity_type"`       // e.g., "gin", "user", "subscription"
	EntityID   *int64     `json:"entity_id,omitempty"`
	Changes    *string    `json:"changes,omitempty"` // JSON string of changes
	IPAddress  *string    `json:"ip_address,omitempty"`
	UserAgent  *string    `json:"user_agent,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// AuditAction represents common audit actions
type AuditAction string

const (
	// Gin actions
	AuditActionCreateGin AuditAction = "create_gin"
	AuditActionUpdateGin AuditAction = "update_gin"
	AuditActionDeleteGin AuditAction = "delete_gin"

	// User actions
	AuditActionCreateUser  AuditAction = "create_user"
	AuditActionUpdateUser  AuditAction = "update_user"
	AuditActionDeleteUser  AuditAction = "delete_user"
	AuditActionInviteUser  AuditAction = "invite_user"

	// Authentication actions
	AuditActionLogin         AuditAction = "login"
	AuditActionLogout        AuditAction = "logout"
	AuditActionFailedLogin   AuditAction = "failed_login"
	AuditActionGenerateAPIKey AuditAction = "generate_api_key"
	AuditActionRevokeAPIKey   AuditAction = "revoke_api_key"

	// Subscription actions
	AuditActionUpgradeSubscription   AuditAction = "upgrade_subscription"
	AuditActionCancelSubscription    AuditAction = "cancel_subscription"
	AuditActionActivateSubscription  AuditAction = "activate_subscription"

	// Photo actions
	AuditActionUploadPhoto AuditAction = "upload_photo"
	AuditActionDeletePhoto AuditAction = "delete_photo"

	// Export/Import actions
	AuditActionExportData AuditAction = "export_data"
	AuditActionImportData AuditAction = "import_data"
)

// EntityType represents the type of entity being audited
type EntityType string

const (
	EntityTypeGin          EntityType = "gin"
	EntityTypeUser         EntityType = "user"
	EntityTypeSubscription EntityType = "subscription"
	EntityTypePhoto        EntityType = "photo"
	EntityTypeTenant       EntityType = "tenant"
	EntityTypeBotanical    EntityType = "botanical"
	EntityTypeCocktail     EntityType = "cocktail"
)
