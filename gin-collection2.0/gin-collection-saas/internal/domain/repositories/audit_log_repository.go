package repositories

import (
	"context"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *models.AuditLog) error

	// List retrieves audit logs for a tenant with pagination
	List(ctx context.Context, tenantID int64, limit, offset int) ([]*models.AuditLog, error)

	// ListByUser retrieves audit logs for a specific user
	ListByUser(ctx context.Context, tenantID, userID int64, limit, offset int) ([]*models.AuditLog, error)

	// ListByEntity retrieves audit logs for a specific entity
	ListByEntity(ctx context.Context, tenantID int64, entityType string, entityID int64, limit, offset int) ([]*models.AuditLog, error)

	// Count counts audit logs for a tenant
	Count(ctx context.Context, tenantID int64) (int, error)

	// DeleteOlderThan deletes audit logs older than a certain date (for cleanup)
	DeleteOlderThan(ctx context.Context, tenantID int64, days int) error
}
