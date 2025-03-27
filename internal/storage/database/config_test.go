package database

import (
	"os"
	"path/filepath"
	"testing"

)

func TestLoadConfig(t *testing.T) {
	// Save original environment variables to restore later
	origEnv := map[string]string{
		"DB_HOST":       os.Getenv("DB_HOST"),
		"DB_PORT":       os.Getenv("DB_PORT"),
		"DB_USER":       os.Getenv("DB_USER"),
		"DB_PASSWORD":   os.Getenv("DB_PASSWORD"),
		"DB_NAME":       os.Getenv("DB_NAME"),
		"DB_SSLMODE":    os.Getenv("DB_SSLMODE"),
		"USE_DB_PREFIX": os.Getenv("USE_DB_PREFIX"),
	}

	// Restore environment variables after test
	defer func() {
		for k, v := range origEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Clear environment variables to ensure we're testing .env loading
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("USE_DB_PREFIX")

	t.Run("LoadConfigWithPath", func(t *testing.T) {
		// Create a temporary .env file
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		// Write test values to the .env file
		envContent := `
DB_HOST=testhost
DB_PORT=5555
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
DB_SSLMODE=disable
USE_DB_PREFIX=false
`
		if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
			t.Fatalf("Failed to create test .env file: %v", err)
		}

		// Get absolute path for logging
		absPath, _ := filepath.Abs(envPath)
		t.Logf("Test .env file location: %s", absPath)

		// Use LoadConfigWithPath function
		config, err := LoadConfigWithPath(envPath)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify values from .env file were loaded
		if config.Host != "testhost" {
			t.Errorf("Expected Host=testhost, got %s", config.Host)
		}
		if config.Port != "5555" {
			t.Errorf("Expected Port=5555, got %s", config.Port)
		}
		if config.User != "testuser" {
			t.Errorf("Expected User=testuser, got %s", config.User)
		}
		if config.Password != "testpass" {
			t.Errorf("Expected Password=testpass, got %s", config.Password)
		}
		if config.Database != "testdb" {
			t.Errorf("Expected Database=testdb, got %s", config.Database)
		}
		if config.SSLMode != "disable" {
			t.Errorf("Expected SSLMode=disable, got %s", config.SSLMode)
		}
		if config.UseDBPrefix != false {
			t.Errorf("Expected UseDBPrefix=false, got %t", config.UseDBPrefix)
		}
	})

	t.Run("LoadConfigWithPath_MissingFile", func(t *testing.T) {
		// Use non-existent file
		envPath := "/path/to/nonexistent/file.env"

		// Attempt to load config
		_, err := LoadConfigWithPath(envPath)

		// Should get an error
		if err == nil {
			t.Error("Expected error when env file doesn't exist, but got nil")
		}
	})

	t.Run("LoadConfigWithPath_MissingVars", func(t *testing.T) {
		// Create a temporary .env file with missing required variables
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		// Write incomplete env file
		envContent := `
DB_HOST=testhost
DB_PORT=5555
# Missing DB_USER
# Missing DB_PASSWORD
DB_NAME=testdb
`
		if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
			t.Fatalf("Failed to create test .env file: %v", err)
		}

		// Attempt to load config
		_, err := LoadConfigWithPath(envPath)

		// Should get an error about missing variables
		if err == nil {
			t.Error("Expected error when required variables are missing, but got nil")
		}
	})

	t.Run("LoadConfig_WithCurrentDir", func(t *testing.T) {
		// Create a temporary directory with .env file
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		// Write test values to the .env file
		envContent := `
DB_HOST=testhost
DB_PORT=5555
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
DB_SSLMODE=disable
USE_DB_PREFIX=false
`
		if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
			t.Fatalf("Failed to create test .env file: %v", err)
		}

		// Save current directory
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}

		// Change to temp directory and restore after test
		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(originalDir)

		// Now the .env file should be in the current directory
		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify values
		if config.Host != "testhost" {
			t.Errorf("Expected Host=testhost, got %s", config.Host)
		}
	})

	t.Run("LoadConfig_NoEnvFile", func(t *testing.T) {
		// Create empty temp directory with no .env file
		tempDir := t.TempDir()

		// Change to that directory
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}

		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		defer os.Chdir(originalDir)

		// Attempt to load config
		_, err = LoadConfig()

		// Should get an error
		if err == nil {
			t.Error("Expected error when no .env file exists, but got nil")
		}
	})

	t.Run("ConnectionString", func(t *testing.T) {
		// Test the connection string generation
		config := &Config{
			Host:        "testhost",
			Port:        "5555",
			User:        "testuser",
			Password:    "testpass",
			Database:    "testdb",
			SSLMode:     "disable",
			UseDBPrefix: true,
		}

		expected := "postgresql://testuser:testpass@db.testhost:5555/testdb?sslmode=disable"
		connStr := config.ConnectionString()
		if connStr != expected {
			t.Errorf("Expected connection string %s, got %s", expected, connStr)
		}

		// Test without DB prefix
		config.UseDBPrefix = false
		expected = "postgresql://testuser:testpass@testhost:5555/testdb?sslmode=disable"
		connStr = config.ConnectionString()
		if connStr != expected {
			t.Errorf("Expected connection string %s, got %s", expected, connStr)
		}
	})
}