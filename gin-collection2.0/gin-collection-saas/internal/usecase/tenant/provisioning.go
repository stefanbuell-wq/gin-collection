package tenant

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// ProvisioningService handles Enterprise database provisioning
type ProvisioningService struct {
	db         *sql.DB
	tenantRepo repositories.TenantRepository
	dbHost     string
	dbUser     string
	dbPassword string
}

// NewProvisioningService creates a new provisioning service
func NewProvisioningService(
	db *sql.DB,
	tenantRepo repositories.TenantRepository,
	dbHost, dbUser, dbPassword string,
) *ProvisioningService {
	return &ProvisioningService{
		db:         db,
		tenantRepo: tenantRepo,
		dbHost:     dbHost,
		dbUser:     dbUser,
		dbPassword: dbPassword,
	}
}

// ProvisionEnterpriseDatabase creates a separate database for an Enterprise tenant
func (s *ProvisioningService) ProvisionEnterpriseDatabase(ctx context.Context, tenantID int64) error {
	logger.Info("Provisioning Enterprise database", "tenant_id", tenantID)

	// Get tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Verify tenant is Enterprise
	if tenant.Tier != "enterprise" {
		return fmt.Errorf("tenant is not Enterprise tier")
	}

	// Check if already provisioned
	if tenant.DBConnectionString != nil && *tenant.DBConnectionString != "" {
		logger.Warn("Database already provisioned for tenant", "tenant_id", tenantID)
		return fmt.Errorf("database already provisioned")
	}

	// Generate database name: gin_collection_tenant_{tenant_id}
	dbName := fmt.Sprintf("gin_collection_tenant_%d", tenantID)

	// Create database
	logger.Info("Creating database", "db_name", dbName)
	_, err = s.db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		logger.Error("Failed to create database", "error", err.Error())
		return fmt.Errorf("failed to create database: %w", err)
	}

	// TODO: Run migrations on the new database
	// This would involve:
	// 1. Connecting to the new database
	// 2. Running all schema migrations
	// 3. Seeding reference data (botanicals, cocktails)
	//
	// For now, we'll log a warning that migrations need to be run manually

	logger.Warn("Database created - migrations must be run manually", "db_name", dbName)

	// Generate connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		s.dbUser,
		s.dbPassword,
		s.dbHost,
		dbName,
	)

	// Update tenant with connection string
	tenant.DBConnectionString = &connectionString
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		logger.Error("Failed to update tenant with connection string", "error", err.Error())
		// Try to rollback database creation
		s.db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	logger.Info("Enterprise database provisioned successfully", "tenant_id", tenantID, "db_name", dbName)

	return nil
}

// MigrateTenantToEnterprise migrates a tenant's data from shared DB to dedicated DB
func (s *ProvisioningService) MigrateTenantToEnterprise(ctx context.Context, tenantID int64) error {
	logger.Info("Migrating tenant to Enterprise", "tenant_id", tenantID)

	// Get tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Verify tenant is Enterprise
	if tenant.Tier != "enterprise" {
		return fmt.Errorf("tenant is not Enterprise tier")
	}

	// Check if database is provisioned
	if tenant.DBConnectionString == nil || *tenant.DBConnectionString == "" {
		return fmt.Errorf("database not provisioned - run ProvisionEnterpriseDatabase first")
	}

	// TODO: Implement data migration
	// This would involve:
	// 1. Connecting to the dedicated database
	// 2. Copying all tenant data from shared DB
	// 3. Verifying data integrity
	// 4. Deleting data from shared DB (optional)
	//
	// This is a complex operation that should be done carefully
	// For now, we'll return an error indicating manual migration is required

	logger.Warn("Data migration not implemented - manual migration required", "tenant_id", tenantID)

	return fmt.Errorf("automatic data migration not implemented - please migrate manually")
}

// DecommissionEnterpriseDatabase removes the dedicated database (for downgrades/deletions)
func (s *ProvisioningService) DecommissionEnterpriseDatabase(ctx context.Context, tenantID int64) error {
	logger.Info("Decommissioning Enterprise database", "tenant_id", tenantID)

	// Get tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Check if database is provisioned
	if tenant.DBConnectionString == nil || *tenant.DBConnectionString == "" {
		logger.Warn("No database to decommission", "tenant_id", tenantID)
		return nil
	}

	// Generate database name
	dbName := fmt.Sprintf("gin_collection_tenant_%d", tenantID)

	// Drop database
	logger.Warn("Dropping database", "db_name", dbName)
	_, err = s.db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		logger.Error("Failed to drop database", "error", err.Error())
		return fmt.Errorf("failed to drop database: %w", err)
	}

	// Update tenant to remove connection string
	tenant.DBConnectionString = nil
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		logger.Error("Failed to update tenant", "error", err.Error())
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	logger.Info("Enterprise database decommissioned successfully", "tenant_id", tenantID)

	return nil
}

// HealthCheck verifies the dedicated database is accessible
func (s *ProvisioningService) HealthCheck(ctx context.Context, tenantID int64) (bool, error) {
	logger.Debug("Checking Enterprise database health", "tenant_id", tenantID)

	// Get tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Check if database is provisioned
	if tenant.DBConnectionString == nil || *tenant.DBConnectionString == "" {
		return false, fmt.Errorf("database not provisioned")
	}

	// Try to connect to the dedicated database
	db, err := sql.Open("mysql", *tenant.DBConnectionString)
	if err != nil {
		logger.Error("Failed to connect to dedicated database", "error", err.Error())
		return false, fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	// Ping the database
	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping dedicated database", "error", err.Error())
		return false, fmt.Errorf("failed to ping: %w", err)
	}

	logger.Debug("Enterprise database is healthy", "tenant_id", tenantID)

	return true, nil
}
