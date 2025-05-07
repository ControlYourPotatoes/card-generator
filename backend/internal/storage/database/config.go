package database

import (
    "fmt"
    "os"
    "log"
    "path/filepath"
    
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
    cwd, err := os.Getwd()
    if err != nil {
        log.Printf("Warning: could not get current working directory: %v", err)
    } else {
        log.Printf("Current working directory: %s", cwd)
    }
    
    // Try to find the .env file in several possible locations
    // First, get the absolute path of the current working directory
    workingDir, err := os.Getwd()
    if err != nil {
        log.Printf("Warning: could not get current working directory: %v", err)
        workingDir = "."
    }
    
    // Create paths to search
    envPaths := []string{
        filepath.Join(workingDir, ".env"),                       // Current directory
        filepath.Join(filepath.Dir(workingDir), ".env"),         // One level up
        filepath.Join(filepath.Dir(filepath.Dir(workingDir)), ".env"), // Two levels up
        filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(workingDir))), ".env"), // Three levels up
    }
    
    envFileFound := false
    for _, envPath := range envPaths {
        
        log.Printf("Looking for .env at: %s", envPath)
        if _, err := os.Stat(envPath); err == nil {
            log.Printf("Loading environment from %s", envPath)
            if err := godotenv.Load(envPath); err != nil {
                log.Printf("Warning: error loading %s: %v", envPath, err)
            } else {
                envFileFound = true
                break
            }
        }
    }
    
    if !envFileFound {
        return nil, fmt.Errorf("no .env file found in any of the checked locations")
    }
    
    // Get database configuration
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbSSL := os.Getenv("DB_SSLMODE")
    useDbPrefix := os.Getenv("USE_DB_PREFIX")
    
    // Check that all required values are present
    missingVars := []string{}
    if dbHost == "" {
        missingVars = append(missingVars, "DB_HOST")
    }
    if dbPort == "" {
        missingVars = append(missingVars, "DB_PORT")
    }
    if dbUser == "" {
        missingVars = append(missingVars, "DB_USER")
    }
    if dbPassword == "" {
        missingVars = append(missingVars, "DB_PASSWORD")
    }
    if dbName == "" {
        missingVars = append(missingVars, "DB_NAME")
    }
    
    if len(missingVars) > 0 {
        return nil, fmt.Errorf("missing required database configuration in .env file: %v", missingVars)
    }
    
    // If these are missing, use defaults
    if dbSSL == "" {
        dbSSL = "require"
    }
    if useDbPrefix == "" {
        useDbPrefix = "true"
    }
    
    log.Printf("Database config: host=%s, port=%s, user=%s, database=%s, sslmode=%s, useDbPrefix=%s",
        dbHost, dbPort, dbUser, dbName, dbSSL, useDbPrefix)
    
    return &Config{
        Host:        dbHost,
        Port:        dbPort,
        User:        dbUser,
        Password:    dbPassword,
        Database:    dbName,
        SSLMode:     dbSSL,
        UseDBPrefix: useDbPrefix == "true",
    }, nil
}

// LoadConfigWithPath loads configuration from a specific .env file path
func LoadConfigWithPath(envPath string) (*Config, error) {
    absPath, _ := filepath.Abs(envPath)
    log.Printf("Loading environment from specified path: %s", absPath)
    
    // Check if file exists
    if _, err := os.Stat(envPath); err != nil {
        return nil, fmt.Errorf("env file not found at %s: %w", envPath, err)
    }
    
    // Load env file
    if err := godotenv.Load(envPath); err != nil {
        return nil, fmt.Errorf("error loading env file %s: %w", envPath, err)
    }
    
    // Get database configuration
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbSSL := os.Getenv("DB_SSLMODE")
    useDbPrefix := os.Getenv("USE_DB_PREFIX")
    
    // Check that all required values are present
    missingVars := []string{}
    if dbHost == "" {
        missingVars = append(missingVars, "DB_HOST")
    }
    if dbPort == "" {
        missingVars = append(missingVars, "DB_PORT")
    }
    if dbUser == "" {
        missingVars = append(missingVars, "DB_USER")
    }
    if dbPassword == "" {
        missingVars = append(missingVars, "DB_PASSWORD")
    }
    if dbName == "" {
        missingVars = append(missingVars, "DB_NAME")
    }
    
    if len(missingVars) > 0 {
        return nil, fmt.Errorf("missing required database configuration in .env file: %v", missingVars)
    }
    
    // If these are missing, use defaults
    if dbSSL == "" {
        dbSSL = "require"
    }
    if useDbPrefix == "" {
        useDbPrefix = "true"
    }
    
    log.Printf("Database config: host=%s, port=%s, user=%s, database=%s, sslmode=%s, useDbPrefix=%s",
        dbHost, dbPort, dbUser, dbName, dbSSL, useDbPrefix)
    
    return &Config{
        Host:        dbHost,
        Port:        dbPort,
        User:        dbUser,
        Password:    dbPassword,
        Database:    dbName,
        SSLMode:     dbSSL,
        UseDBPrefix: useDbPrefix == "true",
    }, nil
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