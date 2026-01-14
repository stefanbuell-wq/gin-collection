package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// UsageMetricsRepository implements usage metrics tracking
type UsageMetricsRepository struct {
	db *sql.DB
}

// NewUsageMetricsRepository creates a new usage metrics repository
func NewUsageMetricsRepository(db *sql.DB) *UsageMetricsRepository {
	return &UsageMetricsRepository{db: db}
}

// GetMetric retrieves current value for a metric
func (r *UsageMetricsRepository) GetMetric(ctx context.Context, tenantID int64, metricName string) (int, error) {
	// Get current period (month)
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	query := `
		SELECT current_value
		FROM usage_metrics
		WHERE tenant_id = ? AND metric_name = ? AND period_start = ?
	`

	var value int
	err := r.db.QueryRowContext(ctx, query, tenantID, metricName, periodStart).Scan(&value)

	if err == sql.ErrNoRows {
		// Metric doesn't exist yet, return 0
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get metric: %w", err)
	}

	return value, nil
}

// IncrementMetric increments a metric by delta
func (r *UsageMetricsRepository) IncrementMetric(ctx context.Context, tenantID int64, metricName string, delta int) error {
	// Get current period
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	// Upsert (insert or update)
	query := `
		INSERT INTO usage_metrics (tenant_id, metric_name, current_value, period_start, period_end, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			current_value = current_value + ?,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, metricName, delta, periodStart, periodEnd, delta)
	if err != nil {
		return fmt.Errorf("failed to increment metric: %w", err)
	}

	return nil
}

// DecrementMetric decrements a metric by delta
func (r *UsageMetricsRepository) DecrementMetric(ctx context.Context, tenantID int64, metricName string, delta int) error {
	// Get current period
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	query := `
		UPDATE usage_metrics
		SET current_value = GREATEST(0, current_value - ?),
		    updated_at = NOW()
		WHERE tenant_id = ? AND metric_name = ? AND period_start = ?
	`

	_, err := r.db.ExecContext(ctx, query, delta, tenantID, metricName, periodStart)
	if err != nil {
		return fmt.Errorf("failed to decrement metric: %w", err)
	}

	return nil
}

// SetMetric sets a metric to a specific value
func (r *UsageMetricsRepository) SetMetric(ctx context.Context, tenantID int64, metricName string, value int) error {
	// Get current period
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	query := `
		INSERT INTO usage_metrics (tenant_id, metric_name, current_value, period_start, period_end, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			current_value = ?,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, metricName, value, periodStart, periodEnd, value)
	if err != nil {
		return fmt.Errorf("failed to set metric: %w", err)
	}

	return nil
}
