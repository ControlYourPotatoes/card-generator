package database

import (
    "fmt"
)

// Config holds database configuration
type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
    SSLMode  string
}

// ConnectionString returns a PostgreSQL connection string
func (c *Config) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
    )
}