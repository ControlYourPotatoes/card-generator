package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ControlYourPotatoes/card-generator/internal/storage/database/migration"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Manager handles database connections and initialization
type Manager struct {
	config      *Config
	db          *sql.DB
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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.db = db
	log.Println("Connected to database successfully")
	return nil
}

// Initialize initializes the database with schema and migrations
func (m *Manager) Initialize(migrationsDir string) error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	// Run migrations
	runner := migration.NewRunner(m.db, migrationsDir)
	if err := runner.Run(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.initialized = true
	log.Println("Database initialized successfully")
	return nil
}

// InitWithTestData initializes the database and seeds test data
func (m *Manager) InitWithTestData() error {
	// Create PostgreSQL store
	store, err := NewPostgresStore(m.config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	defer store.Close()

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

// GetDB returns the database connection
func (m *Manager) GetDB() *sql.DB {
	return m.db
}

// GetStore returns a new PostgreSQL store
func (m *Manager) GetStore() (*PostgresStore, error) {
	if !m.initialized {
		return nil, fmt.Errorf("database not initialized")
	}

	return NewPostgresStore(m.config.ConnectionString())
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}
	return m.db.Close()
}