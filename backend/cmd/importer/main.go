package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/internal/parser"
	"github.com/ControlYourPotatoes/card-generator/internal/storage/database"
)

func main() {
	// Command line flags
	cardType := flag.String("type", "", "Type of cards to import (anthem, creature, spell, artifact, incantation)")
	inputFile := flag.String("file", "", "Path to the CSV file containing card data")
	dbEnvFile := flag.String("env", ".env", "Path to the environment file for database configuration")
	dryRun := flag.Bool("dry-run", false, "Parse the CSV and create cards but don't save to database")
	flag.Parse()

	// Validate input
	if *cardType == "" {
		log.Fatal("Card type is required. Use -type=anthem|creature|spell|artifact|incantation")
	}

	if *inputFile == "" {
		log.Fatal("Input file is required. Use -file=path/to/cards.csv")
	}

	// Normalize card type
	*cardType = strings.ToLower(*cardType)
	supportedTypes := map[string]bool{
		"anthem":      true,
		"creature":    true,
		"spell":       true,
		"artifact":    true,
		"incantation": true,
	}

	if !supportedTypes[*cardType] {
		log.Fatalf("Unsupported card type: %s. Supported types are: anthem, creature, spell, artifact, incantation", *cardType)
	}

	// Open input file
	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer file.Close()

	log.Printf("Importing %s cards from %s", *cardType, *inputFile)

	// Create CSV parser
	p := parser.NewCSVParser(file)

	// Parse the CSV into cards
	cards, err := p.ParseCSV(*cardType)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	log.Printf("Successfully parsed %d cards", len(cards))

	// If dry run, just display the cards and exit
	if *dryRun {
		fmt.Println("Dry run - cards will not be saved to database")
		for i, c := range cards {
			fmt.Printf("Card %d: %s (Cost: %d)\n", i+1, c.GetName(), c.GetCost())
		}
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

	// Initialize database if needed
	if err := dbManager.Initialize(""); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Get database store
	store, err := dbManager.GetStore()
	if err != nil {
		log.Fatalf("Failed to get database store: %v", err)
	}

	// Save cards to database
	successCount := 0
	for _, c := range cards {
		// Validate the card
		if err := c.Validate(); err != nil {
			log.Printf("WARNING: Card '%s' validation failed: %v", c.GetName(), err)
			continue
		}

		id, err := store.Save(c)
		if err != nil {
			log.Printf("WARNING: Failed to save card '%s': %v", c.GetName(), err)
			continue
		}

		fmt.Printf("Saved card: %s (ID: %s)\n", c.GetName(), id)
		successCount++
	}

	log.Printf("Successfully imported %d of %d cards", successCount, len(cards))
}