package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	fmt.Println("Testing database connection...")

	// Use the connection string format with the db. prefix for Supabase
	connString := "postgresql://postgres:mxGzTCKS38jyBCQ8@db.qkkixxahqhuhwvtnqokb.supabase.co:5432/postgres?sslmode=require"

	// Create a connection
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Test query
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	fmt.Println("Database version:", version)

	// Try to create a sample table
	fmt.Println("Testing table creation...")
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS test_connection (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Insert a row
	_, err = conn.Exec(context.Background(), `
		INSERT INTO test_connection (name) VALUES ($1)
	`, "Connection Test")
	if err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	// Query the data
	var id int
	var name string
	err = conn.QueryRow(context.Background(), `
		SELECT id, name FROM test_connection ORDER BY id DESC LIMIT 1
	`).Scan(&id, &name)
	if err != nil {
		log.Fatalf("Failed to query test data: %v", err)
	}

	fmt.Printf("Test data retrieved: ID=%d, Name=%s\n", id, name)
	fmt.Println("Database connection and operations successful!")
}
