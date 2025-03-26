package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Runner manages database migrations
type Runner struct {
	db            *sql.DB
	migrationsDir string
}

// NewRunner creates a new migration runner
func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// EnsureMigrationsTable creates the migrations table if it doesn't exist
func (r *Runner) EnsureMigrationsTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// GetAppliedMigrations returns a map of already applied migrations
func (r *Runner) GetAppliedMigrations() (map[int]bool, error) {
	rows, err := r.db.Query("SELECT version FROM migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

// FindMigrationFiles finds all SQL migration files in the migrations directory
func (r *Runner) FindMigrationFiles() (map[int]string, error) {
	// Pattern to match migration files like: 001_initial_schema.sql
	pattern := regexp.MustCompile(`^(\d+)_(.+)\.sql$`)
	
	files, err := ioutil.ReadDir(r.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	migrations := make(map[int]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		matches := pattern.FindStringSubmatch(file.Name())
		if matches == nil {
			continue
		}
		
		version, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version: %s", matches[1])
		}
		
		migrations[version] = filepath.Join(r.migrationsDir, file.Name())
	}
	
	return migrations, nil
}

// Run executes all pending migrations
func (r *Runner) Run() error {
	if err := r.EnsureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	applied, err := r.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	
	migrations, err := r.FindMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}
	
	// Get sorted versions
	var versions []int
	for v := range migrations {
		versions = append(versions, v)
	}
	sort.Ints(versions)
	
	// Run each migration in order
	for _, version := range versions {
		if applied[version] {
			log.Printf("Migration %d already applied, skipping", version)
			continue
		}
		
		filePath := migrations[version]
		fileName := filepath.Base(filePath)
		
		log.Printf("Applying migration %d: %s", version, fileName)
		
		// Read migration content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}
		
		// Start a transaction
		tx, err := r.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		
		// Execute migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", filePath, err)
		}
		
		// Record migration
		_, err = tx.Exec(
			"INSERT INTO migrations (version, name) VALUES ($1, $2)",
			version, fileName,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", filePath, err)
		}
		
		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		
		log.Printf("Successfully applied migration %d", version)
	}
	
	return nil
}

// RollbackLast reverts the last applied migration
func (r *Runner) RollbackLast() error {
	// Find the last applied migration
	var (
		version int
		name    string
	)
	
	err := r.db.QueryRow(`
		SELECT version, name FROM migrations
		ORDER BY version DESC
		LIMIT 1
	`).Scan(&version, &name)
	
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}
	
	// Look for rollback file
	rollbackPath := filepath.Join(r.migrationsDir, fmt.Sprintf("%03d_%s.down.sql", version, strings.TrimSuffix(name, ".sql")))
	
	// Check if rollback file exists
	if _, err := ioutil.ReadFile(rollbackPath); err != nil {
		return fmt.Errorf("rollback file not found for migration %d: %w", version, err)
	}
	
	// Read rollback SQL
	content, err := ioutil.ReadFile(rollbackPath)
	if err != nil {
		return fmt.Errorf("failed to read rollback file %s: %w", rollbackPath, err)
	}
	
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Execute rollback
	_, err = tx.Exec(string(content))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback %s: %w", rollbackPath, err)
	}
	
	// Remove migration record
	_, err = tx.Exec("DELETE FROM migrations WHERE version = $1", version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %w", err)
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Printf("Successfully rolled back migration %d", version)
	return nil
}