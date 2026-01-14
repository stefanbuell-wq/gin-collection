package database

import (
	"database/sql"
	"fmt"
	"sync"
)

// TenantRouter routes database queries to the correct database
// - Shared database for Free/Basic/Pro tiers
// - Separate database for Enterprise tier
type TenantRouter struct {
	sharedDB      *sql.DB
	enterpriseDBs map[int64]*sql.DB // tenantID -> DB connection
	mutex         sync.RWMutex
}

// NewTenantRouter creates a new tenant router with a shared database
func NewTenantRouter(sharedDB *sql.DB) *TenantRouter {
	return &TenantRouter{
		sharedDB:      sharedDB,
		enterpriseDBs: make(map[int64]*sql.DB),
	}
}

// GetDB returns the appropriate database connection for a tenant
// If isEnterprise is false, returns the shared database
// If isEnterprise is true, returns the tenant-specific database (lazy-loaded)
func (tr *TenantRouter) GetDB(tenantID int64, isEnterprise bool, dsn string) (*sql.DB, error) {
	if !isEnterprise {
		return tr.sharedDB, nil
	}

	// Check if enterprise DB is already loaded
	tr.mutex.RLock()
	db, exists := tr.enterpriseDBs[tenantID]
	tr.mutex.RUnlock()

	if exists {
		return db, nil
	}

	// Lazy load enterprise database
	return tr.loadEnterpriseDB(tenantID, dsn)
}

// loadEnterpriseDB connects to an enterprise tenant's dedicated database
func (tr *TenantRouter) loadEnterpriseDB(tenantID int64, dsn string) (*sql.DB, error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// Double-check if another goroutine already loaded it
	if db, exists := tr.enterpriseDBs[tenantID]; exists {
		return db, nil
	}

	if dsn == "" {
		return nil, fmt.Errorf("enterprise tenant %d has no database connection string", tenantID)
	}

	// Create new connection for enterprise tenant
	db, err := NewMySQL(dsn, 10, 5) // Lower connection limits per enterprise tenant
	if err != nil {
		return nil, fmt.Errorf("failed to connect to enterprise database for tenant %d: %w", tenantID, err)
	}

	tr.enterpriseDBs[tenantID] = db
	return db, nil
}

// CloseEnterpriseDB closes an enterprise tenant's database connection
func (tr *TenantRouter) CloseEnterpriseDB(tenantID int64) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	db, exists := tr.enterpriseDBs[tenantID]
	if !exists {
		return nil // Already closed or never loaded
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close enterprise database for tenant %d: %w", tenantID, err)
	}

	delete(tr.enterpriseDBs, tenantID)
	return nil
}

// Close closes all database connections (shared + all enterprise)
func (tr *TenantRouter) Close() error {
	// Close all enterprise databases
	tr.mutex.Lock()
	for tenantID, db := range tr.enterpriseDBs {
		if err := db.Close(); err != nil {
			// Log error but continue closing others
			fmt.Printf("Error closing enterprise DB for tenant %d: %v\n", tenantID, err)
		}
	}
	tr.enterpriseDBs = make(map[int64]*sql.DB)
	tr.mutex.Unlock()

	// Close shared database
	if err := tr.sharedDB.Close(); err != nil {
		return fmt.Errorf("failed to close shared database: %w", err)
	}

	return nil
}

// HealthCheck checks health of shared database and all loaded enterprise databases
func (tr *TenantRouter) HealthCheck() map[string]error {
	errors := make(map[string]error)

	// Check shared database
	if err := HealthCheck(tr.sharedDB); err != nil {
		errors["shared"] = err
	}

	// Check all enterprise databases
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	for tenantID, db := range tr.enterpriseDBs {
		if err := HealthCheck(db); err != nil {
			errors[fmt.Sprintf("enterprise_tenant_%d", tenantID)] = err
		}
	}

	return errors
}
