package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Manager handles database connections and initialization
type Manager struct {
	config      *Config
	pool        *pgxpool.Pool
	initialized bool
}

// NewManager creates a new database manager
func NewManager(configPath string) (*Manager, error) {
	// Load environment variables from .env file if it exists
	if configPath != "" {
		err := godotenv.Load(configPath)
		if err != nil {
			log.Printf("Warning: .env file not found at %s, using environment variables", configPath)
		}
	}

	// Load configuration from environment variables
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	return &Manager{
		config:      config,
		initialized: false,
	}, nil
}

// Connect establishes a database connection
func (m *Manager) Connect() error {
	connStr := m.config.ConnectionString()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.pool = pool
	log.Println("Connected to database successfully")
	return nil
}

// Initialize initializes the database with schema and migrations
func (m *Manager) Initialize(migrationsDir string) error {
	if m.pool == nil {
		return fmt.Errorf("database not connected")
	}

	// Create PostgreSQL store
	store := &PostgresStore{pool: m.pool}

	// Initialize schema directly rather than using the runner
	// This avoids the sql.DB vs pgxpool.Pool type mismatch
	if err := store.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	m.initialized = true
	log.Println("Database initialized successfully")
	return nil
}

// GetDB returns the database connection
func (m *Manager) GetDB() *pgxpool.Pool {
	return m.pool
}

// GetStore returns a new PostgreSQL store
func (m *Manager) GetStore() (*PostgresStore, error) {
	if !m.initialized {
		return nil, fmt.Errorf("database not initialized")
	}

	return &PostgresStore{pool: m.pool}, nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.pool == nil {
		return nil
	}
	m.pool.Close()
	return nil
}

// InitWithTestData initializes the database and seeds test data
func (m *Manager) InitWithTestData() error {
	// Create PostgreSQL store directly with the pool
	store := &PostgresStore{pool: m.pool}

	// Initialize schema
	if err := store.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Create seeder
	seeder := NewSeeder(store)

	// Seed test data
	if err := seeder.SeedAll(); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	m.initialized = true
	log.Println("Database initialized with test data")
	return nil
}
