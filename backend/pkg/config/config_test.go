package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig_DefaultConfig(t *testing.T) {
	// Load config without any config file or env vars
	config, err := LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	// Verify default values
	if config.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Server.Port)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Server.Host)
	}

	if config.Database.Type != "memory" {
		t.Errorf("Expected default database type 'memory', got '%s'", config.Database.Type)
	}

	if config.Storage.Type != "file" {
		t.Errorf("Expected default storage type 'file', got '%s'", config.Storage.Type)
	}

	if config.Generator.ImageWidth != 750 {
		t.Errorf("Expected default image width 750, got %d", config.Generator.ImageWidth)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", config.Logging.Level)
	}
}

func TestLoadConfig_WithYAMLFile(t *testing.T) {
	// Create a temporary config file in the proper directory structure
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.test.yaml")

	configContent := `
server:
  port: 9090
  host: "0.0.0.0"
  environment: "test"

database:
  type: "postgres"
  name: "test_db"
  user: "test_user"

storage:
  type: "memory"
  base_path: "/tmp/test"

generator:
  image_width: 800
  image_height: 1100
  parallel_jobs: 8

logging:
  level: "debug"
  format: "text"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Set CONFIG_DIR environment variable to point to our temp config directory
	os.Setenv("CONFIG_DIR", configDir)
	defer os.Unsetenv("CONFIG_DIR")

	config, err := LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load config from YAML: %v", err)
	}

	// Verify YAML values override defaults
	if config.Server.Port != 9090 {
		t.Errorf("Expected port 9090 from YAML, got %d", config.Server.Port)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host '0.0.0.0' from YAML, got '%s'", config.Server.Host)
	}

	if config.Database.Type != "postgres" {
		t.Errorf("Expected database type 'postgres' from YAML, got '%s'", config.Database.Type)
	}

	if config.Generator.ImageWidth != 800 {
		t.Errorf("Expected image width 800 from YAML, got %d", config.Generator.ImageWidth)
	}

	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug' from YAML, got '%s'", config.Logging.Level)
	}
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables (using correct prefixes from the implementation)
	testEnvVars := map[string]string{
		"SERVER_PORT":  "7777",
		"SERVER_HOST":  "example.com",
		"DB_TYPE":      "postgres",
		"DB_NAME":      "env_test_db",
		"STORAGE_TYPE": "s3",
		"LOG_LEVEL":    "error",
		"LOG_FORMAT":   "text",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up environment variables after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	config, err := LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load config with environment variables: %v", err)
	}

	// Verify environment variables override defaults
	if config.Server.Port != 7777 {
		t.Errorf("Expected port 7777 from env var, got %d", config.Server.Port)
	}

	if config.Server.Host != "example.com" {
		t.Errorf("Expected host 'example.com' from env var, got '%s'", config.Server.Host)
	}

	if config.Database.Type != "postgres" {
		t.Errorf("Expected database type 'postgres' from env var, got '%s'", config.Database.Type)
	}

	if config.Logging.Level != "error" {
		t.Errorf("Expected log level 'error' from env var, got '%s'", config.Logging.Level)
	}
}

func TestLoadConfig_InvalidYAMLFile(t *testing.T) {
	// Create a temporary invalid config file
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.test.yaml")

	invalidContent := `
invalid yaml content
  - missing proper structure
    key without value:
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	// Set CONFIG_DIR environment variable to point to our temp config directory
	os.Setenv("CONFIG_DIR", configDir)
	defer os.Unsetenv("CONFIG_DIR")

	_, err = LoadConfig("test")
	if err == nil {
		t.Fatal("Expected error when loading invalid YAML file")
	}

	if !strings.Contains(err.Error(), "failed to load config from YAML") {
		t.Errorf("Expected YAML loading error, got: %v", err)
	}
}

func TestValidateConfig_InvalidServerPort(t *testing.T) {
	config := getDefaultConfig()
	config.Server.Port = -1

	err := validateConfig(config)
	if err == nil {
		t.Fatal("Expected validation error for invalid server port")
	}
}

func TestValidateConfig_InvalidLogLevel(t *testing.T) {
	config := getDefaultConfig()
	config.Logging.Level = "invalid-level"

	err := validateConfig(config)
	if err == nil {
		t.Fatal("Expected validation error for invalid log level")
	}
}

func TestValidateConfig_InvalidDatabaseType(t *testing.T) {
	config := getDefaultConfig()
	config.Database.Type = "invalid-db"

	err := validateConfig(config)
	if err == nil {
		t.Fatal("Expected validation error for invalid database type")
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	config := getDefaultConfig()

	err := validateConfig(config)
	if err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
}

func TestDatabaseConfig_GetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "postgres connection string",
			config: DatabaseConfig{
				Type:     "postgres",
				Host:     "localhost",
				Port:     5432,
				Name:     "testdb",
				User:     "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "postgres connection string without password",
			config: DatabaseConfig{
				Type:    "postgres",
				Host:    "localhost",
				Port:    5432,
				Name:    "testdb",
				User:    "testuser",
				SSLMode: "require",
			},
			expected: "host=localhost port=5432 user=testuser password= dbname=testdb sslmode=require",
		},
		{
			name: "sqlite connection string",
			config: DatabaseConfig{
				Type: "sqlite",
				Name: "test.db",
			},
			expected: "file:test.db?cache=shared&mode=rwc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetConnectionString()
			if result != tt.expected {
				t.Errorf("Expected connection string '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := getEnv("TEST_ENV_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid integer environment variable
	os.Setenv("TEST_INT_VAR", "42")
	defer os.Unsetenv("TEST_INT_VAR")

	result := getEnvInt("TEST_INT_VAR", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid integer environment variable
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	result = getEnvInt("TEST_INVALID_INT", 10)
	if result != 10 {
		t.Errorf("Expected default value 10, got %d", result)
	}

	// Test with non-existing environment variable
	result = getEnvInt("NON_EXISTENT_INT", 20)
	if result != 20 {
		t.Errorf("Expected default value 20, got %d", result)
	}
}

func TestConfigPriority(t *testing.T) {
	// Create temp config file in proper directory structure
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.test.yaml")

	configContent := `
server:
  port: 9000
  host: "yaml.host"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Set CONFIG_DIR environment variable and override value
	os.Setenv("CONFIG_DIR", configDir)
	os.Setenv("SERVER_PORT", "8888")
	defer func() {
		os.Unsetenv("CONFIG_DIR")
		os.Unsetenv("SERVER_PORT")
	}()

	config, err := LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Environment variable should override YAML
	if config.Server.Port != 8888 {
		t.Errorf("Expected port 8888 from env var (should override YAML), got %d", config.Server.Port)
	}

	// YAML should override default
	if config.Server.Host != "yaml.host" {
		t.Errorf("Expected host 'yaml.host' from YAML (should override default), got '%s'", config.Server.Host)
	}
}
