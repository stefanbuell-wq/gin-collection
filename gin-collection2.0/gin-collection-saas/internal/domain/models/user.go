package models

import "time"

// User represents a user within a tenant
type User struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	UUID            string    `json:"uuid"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"` // Never expose password hash in JSON
	FirstName       *string   `json:"first_name,omitempty"`
	LastName        *string   `json:"last_name,omitempty"`
	Role            UserRole  `json:"role"`
	APIKey          *string   `json:"api_key,omitempty"` // Enterprise only
	IsActive        bool      `json:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserRole represents the role of a user within a tenant
type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
	RoleViewer UserRole = "viewer"
)

// HasPermission checks if the user role has permission for an action
func (r UserRole) HasPermission(action string) bool {
	permissions := map[UserRole]map[string]bool{
		RoleOwner: {
			"create":         true,
			"read":           true,
			"update":         true,
			"delete":         true,
			"manage_users":   true,
			"manage_billing": true,
		},
		RoleAdmin: {
			"create":       true,
			"read":         true,
			"update":       true,
			"delete":       true,
			"manage_users": true,
		},
		RoleMember: {
			"create": true,
			"read":   true,
			"update": true,
			"delete": false,
		},
		RoleViewer: {
			"read": true,
		},
	}

	rolePerms, exists := permissions[r]
	if !exists {
		return false
	}

	return rolePerms[action]
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	TenantName string `json:"tenant_name" binding:"required,min=2,max=255"`
	Subdomain  string `json:"subdomain" binding:"required,min=3,max=63,alphanum"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=12"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string  `json:"token"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	User         *User   `json:"user"`
	Tenant       *Tenant `json:"tenant"`
}
