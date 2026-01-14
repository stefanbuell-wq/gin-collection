package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// TenantRepository implements the tenant repository interface
type TenantRepository struct {
	db *sql.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (uuid, name, subdomain, tier, is_enterprise, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		uuid.New().String(),
		tenant.Name,
		tenant.Subdomain,
		tenant.Tier,
		tenant.IsEnterprise,
		tenant.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tenant.ID = id
	return nil
}

// GetByID retrieves a tenant by ID
func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*models.Tenant, error) {
	query := `
		SELECT id, uuid, name, subdomain, tier, is_enterprise, db_connection_string,
		       status, created_at, updated_at
		FROM tenants
		WHERE id = ?
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.UUID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Tier,
		&tenant.IsEnterprise,
		&tenant.DBConnectionString,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return tenant, nil
}

// GetBySubdomain retrieves a tenant by subdomain
func (r *TenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error) {
	query := `
		SELECT id, uuid, name, subdomain, tier, is_enterprise, db_connection_string,
		       status, created_at, updated_at
		FROM tenants
		WHERE subdomain = ?
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, subdomain).Scan(
		&tenant.ID,
		&tenant.UUID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Tier,
		&tenant.IsEnterprise,
		&tenant.DBConnectionString,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by subdomain: %w", err)
	}

	return tenant, nil
}

// GetByUUID retrieves a tenant by UUID
func (r *TenantRepository) GetByUUID(ctx context.Context, uuid string) (*models.Tenant, error) {
	query := `
		SELECT id, uuid, name, subdomain, tier, is_enterprise, db_connection_string,
		       status, created_at, updated_at
		FROM tenants
		WHERE uuid = ?
	`

	tenant := &models.Tenant{}
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&tenant.ID,
		&tenant.UUID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Tier,
		&tenant.IsEnterprise,
		&tenant.DBConnectionString,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by UUID: %w", err)
	}

	return tenant, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
		UPDATE tenants
		SET name = ?, tier = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		tenant.Name,
		tenant.Tier,
		tenant.Status,
		tenant.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	return nil
}

// UpdateStatus updates tenant status
func (r *TenantRepository) UpdateStatus(ctx context.Context, id int64, status models.TenantStatus) error {
	query := `UPDATE tenants SET status = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update tenant status: %w", err)
	}

	return nil
}

// UpdateTier updates tenant subscription tier
func (r *TenantRepository) UpdateTier(ctx context.Context, id int64, tier models.SubscriptionTier) error {
	query := `UPDATE tenants SET tier = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, tier, id)
	if err != nil {
		return fmt.Errorf("failed to update tenant tier: %w", err)
	}

	return nil
}
