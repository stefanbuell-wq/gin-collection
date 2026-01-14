package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// AuditLogRepository implements the audit log repository interface
type AuditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, entity_id,
		                        changes, ip_address, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		log.TenantID,
		log.UserID,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.Changes,
		log.IPAddress,
		log.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.ID = id
	return nil
}

// List retrieves audit logs for a tenant with pagination
func (r *AuditLogRepository) List(ctx context.Context, tenantID int64, limit, offset int) ([]*models.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity_type, entity_id,
		       changes, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE tenant_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Changes,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ListByUser retrieves audit logs for a specific user
func (r *AuditLogRepository) ListByUser(ctx context.Context, tenantID, userID int64, limit, offset int) ([]*models.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity_type, entity_id,
		       changes, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE tenant_id = ? AND user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs by user: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Changes,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ListByEntity retrieves audit logs for a specific entity
func (r *AuditLogRepository) ListByEntity(ctx context.Context, tenantID int64, entityType string, entityID int64, limit, offset int) ([]*models.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity_type, entity_id,
		       changes, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE tenant_id = ? AND entity_type = ? AND entity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, entityType, entityID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs by entity: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Changes,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// Count counts audit logs for a tenant
func (r *AuditLogRepository) Count(ctx context.Context, tenantID int64) (int, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE tenant_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return count, nil
}

// DeleteOlderThan deletes audit logs older than a certain date (for cleanup)
func (r *AuditLogRepository) DeleteOlderThan(ctx context.Context, tenantID int64, days int) error {
	query := `DELETE FROM audit_logs WHERE tenant_id = ? AND created_at < DATE_SUB(NOW(), INTERVAL ? DAY)`

	_, err := r.db.ExecContext(ctx, query, tenantID, days)
	if err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	return nil
}
