package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PlatformAdminRepository handles platform admin database operations
type PlatformAdminRepository struct {
	db *sql.DB
}

// NewPlatformAdminRepository creates a new platform admin repository
func NewPlatformAdminRepository(db *sql.DB) *PlatformAdminRepository {
	return &PlatformAdminRepository{db: db}
}

// GetByEmail retrieves a platform admin by email
func (r *PlatformAdminRepository) GetByEmail(ctx context.Context, email string) (*models.PlatformAdmin, error) {
	query := `
		SELECT id, email, password_hash, name, is_active, last_login_at, created_at, updated_at
		FROM platform_admins
		WHERE email = ?
	`

	admin := &models.PlatformAdmin{}
	var name sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&name,
		&admin.IsActive,
		&lastLoginAt,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if name.Valid {
		admin.Name = &name.String
	}
	if lastLoginAt.Valid {
		admin.LastLoginAt = &lastLoginAt.Time
	}

	return admin, nil
}

// GetByID retrieves a platform admin by ID
func (r *PlatformAdminRepository) GetByID(ctx context.Context, id int64) (*models.PlatformAdmin, error) {
	query := `
		SELECT id, email, password_hash, name, is_active, last_login_at, created_at, updated_at
		FROM platform_admins
		WHERE id = ?
	`

	admin := &models.PlatformAdmin{}
	var name sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&name,
		&admin.IsActive,
		&lastLoginAt,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if name.Valid {
		admin.Name = &name.String
	}
	if lastLoginAt.Valid {
		admin.LastLoginAt = &lastLoginAt.Time
	}

	return admin, nil
}

// UpdateLastLogin updates the last login timestamp
func (r *PlatformAdminRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := `UPDATE platform_admins SET last_login_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// UpdatePassword updates the admin password
func (r *PlatformAdminRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `UPDATE platform_admins SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), id)
	return err
}

// List returns all platform admins
func (r *PlatformAdminRepository) List(ctx context.Context) ([]*models.PlatformAdmin, error) {
	query := `
		SELECT id, email, password_hash, name, is_active, last_login_at, created_at, updated_at
		FROM platform_admins
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*models.PlatformAdmin
	for rows.Next() {
		admin := &models.PlatformAdmin{}
		var name sql.NullString
		var lastLoginAt sql.NullTime

		if err := rows.Scan(
			&admin.ID,
			&admin.Email,
			&admin.PasswordHash,
			&name,
			&admin.IsActive,
			&lastLoginAt,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if name.Valid {
			admin.Name = &name.String
		}
		if lastLoginAt.Valid {
			admin.LastLoginAt = &lastLoginAt.Time
		}

		admins = append(admins, admin)
	}

	return admins, rows.Err()
}

// GetPlatformStats retrieves overall platform statistics
func (r *PlatformAdminRepository) GetPlatformStats(ctx context.Context) (*models.PlatformStats, error) {
	stats := &models.PlatformStats{
		TenantsByTier: make(map[string]int64),
	}

	// Total tenants by status
	statusQuery := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'suspended' THEN 1 ELSE 0 END) as suspended,
			SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) as cancelled
		FROM tenants
	`
	err := r.db.QueryRowContext(ctx, statusQuery).Scan(
		&stats.TotalTenants,
		&stats.ActiveTenants,
		&stats.SuspendedTenants,
		&stats.CancelledTenants,
	)
	if err != nil {
		return nil, err
	}

	// Total users
	userQuery := `SELECT COUNT(*) FROM users`
	r.db.QueryRowContext(ctx, userQuery).Scan(&stats.TotalUsers)

	// Total gins
	ginQuery := `SELECT COUNT(*) FROM gins`
	r.db.QueryRowContext(ctx, ginQuery).Scan(&stats.TotalGins)

	// Tenants by tier
	tierQuery := `SELECT tier, COUNT(*) FROM tenants GROUP BY tier`
	rows, err := r.db.QueryContext(ctx, tierQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tier string
			var count int64
			if err := rows.Scan(&tier, &count); err == nil {
				stats.TenantsByTier[tier] = count
			}
		}
	}

	// New tenants last 7 days
	recent7Query := `SELECT COUNT(*) FROM tenants WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)`
	r.db.QueryRowContext(ctx, recent7Query).Scan(&stats.NewTenantsLast7d)

	// New tenants last 30 days
	recent30Query := `SELECT COUNT(*) FROM tenants WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`
	r.db.QueryRowContext(ctx, recent30Query).Scan(&stats.NewTenantsLast30d)

	return stats, nil
}

// GetAllTenants retrieves all tenants with stats
func (r *PlatformAdminRepository) GetAllTenants(ctx context.Context, limit, offset int) ([]*models.TenantWithStats, int64, error) {
	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM tenants`
	r.db.QueryRowContext(ctx, countQuery).Scan(&total)

	// Get tenants with user and gin counts
	query := `
		SELECT
			t.id, t.uuid, t.name, t.subdomain, t.tier, t.is_enterprise,
			t.status, t.settings, t.created_at, t.updated_at,
			(SELECT COUNT(*) FROM users WHERE tenant_id = t.id) as user_count,
			(SELECT COUNT(*) FROM gins WHERE tenant_id = t.id) as gin_count
		FROM tenants t
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*models.TenantWithStats
	for rows.Next() {
		tenant := &models.Tenant{}
		tws := &models.TenantWithStats{Tenant: tenant}
		var settings sql.NullString

		if err := rows.Scan(
			&tenant.ID,
			&tenant.UUID,
			&tenant.Name,
			&tenant.Subdomain,
			&tenant.Tier,
			&tenant.IsEnterprise,
			&tenant.Status,
			&settings,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
			&tws.UserCount,
			&tws.GinCount,
		); err != nil {
			return nil, 0, err
		}

		if settings.Valid {
			tenant.Settings = []byte(settings.String)
		}

		results = append(results, tws)
	}

	return results, total, rows.Err()
}

// GetAllUsers retrieves all users across all tenants
func (r *PlatformAdminRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, int64, error) {
	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	r.db.QueryRowContext(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, tenant_id, uuid, email, first_name, last_name, role,
			   is_active, email_verified_at, last_login_at, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var firstName, lastName sql.NullString
		var emailVerifiedAt, lastLoginAt sql.NullTime

		if err := rows.Scan(
			&user.ID,
			&user.TenantID,
			&user.UUID,
			&user.Email,
			&firstName,
			&lastName,
			&user.Role,
			&user.IsActive,
			&emailVerifiedAt,
			&lastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if firstName.Valid {
			user.FirstName = &firstName.String
		}
		if lastName.Valid {
			user.LastName = &lastName.String
		}
		if emailVerifiedAt.Valid {
			user.EmailVerifiedAt = &emailVerifiedAt.Time
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, user)
	}

	return users, total, rows.Err()
}

// UpdateTenantStatus updates a tenant's status
func (r *PlatformAdminRepository) UpdateTenantStatus(ctx context.Context, tenantID int64, status models.TenantStatus) error {
	query := `UPDATE tenants SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), tenantID)
	return err
}

// UpdateTenantTier updates a tenant's subscription tier
func (r *PlatformAdminRepository) UpdateTenantTier(ctx context.Context, tenantID int64, tier models.SubscriptionTier) error {
	query := `UPDATE tenants SET tier = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, tier, time.Now(), tenantID)
	return err
}
