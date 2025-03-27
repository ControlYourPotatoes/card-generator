package database

import (
    "fmt"
    "os"
    "log"
    
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
    // For Supabase, we need to prepend "db." to the host
    UseDBPrefix bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    // Just load from the known location - third one up from the database package
    envPath := "../../../.env"
    
    if _, err := os.Stat(envPath); err == nil {
        log.Printf("Loading environment from %s", envPath)
        if err := godotenv.Load(envPath); err != nil {
            log.Printf("Warning: error loading %s: %v", envPath, err)
        }
    } else {
        log.Printf("No .env file found at %s, using environment variables", envPath)
    }
    
    // Get database configuration
    dbHost := getEnvOrDefault("DB_HOST", "localhost")
    dbPort := getEnvOrDefault("DB_PORT", "5432")
    dbUser := getEnvOrDefault("DB_USER", "postgres")
    dbName := getEnvOrDefault("DB_NAME", "postgres")
    dbSSL := getEnvOrDefault("DB_SSLMODE", "require")
    useDbPrefix := getEnvOrDefault("USE_DB_PREFIX", "true")
    
    if log.Flags() != 0 { // Only log in verbose mode
        log.Printf("Database config: host=%s, port=%s, user=%s, database=%s, sslmode=%s, useDbPrefix=%s",
            dbHost, dbPort, dbUser, dbName, dbSSL, useDbPrefix)
    }
    
    return &Config{
        Host:        dbHost,
        Port:        dbPort,
        User:        dbUser,
        Password:    getEnvOrDefault("DB_PASSWORD", "postgres"),
        Database:    dbName,
        SSLMode:     dbSSL,
        UseDBPrefix: useDbPrefix == "true",
    }, nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists && value != "" {
        return value
    }
    return defaultValue
}

// ConnectionString returns a PostgreSQL connection string in pgx format
func (c *Config) ConnectionString() string {
    host := c.Host
    if c.UseDBPrefix {
        // For Supabase, prepend "db." to the host
        host = "db." + host
    }
    
    return fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s?sslmode=%s",
        c.User, c.Password, host, c.Port, c.Database, c.SSLMode,
    )
}