package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const migrationsTable = "schema_migrations"

func main() {
	// Parse flags
	var (
		dbHost     = flag.String("host", getEnv("DB_HOST", "localhost"), "Database host")
		dbPort     = flag.Int("port", getEnvInt("DB_PORT", 3306), "Database port")
		dbUser     = flag.String("user", getEnv("DB_USER", "gin_app"), "Database user")
		dbPassword = flag.String("password", getEnv("DB_PASSWORD", ""), "Database password")
		dbName     = flag.String("database", getEnv("DB_NAME", "gin_collection"), "Database name")
		migrationsPath = flag.String("path", "./internal/infrastructure/database/migrations", "Path to migrations directory")
		command    = flag.String("command", "up", "Command: up, down, status, create")
		steps      = flag.Int("steps", 0, "Number of migrations to run (0 = all)")
		name       = flag.String("name", "", "Migration name (for create command)")
	)

	flag.Parse()

	// Handle create command separately (doesn't need database)
	if *command == "create" {
		if *name == "" {
			log.Fatal("Migration name is required for create command")
		}
		if err := createMigration(*migrationsPath, *name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		return
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		*dbUser, *dbPassword, *dbHost, *dbPort, *dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Ensure migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Execute command
	switch *command {
	case "up":
		if err := runMigrations(db, *migrationsPath, "up", *steps); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down":
		if err := runMigrations(db, *migrationsPath, "down", *steps); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "status":
		if err := showStatus(db, *migrationsPath); err != nil {
			log.Fatalf("Failed to show status: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func ensureMigrationsTable(db *sql.DB) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, migrationsTable)

	_, err := db.Exec(query)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	applied := make(map[string]bool)

	query := fmt.Sprintf("SELECT version FROM %s", migrationsTable)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func getMigrationFiles(path, direction string) ([]string, error) {
	pattern := filepath.Join(path, fmt.Sprintf("*.%s.sql", direction))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func extractVersion(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return base
}

func runMigrations(db *sql.DB, migrationsPath, direction string, steps int) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	files, err := getMigrationFiles(migrationsPath, direction)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// For down migrations, reverse the order
	if direction == "down" {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	count := 0
	for _, file := range files {
		version := extractVersion(file)

		// Check if should run this migration
		shouldRun := false
		if direction == "up" {
			shouldRun = !applied[version]
		} else {
			shouldRun = applied[version]
		}

		if !shouldRun {
			continue
		}

		// Check steps limit
		if steps > 0 && count >= steps {
			break
		}

		// Read and execute migration
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		log.Printf("Running migration: %s", filepath.Base(file))

		// Execute migration (split by semicolons for multiple statements)
		statements := splitStatements(string(content))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w\nStatement: %s", file, err, stmt)
			}
		}

		// Update migrations table
		if direction == "up" {
			query := fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", migrationsTable)
			if _, err := db.Exec(query, version); err != nil {
				return fmt.Errorf("failed to record migration %s: %w", version, err)
			}
		} else {
			query := fmt.Sprintf("DELETE FROM %s WHERE version = ?", migrationsTable)
			if _, err := db.Exec(query, version); err != nil {
				return fmt.Errorf("failed to remove migration record %s: %w", version, err)
			}
		}

		count++
		log.Printf("Migration %s completed successfully", filepath.Base(file))
	}

	if count == 0 {
		log.Println("No migrations to run")
	} else {
		log.Printf("Ran %d migration(s) successfully", count)
	}

	return nil
}

func splitStatements(content string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := rune(0)

	for i, char := range content {
		current.WriteRune(char)

		// Track string literals
		if (char == '\'' || char == '"' || char == '`') && (i == 0 || content[i-1] != '\\') {
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				inString = false
			}
		}

		// Split on semicolons outside of strings
		if char == ';' && !inString {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" && stmt != ";" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}

	// Add remaining content
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func showStatus(db *sql.DB, migrationsPath string) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	files, err := getMigrationFiles(migrationsPath, "up")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("=================")

	for _, file := range files {
		version := extractVersion(file)
		status := "Pending"
		if applied[version] {
			status = "Applied"
		}
		fmt.Printf("[%s] %s\n", status, filepath.Base(file))
	}

	fmt.Printf("\nTotal: %d migrations, %d applied\n", len(files), len(applied))

	return nil
}

func createMigration(migrationsPath, name string) error {
	// Get next version number
	files, err := getMigrationFiles(migrationsPath, "up")
	if err != nil {
		return err
	}

	nextVersion := 1
	for _, file := range files {
		version := extractVersion(file)
		if v, err := strconv.Atoi(version); err == nil && v >= nextVersion {
			nextVersion = v + 1
		}
	}

	// Create migration files
	versionStr := fmt.Sprintf("%03d", nextVersion)
	baseName := fmt.Sprintf("%s_%s", versionStr, sanitizeName(name))

	upFile := filepath.Join(migrationsPath, baseName+".up.sql")
	downFile := filepath.Join(migrationsPath, baseName+".down.sql")

	header := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n", name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upFile, []byte(header+"-- Add your UP migration SQL here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}

	if err := os.WriteFile(downFile, []byte(header+"-- Add your DOWN migration SQL here\n"), 0644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)

	return nil
}

func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	var result strings.Builder
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
