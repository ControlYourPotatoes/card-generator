package database

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestDatabaseConnection(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Logf("Warning: .env file not found, using environment variables")
	}

	// Load database configuration
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create store with connection string
	store, err := NewPostgresStore(config.ConnectionString())
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	// Test basic ping to verify connection
	t.Log("Testing database connection...")
	start := time.Now()
	
	// Ping database (already done in NewPostgresStore, but we'll do it again for clarity)
	err = store.db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
	
	t.Logf("Database connection successful! (took %v)", time.Since(start))

	// Test simple query to verify permissions
	t.Log("Testing simple query...")
	var version string
	err = store.db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to execute simple query: %v", err)
	}
	
	t.Logf("Query successful! PostgreSQL version: %s", version)
}