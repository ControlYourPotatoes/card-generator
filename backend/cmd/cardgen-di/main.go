package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/pkg/bootstrap"
)

type CardOutput struct {
	ID   string        `json:"id"`
	Card *card.CardDTO `json:"card"`
}

func main() {
	// Command line flags for input/output
	inputFile := flag.String("input", "", "Input CSV file containing card definitions")
	outputFile := flag.String("output", "output/cards.json", "Output JSON file for processed cards")
	outputImageDir := flag.String("images", "output/cards", "Output directory for card images")
	clean := flag.Bool("clean", false, "Clean output directories before generation")
	cardType := flag.String("type", "creature", "Type of cards to parse (creature, spell, artifact, incantation, anthem)")
	env := flag.String("env", "development", "Environment (development, production, test)")
	flag.Parse()

	// Initialize application with DI
	app, err := bootstrap.NewApplication(*env)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer func() {
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Clean output directories if requested
	if *clean {
		if err := cleanOutputDirectories(*outputFile, *outputImageDir); err != nil {
			log.Printf("Warning: cleanup failed: %v", err)
		}
	}

	// Use default test cards if no input specified
	if *inputFile == "" {
		*inputFile = filepath.Join("test", "testdata", "test_cards.csv")
	}

	// Create output directories
	if err := os.MkdirAll(*outputImageDir, 0755); err != nil {
		log.Fatalf("Failed to create image output directory: %v", err)
	}

	// Process cards and collect results
	results, err := processCards(*inputFile, *cardType, *outputImageDir, app)
	if err != nil {
		log.Fatal(err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write results to JSON file
	if err := writeJSON(*outputFile, results); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully processed %d cards\n", len(results))
	fmt.Printf("Output written to: %s\n", *outputFile)
	fmt.Printf("Card images written to: %s\n", *outputImageDir)
}

func cleanOutputDirectories(jsonPath string, imagePath string) error {
	// Clean up function that handles errors gracefully
	cleanup := func(path string) error {
		// Check if path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil // Path doesn't exist, nothing to clean
		}

		// If it's a file, remove it
		if filepath.Ext(path) != "" {
			return os.Remove(path)
		}

		// For directories, remove all contents but keep the directory
		dir, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %w", path, err)
		}

		for _, d := range dir {
			fullPath := filepath.Join(path, d.Name())
			if err := os.RemoveAll(fullPath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", fullPath, err)
			}
		}

		return nil
	}

	// Clean JSON output
	if err := cleanup(jsonPath); err != nil {
		return fmt.Errorf("failed to clean JSON output: %w", err)
	}

	// Clean image directory
	if err := cleanup(imagePath); err != nil {
		return fmt.Errorf("failed to clean image directory: %w", err)
	}

	fmt.Println("Cleaned output directories")
	return nil
}

func processCards(filename string, cardType string, outputImageDir string, app *bootstrap.Application) ([]CardOutput, error) {
	// Open input file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// Get dependencies from DI container
	store, err := app.GetCardStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get card store: %w", err)
	}

	generator, err := app.GetCardGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to get card generator: %w", err)
	}

	csvParser, err := app.GetCSVParser(file)
	if err != nil {
		return nil, fmt.Errorf("failed to get CSV parser: %w", err)
	}

	// Parse cards
	cards, err := csvParser.ParseCSV(cardType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cards: %w", err)
	}

	var results []CardOutput

	// Process each card
	for _, c := range cards {
		// Save to store
		id, err := store.Save(c)
		if err != nil {
			return nil, fmt.Errorf("failed to save card %s: %w", c.GetName(), err)
		}

		// Convert to DTO for JSON output and image generation
		cardDTO := c.ToDTO()

		// Generate image for the card
		imagePath := filepath.Join(outputImageDir, id+".png")
		if err := generator.GenerateCard(cardDTO, imagePath); err != nil {
			return nil, fmt.Errorf("failed to generate image for card %s: %w", c.GetName(), err)
		}

		// Add to results
		results = append(results, CardOutput{
			ID:   id,
			Card: cardDTO,
		})

		fmt.Printf("Processed card: %s (ID: %s)\n", c.GetName(), id)
	}

	return results, nil
}

func writeJSON(filename string, cards []CardOutput) error {
	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create encoder with pretty printing
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	// Encode cards
	if err := encoder.Encode(cards); err != nil {
		return fmt.Errorf("failed to encode cards: %w", err)
	}

	return nil
}
