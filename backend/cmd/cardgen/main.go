package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/image"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/parser"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/store/memory"
)

type CardOutput struct {
    ID   string         `json:"id"`
    Card *card.CardData `json:"card"`
}

func main() {
    // Command line flags for input/output
    inputFile := flag.String("input", "", "Input CSV file containing card definitions")
    outputFile := flag.String("output", "output/cards.json", "Output JSON file for processed cards")
    outputImageDir := flag.String("images", "output/cards", "Output directory for card images")
    clean := flag.Bool("clean", false, "Clean output directories before generation")
    flag.Parse()

    if *clean {
        if err := cleanOutputDirectories(*outputFile, *outputImageDir); err != nil {
            log.Printf("Warning: cleanup failed: %v", err)
        }
    }

    if *inputFile == "" {
        // Use default test cards if no input specified
        *inputFile = filepath.Join("test", "testdata", "test_cards.csv")
    }

    // Initialize components
    store := memory.New()
    defer store.Close()

    // Create card factory
    factory := card.NewCardFactory(store)

    // Create image generator
    generator, err := image.NewGenerator()
    if err != nil {
        log.Fatalf("Failed to create image generator: %v", err)
    }

    // Create output directories
    if err := os.MkdirAll(*outputImageDir, 0755); err != nil {
        log.Fatalf("Failed to create image output directory: %v", err)
    }

    // Process cards and collect results
    results, err := processCards(*inputFile, factory, store, generator, *outputImageDir)
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
func processCards(filename string, factory *card.CardFactory, store card.CardStore, generator *image.Generator, outputImageDir string) ([]CardOutput, error) {
    // Open input file
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open input file: %w", err)
    }
    defer file.Close()

    // Create parser
    p := parser.NewParser(file)

    // Parse cards
    cards, err := p.Parse()
    if err != nil {
        return nil, fmt.Errorf("failed to parse cards: %w", err)
    }

    var results []CardOutput

    // Process each card
    for _, c := range cards {
        // Convert to CardData
        data := card.ToData(c)

        // Create card using factory
        newCard, err := factory.CreateFromData(data)
        if err != nil {
            return nil, fmt.Errorf("failed to create card %s: %w", c.GetName(), err)
        }

        // Save to store
        id, err := store.Save(newCard)
        if err != nil {
            return nil, fmt.Errorf("failed to save card %s: %w", newCard.GetName(), err)
        }

        // Generate image for the card
        imagePath := filepath.Join(outputImageDir, id+".png")
        if err := generator.GenerateImage(data, imagePath); err != nil {
            return nil, fmt.Errorf("failed to generate image for card %s: %w", newCard.GetName(), err)
        }

        // Add to results
        results = append(results, CardOutput{
            ID:   id,
            Card: data,
        })

        fmt.Printf("Processed card: %s (ID: %s)\n", newCard.GetName(), id)
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