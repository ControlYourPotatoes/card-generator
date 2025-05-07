// Package main provides a utility to clean the database
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ControlYourPotatoes/card-generator/internal/storage/database"
)

func main() {
	// Command line flags
	dbEnvFile := flag.String("env", ".env", "Path to the environment file for database configuration")
	confirm := flag.Bool("confirm", false, "Confirmation flag to prevent accidental cleanup")
	flag.Parse()

	if !*confirm {
		fmt.Println("WARNING: This will delete ALL cards from the database.")
		fmt.Println("Run with -confirm flag to proceed.")
		return
	}

	// Setup database connection
	dbManager, err := database.NewManager(*dbEnvFile)
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}

	if err := dbManager.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	// Initialize the database first
	if err := dbManager.Initialize(""); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Get database store for direct connection
	store, err := dbManager.GetStore()
	if err != nil {
		log.Fatalf("Failed to get database store: %v", err)
	}

	// Clean the database
	deleted, err := cleanDatabase(store)
	if err != nil {
		log.Fatalf("Failed to clean database: %v", err)
	}

	log.Printf("Successfully deleted %d cards from the database", deleted)
}

func cleanDatabase(store *database.PostgresStore) (int, error) {
	// Get the connection pool from the store
	pool := store.GetPool()
	if pool == nil {
		return 0, fmt.Errorf("database connection pool is nil")
	}

	ctx := context.Background()

	// Delete in correct order to respect foreign key constraints
	// First delete from dependent tables
	tables := []string{
		"card_keywords",
		"card_metadata",
		"creature_cards",
		"artifact_cards",
		"spell_cards",
		"incantation_cards",
		"anthem_cards",
		"card_set_cards",
		"card_images",
		// Then delete from the main cards table
		"cards",
	}

	var totalDeleted int

	for _, table := range tables {
		tag, err := pool.Exec(
			ctx,
			fmt.Sprintf("DELETE FROM %s", table),
		)
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to clean table %s: %w", table, err)
		}
		if table == "cards" {
			totalDeleted = int(tag.RowsAffected())
		}
		fmt.Printf("Cleaned table: %s\n", table)
	}

	return totalDeleted, nil
}