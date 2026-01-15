package models

import "time"

// PlatformAdmin represents a super-admin user for platform management
type PlatformAdmin struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose
	Name         *string    `json:"name,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// PlatformAdminAuthResponse is returned after successful admin authentication
type PlatformAdminAuthResponse struct {
	Token string         `json:"token"`
	Admin *PlatformAdmin `json:"admin"`
}

// PlatformStats represents overall platform statistics
type PlatformStats struct {
	TotalTenants      int64            `json:"total_tenants"`
	ActiveTenants     int64            `json:"active_tenants"`
	SuspendedTenants  int64            `json:"suspended_tenants"`
	CancelledTenants  int64            `json:"cancelled_tenants"`
	TotalUsers        int64            `json:"total_users"`
	TotalGins         int64            `json:"total_gins"`
	TenantsByTier     map[string]int64 `json:"tenants_by_tier"`
	NewTenantsLast7d  int64            `json:"new_tenants_last_7d"`
	NewTenantsLast30d int64            `json:"new_tenants_last_30d"`
}

// TenantWithStats represents a tenant with additional statistics
type TenantWithStats struct {
	Tenant    *Tenant `json:"tenant"`
	UserCount int64   `json:"user_count"`
	GinCount  int64   `json:"gin_count"`
}

// SystemHealth represents system health status
type SystemHealth struct {
	Status    string            `json:"status"`
	Database  HealthStatus      `json:"database"`
	Redis     HealthStatus      `json:"redis"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime,omitempty"`
}

// HealthStatus represents health of a single component
type HealthStatus struct {
	Status  string `json:"status"` // "healthy", "unhealthy", "degraded"
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}
