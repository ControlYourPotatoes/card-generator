package database

import (
    "fmt"
    "os"
    
    "github.com/joho/godotenv"
)

// Config holds database configuration
type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
    SSLMode  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    // Load .env file if it exists
    _ = godotenv.Load()
    
    return &Config{
        Host:     getEnvOrDefault("DB_HOST", "localhost"),
        Port:     getEnvOrDefault("DB_PORT", "5432"),
        User:     getEnvOrDefault("DB_USER", "postgres"),
        Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
        Database: getEnvOrDefault("DB_NAME", "postgres"),
        SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
    }, nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

// ConnectionString returns a PostgreSQL connection string in pgx format
func (c *Config) ConnectionString() string {
    return fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s?sslmode=%s",
        c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
    )
}